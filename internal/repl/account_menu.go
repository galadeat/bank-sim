package repl

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"time"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/gofrs/uuid"
)

func runAccountMenu(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {
	for {
		fmt.Println("\n\t\t\t\tAccount Menu")
		fmt.Println("1) Create Account")
		fmt.Println("2) Get Account")
		fmt.Println("3) List Accounts")
		fmt.Println("4) Delete Account")
		fmt.Println("5) Deposit Money")
		fmt.Println("6) Withdraw Money")
		fmt.Println("7) Back")

		choice := readInput(reader, "Choose option: ")

		switch choice {
		case "1":
			handleCreateAccount(reader, accountClient, userClient)
		case "2":
			handleGetAccount(reader, accountClient, userClient)
		case "3":
			handleListAccounts(reader, accountClient, userClient)
		case "4":
			handleDeleteAccount(reader, accountClient, userClient)
		case "5":
			handleDepositMoney(reader, accountClient, userClient)
		case "6":
			handleWithdrawMoney(reader, accountClient, userClient)
		case "7":
			return
		default:
			fmt.Println("Invalid choice")

		}
	}

}

func handleCreateAccount(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {

	userId := runChooseUserMenu(reader, userClient)
	if userId == "" {
		return
	}
	balance := runBalanceMenu(reader)
	request, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return
	}
	req := &accountv2.CreateAccountRequest{
		UserId:         userId,
		InitialBalance: balance,
		RequestId:      request.String(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := accountClient.CreateAccount(ctx, req)
	if err != nil {
		fmt.Printf("CreateAccount error: %v\n", err)
		return
	}
	fmt.Printf("Account created: %v\n", resp.GetAccount().GetId())

}

func handleGetAccount(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {
	id := runChooseAccountMenu(reader, accountClient, userClient)
	if id == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := accountClient.GetAccount(ctx, &accountv2.GetAccountRequest{Id: id})
	if err != nil {
		fmt.Printf("Error getting account: %v\n", err)
		return
	}
	fmt.Printf("\nAccount id: %s\n", resp.GetAccount().GetId())
	fmt.Printf("Account owner: %v\n", resp.GetAccount().GetOwner().GetId())
	fmt.Printf("Account balance: %v\n", resp.GetAccount().GetBalance())
}

func handleListAccounts(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {

	id := runChooseUserMenu(reader, userClient)
	if id == "" {
		return
	}
	req := &accountv2.ListAccountsRequest{
		UserId: id,
	}
	fmt.Printf("\n\tAccounts:")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := accountClient.ListAccounts(ctx, req)
	if err != nil {
		fmt.Printf("Error listing accounts: %v\n", err)
		return
	}
	for i, a := range resp.Accounts {
		fmt.Printf("Account %d: %v", i+1, a.GetId())
	}
}

func handleDeleteAccount(reader *bufio.Reader, accountClient accountv2.AccountClient, client userv1.UserClient) {

	id := runChooseAccountMenu(reader, accountClient, client)
	if id == "" {
		return
	}
	req := &accountv2.DeleteAccountRequest{
		AccountId: id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := accountClient.DeleteAccount(ctx, req)
	if err != nil {
		fmt.Printf("Error deleting account: %v\n", err)
		return
	}
	fmt.Printf("\nAccount deleted: %v\n", id)
}
func handleDepositMoney(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {
	id := runChooseAccountMenu(reader, accountClient, userClient)
	if id == "" {
		return
	}
	depositMoney := runBalanceMenu(reader)
	reqId, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return
	}
	req := &accountv2.DepositRequest{
		AccountId: id,
		Amount:    depositMoney,
		RequestId: reqId.String(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := accountClient.Deposit(ctx, req)
	if err != nil {
		fmt.Printf("Error depositing money: %v\n", err)
		return
	}
	fmt.Printf("\nAccount deposited: %v\n", resp.GetAccount().GetId())
	fmt.Printf("New balance: %v\n", resp.GetAccount().GetBalance())

}

func handleWithdrawMoney(reader *bufio.Reader, accountClient accountv2.AccountClient, userClient userv1.UserClient) {

	id := runChooseAccountMenu(reader, accountClient, userClient)
	if id == "" {
		return
	}
	withdrawMoney := runBalanceMenu(reader)
	reqId, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return
	}
	req := &accountv2.WithdrawRequest{
		AccountId: id,
		Amount:    withdrawMoney,
		RequestId: reqId.String(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := accountClient.Withdraw(ctx, req)
	if err != nil {
		fmt.Printf("Error withdrawing money: %v\n", err)
	}
	fmt.Printf("\nAccount withdrawed: %v\n", resp.GetAccount().GetId())
	fmt.Printf("New balance: %v\n", resp.GetAccount().GetBalance())
}
