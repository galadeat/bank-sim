package repl

import (
	"bufio"
	"context"
	"fmt"
	"strconv"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	commonv1 "github.com/galadeat/bank-sim/api/proto/common/v1"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
)

// runBalanceMenu repl function to initialize balance.
func runBalanceMenu(reader *bufio.Reader) *commonv1.Money {
	var currency = "None"
	var units int64 = 0
	var nanos int32 = 0
	for {
		fmt.Println("\n\tBalance:")
		fmt.Printf("1) Currency: %s", currency)
		fmt.Printf("\n2) Units: %d", units)
		fmt.Printf("\n3) Nanos: %d", nanos)
		fmt.Println("\n4) Accept")

		choice := readInput(reader, "Choose option: ")
		switch choice {
		case "1":
			currency = readInput(reader, "Enter currency: ")
		case "2":
			input, err := strconv.Atoi(readInput(reader, "Enter units: "))
			if err != nil || input < 0 {
				fmt.Println("Enter valid number")
				continue
			}
			units = int64(input)
		case "3":
			input, err := strconv.Atoi(readInput(reader, "Enter nanos: "))
			if err != nil || input < 0 {
				fmt.Println("Enter valid number")
				continue
			}
			nanos = int32(input)
		case "4":
			if currency == "None" {
				fmt.Println("Currency must be chosen")
				continue
			}
			return &commonv1.Money{
				Currency: currency,
				Units:    units,
				Nanos:    nanos,
			}
		default:
			fmt.Println("Enter valid option")

		}
	}
}

//

func runChooseUserMenu(reader *bufio.Reader, client userv1.UserClient) string {
	fmt.Println("\n\tUsers:")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := client.ListUsers(ctx, &userv1.ListUsersRequest{})
	if err != nil {
		fmt.Println("Error listing users: ", err)
		return ""
	}
	if len(resp.Users) == 0 {
		fmt.Println("No users found")
		return ""
	}
	for i, user := range resp.Users {
		fmt.Printf("%d) User ID: %s\n", i+1, user.GetId())
	}
	choice, err := strconv.Atoi(readInput(reader, "Choose User: "))
	if err != nil {
		fmt.Println("Error choosing user: ", err)
		return ""
	}
	if choice-1 < 0 || choice-1 >= len(resp.Users) {
		fmt.Println("Invalid choice")
		return ""
	}

	return resp.Users[choice-1].GetId()

}

func runChooseAccountMenu(reader *bufio.Reader, client accountv2.AccountClient, userClient userv1.UserClient) string {
	id := runChooseUserMenu(reader, userClient)
	if id == "" {
		return ""
	}
	fmt.Println("\n\tAccounts:")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resp, err := client.ListAccounts(ctx, &accountv2.ListAccountsRequest{UserId: id})
	if err != nil {
		fmt.Println("Error listing accounts: ", err)
		return ""
	}
	if len(resp.Accounts) == 0 {
		fmt.Println("No accounts found")
		return ""
	}

	for i, account := range resp.Accounts {
		fmt.Printf("%d) Account ID: %s\n", i+1, account.GetId())
	}
	choice, err := strconv.Atoi(readInput(reader, "Choose Account: "))
	if err != nil {
		fmt.Println("Error choosing account: ", err)
		return ""
	}
	if choice-1 < 0 || choice-1 >= len(resp.Accounts) {
		fmt.Println("Invalid choice")
		return ""
	}

	return resp.Accounts[choice-1].GetId()
}
