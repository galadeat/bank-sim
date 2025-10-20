package repl

import (
	"bufio"
	"fmt"
	"os"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
)

func Run(userClient userv1.UserClient, accountClient accountv2.AccountClient) {
	fmt.Println("\n\n\t\t\tWelcome to Bank Sim REPL")
	reader := bufio.NewReader(os.Stdin)
	for {
		switch showMainMenu(reader) {
		case "user":
			runUserMenu(reader, userClient)
		case "account":
			runAccountMenu(reader, accountClient, userClient)
		case "exit":
			fmt.Println("Bye! Thanks for using this application!")
			os.Exit(0)
		}
	}
}
