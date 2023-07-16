package main

import (
	"log"
	"main/client"
)

func main() {
	clientConfig := &client.Config{
		Hostname: "localhost",
		Port:     "8080",
		Resource: "https://example.com",
	}
	
	tcpClient := client.NewClient(clientConfig)
	err := tcpClient.Start()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
