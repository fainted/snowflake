package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/fainted/snowflake/api/thrift/protocols"
)

func getClient() *protocols.SnowflakeClient {
	var transport thrift.TTransport
	var err error

	transport, err = thrift.NewTSocket(flagAddr)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		panic("getClient fail")
	}

	transportFactory := thrift.NewTTransportFactory()
	transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport = transportFactory.GetTransport(transport)

	if err := transport.Open(); err != nil {
		panic("open transport")
	}

	return protocols.NewSnowflakeClientFactory(transport, protocolFactory)
}

func getNext(client *protocols.SnowflakeClient) {
	start := time.Now().UnixNano()
	resp, err := client.GetNextID()
	costs := int64((time.Now().UnixNano() - start) / 1000)
	if err != nil {
		fmt.Println("RPC[GetNextID] error", err)
	} else {
		fmt.Printf("RPC[GetNextID] returns[0x%x, %d], costs[%d]us\n", *resp.ID, *resp.ID, costs)
	}
}

func main() {
	flag.Parse()

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			client := getClient()
			defer client.Transport.Close()
			if client == nil {
				panic("get a nil client")
			}
			for i := 0; i < 1000; i++ {
				getNext(client)
			}
		}(&wg)
	}

	wg.Wait()
}
