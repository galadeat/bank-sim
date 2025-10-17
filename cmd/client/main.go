package main

import (
	"context"
	"log"
	"os"
	"time"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	commonv1 "github.com/galadeat/bank-sim/api/proto/common/v1"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/galadeat/bank-sim/pkg/clients"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clients, err := clients.New()

	if err != nil {
		log.Fatal("falied to init clients: %v", err)
	}
	defer clients.Close()

	// create user
	login := os.Getenv("LOGIN")
	email := os.Getenv("EMAIL")
	if login == "" || email == "" {
		log.Fatal("LOGIN and EMAIL must be set in .env")
	}

	user, err := clients.User.CreateUser(ctx, &userv1.CreateUserRequest{
		Login: login,
		Email: email,
	})
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	// retrieve user using returned ID
	_, err = clients.User.GetUser(ctx, &userv1.GetUserRequest{Id: user.Id})
	if err != nil {
		log.Fatalf("Could not get account: %v", err)
	}

	// Create accounts

	accRUB, err := clients.Account.CreateAccount(ctx, &accountv2.CreateAccountRequest{
		UserId: user.Id,
		InitialBalance: &commonv1.Money{
			Currency: "RUB",
			Units:    1_000_000,
			Nanos:    0,
		},
		RequestId: "1",
	})
	if err != nil {
		log.Fatalf("Could not create account: %v", err)
	}

	accUSD, err := clients.Account.CreateAccount(ctx, &accountv2.CreateAccountRequest{
		UserId: user.Id,
		InitialBalance: &commonv1.Money{
			Currency: "USD",
			Units:    1_000_000,
			Nanos:    0,
		},
		RequestId: "2",
	})
	if err != nil {
		log.Fatalf("Could not create account: %v", err)
	}

	// withdraw money
	_, err = clients.Account.Withdraw(ctx, &accountv2.WithdrawRequest{
		AccountId: accUSD.Account.GetId(),
		Amount:    &commonv1.Money{Currency: "USD", Units: 1_000, Nanos: 100},
		RequestId: "3",
	})
	if err != nil {
		log.Fatalf("Could not withdraw: %v", err)
	}

	// deposit money
	_, err = clients.Account.Deposit(ctx, &accountv2.DepositRequest{
		AccountId: accRUB.Account.GetId(),
		Amount:    &commonv1.Money{Currency: "RUB", Units: 1_000, Nanos: 100},
		RequestId: "4",
	})
	if err != nil {
		log.Fatalf("Could not deposit: %v", err)
	}

}
