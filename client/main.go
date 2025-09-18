package main

import (
	"context"
	"log"
	"os"
	"time"

	accountv1 "github.com/galadeat/bank-sim/api/proto/account/v1"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// establish gRPC connection using address from .env
	conn, err := grpc.NewClient(os.Getenv("GRPC_ADDRESS"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect^ %v", err)
	}
	defer conn.Close()

	// initialize client
	client := accountv1.NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// add account
	login := os.Getenv("LOGIN")
	email := os.Getenv("EMAIL")
	if login == "" || email == "" {
		log.Fatal("LOGIN and EMAIL must be set in .env")
	}

	response, err := client.AddAccount(ctx, &accountv1.AccountInfo{Login: login, Email: email})
	if err != nil {
		log.Fatalf("Could not add account: %v", err)
	}
	log.Printf("Account ID: %s added successfully", response.Value)

	// retrieve account using returned ID
	account, err := client.GetAccount(ctx, &accountv1.AccountID{Value: response.Value})
	if err != nil {
		log.Fatalf("Could not get account: %v", err)
	}
	log.Print("Account: ", account.String())

}
