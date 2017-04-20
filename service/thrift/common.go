package main

import (
	"flag"
)

var (
	flagAddr string
)

func init() {
	flag.StringVar(&flagAddr, "addr", "127.0.0.1:9090", "Service address")
}
