package main

import (
	"log"
	"main/client"
)

func main() {
	tcpClient := client.NewClient("localhost", "8080", "https://example.com")
	err := tcpClient.Start()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
