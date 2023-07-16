package main

import (
	"log"
	hashcash2 "main/internal/hashcash"
	quote2 "main/internal/quote"
	"main/internal/shared/db"
	"main/server"
)

func main() {
	database := db.NewDB()
	hashRepository := hashcash2.NewRepository(database)
	hashService := hashcash2.NewService(hashRepository)

	quoteRepository := quote2.NewRepository()
	quoteService := quote2.NewService(quoteRepository)

	tcpServer := server.NewServer(hashService, quoteService)

	log.Fatal(tcpServer.Listen(":8080"))
}
