package main

import (
	"log"

	"github.com/galadeat/bank-sim/internal/repl"
	"github.com/galadeat/bank-sim/pkg/clients"
	"github.com/galadeat/bank-sim/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	file := logger.Init("appClient.log")
	defer file.Close()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Print("Logs started")
	clients, err := clients.New()

	if err != nil {
		log.Fatalf("falied to init clients: %v", err)
	}
	defer clients.Close()

	// run REPL
	repl.Run(clients.User, clients.Account)

}
