package main

import (
	"flag"
	"log"
)

const defaultListenAddr = ":5001"

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of the goredis server")
	flag.Parse()

	server := NewServer(*listenAddr)
	log.Fatal(server.Start())
}
