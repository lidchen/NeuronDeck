package cli

import (
	"fmt"
	"log"

	"github.com/lidchen/neuron_deck/backend/db"
)

// TODO: create card sys when create card
func (a *CliApp) handleCreateCard(args []string) {
	// no args: prompt for front and back
	// 1 arg: error
	// 2 arg: create card
	// 2+ arg: ignore others
	if a.deck == nil {
		fmt.Println("No deck is opened, please open a deck first")
		return
	}
	var front, back string
	if len(args) == 0 {
		front = readLineWithPrompt("front: ")
		back = readLineWithPrompt("back: ")
	}
	if len(args) == 1 {
		fmt.Printf("Expect 2 args for create card\nUsage:\n%s\n", create_usage)
		return
	}
	if len(args) >= 2 {
		front = args[0]
		back = args[1]
	}
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}
	if err := a.validateDeck(); err != nil {
		fmt.Println(err)
		return
	}
	err := db.CreateCard(a.db, a.deck.Id, front, back, nil, false)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Card Created")
}

func (a *CliApp) handleCreateDeck(args []string) {
	// no args: prompt for deckname
	// 1 arg: create deck
	// 1+ arg: ignore others
	var deckname string
	if len(args) == 0 {
		deckname = readLineWithPrompt("deckname: ")
	}
	if len(args) >= 1 {
		deckname = args[0]
	}
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}
	newDeck, err := db.CreateDeck(a.db, a.user.Id, deckname)
	if err != nil {
		if err.Code == "DECK_ALREADY_EXISTS" {
			fmt.Println("deck already exists")
			return
		}
		log.Fatal(err)
		return
	}
	a.deck = newDeck
	fmt.Println("Deck Created")
}

func (a *CliApp) handleCreateUser(args []string) {
	// no args: prompt for username and password
	// 1 arg: error
	// 2 arg: create user
	// 2+ arg: ignore others
	var username, password string
	if len(args) == 0 {
		username = readLineWithPrompt("username: ")
		password = readLineWithPrompt("password: ")
	}
	if len(args) == 1 {
		fmt.Printf("Expect 2 args for create user\nUsage:\n%s\n", create_usage)
		return
	}
	if len(args) >= 2 {
		username = args[0]
		password = args[1]
	}
	newUser, err := db.CreateUser(a.db, username, password)
	if err != nil {
		if err.Code == "USER_ALREADY_EXISTS" {
			fmt.Println("user already exists")
			return
		}
		log.Fatal(err)
		return
	}
	a.user = newUser
	fmt.Println("User Created")
}

func (a *CliApp) handleCreate(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", create_usage)
		return
	}

	target := (args)[0]
	switch target {
	case "card":
		a.handleCreateCard(args[1:])
	case "deck":
		a.handleCreateDeck(args[1:])
	case "user":
		a.handleCreateUser(args[1:])
	default:
		fmt.Printf("Unknown create target: %q\nUsage:\n%s\n", target, create_usage)
	}
}
