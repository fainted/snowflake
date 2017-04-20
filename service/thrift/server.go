package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/golang/glog"

	"github.com/fainted/snowflake"
	"github.com/fainted/snowflake/api/thrift/protocols"
)

var (
	flagWorkerID int64
)

var (
	snowflakeWorker snowflake.Worker
)

func registerSignalHandler(server *thrift.TSimpleServer) {
	if server == nil {
		return
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c
		fmt.Fprintf(os.Stderr, "Receive signal [%d]\n", s)

		server.Stop()
	}()
}

func buildServer(address string) (*thrift.TSimpleServer, error) {
	serverTransport, err := thrift.NewTServerSocket(address)
	if err != nil {
		return nil, err
	}

	// worker, err := snowflake.NewAtomicWorker(flagWorkerID)
	worker, err := snowflake.NewChannelWorker(flagWorkerID)
	if err != nil {
		return nil, err
	}

	handler, err := NewSnowflakeHandler(worker)
	if err != nil {
		return nil, err
	}

	return thrift.NewTSimpleServer4(
		protocols.NewSnowflakeProcessor(handler), // processor
		serverTransport,                          // serverTransport
		thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()), // transportFactory
		thrift.NewTBinaryProtocolFactoryDefault(),                        // protocolFactory
	), nil
}

func main() {
	flag.Parse()

	server, err := buildServer(flagAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server @[%s]\n", flagAddr)
	}

	registerSignalHandler(server)
	server.Serve()

	glog.Flush()
}

func init() {
	flag.Int64Var(&flagWorkerID, "worker_id", 1, "Snowflake worker ID")
}
