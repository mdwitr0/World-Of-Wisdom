package main

import (
	"log"
	"main/client"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	clientConfig := &client.Config{
		Hostname: host,
		Port:     "8080",
		Resource: "https://example.com",
	}

	tcpClient := client.NewClient(clientConfig)
	err := tcpClient.Start()
	if err != nil {
		log.Fatalf("failed to start client: %s", err.Error())
	}
}
