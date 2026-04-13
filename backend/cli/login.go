package cli

import (
	"fmt"
	"log"

	"github.com/lidchen/neuron_deck/backend/db"
)

func (a *CliApp) handleLogin(args []string) {
	// no args: prompt user name and user password
	// 1 args: prompt user password
	// 2 args: try login

	var username, password string
	if len(args) == 0 {
		username = readLineWithPrompt("username: ")
		password = readLineWithPrompt("password: ")
	}
	if len(args) == 1 {
		username = args[0]
		password = readLineWithPrompt("password: ")
	}
	if len(args) >= 2 {
		username = args[0]
		password = args[1]
	}

	valid_user, err := db.GetUserByUsername(a.db, username)
	if err != nil {
		if err.Code == "NOT_FOUND" {
			fmt.Println("User not exist")
			return
		}
		log.Fatal(err)
		return
	}
	if valid_user.Password != password {
		fmt.Println("password mismatch")
		return
	}
	a.user = valid_user
}
