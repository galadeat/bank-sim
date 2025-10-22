package main

import (
	"log"

	"github.com/galadeat/bank-sim/internal/repl"
	"github.com/galadeat/bank-sim/pkg/clients"
	"github.com/galadeat/bank-sim/pkg/logger"
)

func main() {
	file := logger.Init("appClient.log")
	defer file.Close()
	
	log.Print("Logs started")
	clients, err := clients.New()

	if err != nil {
		log.Fatalf("falied to init clients: %v", err)
	}
	defer clients.Close()

	// run REPL
	repl.Run(clients.User, clients.Account)

}
