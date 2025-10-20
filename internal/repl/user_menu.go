package repl

import (
	"bufio"
	"context"
	"fmt"
	"time"

	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func runUserMenu(reader *bufio.Reader, userClient userv1.UserClient) {
	for {
		fmt.Println("\n\t\t\t\tUser Menu")
		fmt.Println("1) Create User")
		fmt.Println("2) Get User")
		fmt.Println("3) List Users")
		fmt.Println("4) Update User")
		fmt.Println("5) Delete User")
		fmt.Println("6) Back")

		choice := readInput(reader, "Choose option: ")

		switch choice {
		case "1":
			handleCreateUser(reader, userClient)
		case "2":
			handleGetUser(reader, userClient)
		case "3":
			handleListUsers(userClient)
		case "4":
			handleUpdateUser(reader, userClient)
		case "5":
			handleDeleteUser(reader, userClient)
		case "6":
			return
		default:
			fmt.Println("Invalid choice")

		}
	}

}

func handleCreateUser(reader *bufio.Reader, userClient userv1.UserClient) {
	login := readInput(reader, "Enter your login: ")
	email := readInput(reader, "Enter your email: ")

	req := &userv1.CreateUserRequest{
		Login: login,
		Email: email,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := userClient.CreateUser(ctx, req)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return
	}
	fmt.Printf("User %s created.\n", resp.GetId())

}

func handleGetUser(reader *bufio.Reader, userClient userv1.UserClient) {
	id := runChooseUserMenu(reader, userClient)
	if id == "" {
		return
	}
	resp, err := userClient.GetUser(context.Background(), &userv1.GetUserRequest{Id: id})
	if err != nil {
		fmt.Println("Error getting user: ", err)
		return
	}

	fmt.Printf("UserdID: %s\n", resp.GetUser().GetId())
	fmt.Printf("Login: %s\n", resp.GetUser().GetLogin())
	fmt.Printf("Email: %s\n", resp.GetUser().GetEmail())

}

func handleListUsers(userClient userv1.UserClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("\n\t\t\t\tUsers:")
	resp, err := userClient.ListUsers(ctx, &userv1.ListUsersRequest{})
	if err != nil {
		fmt.Println("Error listing users: ", err)
		return
	}
	if len(resp.Users) == 0 {
		fmt.Println("No users found")
		return
	}
	for i, user := range resp.Users {
		fmt.Printf("\nUser %d ID: %s", i+1, user.Id)
	}
	fmt.Print("\n")
}

func handleUpdateUser(reader *bufio.Reader, userClient userv1.UserClient) {
	id := runChooseUserMenu(reader, userClient)
	if id == "" {
		return
	}
	login := readInput(reader, "Enter new login or press Enter to skip: ")
	email := readInput(reader, "Enter new email or press Enter to skip: ")

	req := &userv1.UpdateUserRequest{
		Login: wrapperspb.String(login),
		Email: wrapperspb.String(email),
		Id:    id,
	}

	resp, err := userClient.UpdateUser(context.Background(), req)
	if err != nil {
		fmt.Println("Error updating user: ", err)
		return
	}
	fmt.Printf("User %s updated\n", resp.GetUser().GetId())

}
func handleDeleteUser(reader *bufio.Reader, userClient userv1.UserClient) {
	id := runChooseUserMenu(reader, userClient)
	if id == "" {
		return
	}

	req := &userv1.DeleteUserRequest{Id: id}
	_, err := userClient.DeleteUser(context.Background(), req)
	if err != nil {
		fmt.Println("Error deleting user: ", err)
		return
	}
	fmt.Printf("User %s deleted", id)
}
