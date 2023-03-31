package main

import (
	"log"
	"wbL0/orderPub/service"
)

const (
	NatsStrUrl = "localhost:4223"
	clusterId  = "test-cluster"
	clientId   = "test-publisher"
	ch         = "testch"
)

func main() {
	nc := service.CreateStan()
	err := nc.Connect(clusterId, clientId, NatsStrUrl)
	defer nc.Close()
	if err != nil {
		log.Println("Error while connecting to nats")
		panic(err)
	}
	_ = nc.PublishFromStdinCycle(ch)
}
