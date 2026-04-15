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

	valid_user, success, err := db.Login(a.db, username, password)
	if err != nil {
		log.Fatal(err.Message)
		return
	}
	if !success {
		fmt.Println("password or username incorrect")
	}
	a.user = valid_user
}
