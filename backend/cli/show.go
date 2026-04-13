package cli

import (
	"fmt"
	"log"

	"github.com/lidchen/neuron_deck/backend/db"
)

func (a *CliApp) handleShowCards(args []string) {
	// default print all cards
	if a.deck == nil {
		fmt.Println("No deck is opened, please open a deck first")
		return
	}
	cards, err := db.GetCards(a.db, a.deck.Id)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, c := range cards {
		fmt.Println(c)
	}
}

func (a *CliApp) handleShowDecks(args []string) {
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}

	decks, err := db.ListDecksByUserID(a.db, a.user.Id)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(decks) == 0 {
		fmt.Println("No decks")
		return
	}

	for _, d := range decks {
		fmt.Println(d)
	}
}

func (a *CliApp) handleShowUsers(args []string) {
	if a.user == nil {
		fmt.Println("No user logged in")
		return
	}
	fmt.Println(*a.user)
}

func (a *CliApp) handleShow(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", show_usage)
		return
	}

	target := args[0]
	switch target {
	case "card", "cards":
		a.handleShowCards(args[1:])
	case "deck", "decks":
		a.handleShowDecks(args[1:])
	case "user", "users":
		a.handleShowUsers(args[1:])
	default:
		fmt.Printf("Unknown show target: %q\nUsage:\n%s\n", target, show_usage)
	}
}
