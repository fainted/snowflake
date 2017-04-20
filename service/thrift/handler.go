package main

import (
	"errors"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/golang/glog"

	"github.com/fainted/snowflake"
	"github.com/fainted/snowflake/api/thrift/protocols"
)

type SnowflakeHandler struct {
	worker snowflake.Worker
}

func NewSnowflakeHandler(w snowflake.Worker) (*SnowflakeHandler, error) {
	if w == nil {
		return nil, errors.New("Nil snowflake.Worker")
	}

	return &SnowflakeHandler{worker: w}, nil
}

func (h *SnowflakeHandler) GetNextID() (r *protocols.SimpleResponse, err error) {
	retval := protocols.SimpleResponse{OK: false, ID: nil}
	id, err := h.worker.Next()
	if nil == err {
		retval.OK, retval.ID = true, thrift.Int64Ptr(id)
	}

	glog.Infof(`GetNextID response{"OK": %t, ID: %x}`, retval.OK, *retval.ID)
	return &retval, nil
}

func (h *SnowflakeHandler) Ping() error {
	return nil
}
