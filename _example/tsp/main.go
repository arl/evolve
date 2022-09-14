package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if _, err := os.Stat("index.html"); os.IsNotExist(err) {
		log.Fatal("tsp example should be run from _example/tsp directory")
	}

	host := "localhost:8080"
	flag.StringVar(&host, "host", host, "serve web gui at [host]:[port]")
	flag.Parse()

	server := newServer()
	server.serve(host)
	if err := server.start(); err != nil {
		log.Fatal(err)
	}
}
