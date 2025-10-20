package repl

import (
	"bufio"
	"fmt"
)

func showMainMenu(reader *bufio.Reader) string {
	for {
		fmt.Println("\n\t\t\t\tMain Menu")
		fmt.Println("1) User Service")
		fmt.Println("2) Account Service")
		fmt.Println("3) Exit")

		choice := readInput(reader, "Choose option: ")
		switch choice {
		case "1":
			return "user"
		case "2":
			return "account"
		case "3":
			return "exit"
		default:
			fmt.Println("Please enter a valid choice")
		}
	}
}
