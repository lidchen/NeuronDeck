package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
)

func (a *CliApp) handleDeleteCard(args []string) {
	if a.user == nil {
		fmt.Println(ErrNoLogin.Error())
		return
	}
	if a.deck == nil {
		fmt.Println(ErrNoDeckOpen.Error())
		return
	}

	var cardIDStr string
	if len(args) == 0 {
		cardIDStr = readLineWithPrompt("card id: ")
	} else {
		cardIDStr = args[0]
	}

	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		fmt.Println("invalid card id")
		return
	}

	errApp := db.DeleteCard(a.db, a.deck.Id, cardID)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("card not found")
			return
		}
		log.Fatal(errApp)
		return
	}

	fmt.Println("Card Deleted")
}

func (a *CliApp) handleDeleteDeck(args []string) {
	if a.user == nil {
		fmt.Println(ErrNoLogin.Error())
		return
	}
	var deckIDStr string
	if len(args) == 0 {
		deckIDStr = readLineWithPrompt("deck id: ")
	} else {
		deckIDStr = args[0]
	}

	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		fmt.Println("invalid deck id")
		return
	}

	errApp := db.DeleteDeck(a.db, a.user.Id, deckID)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("deck not found")
			return
		}
		log.Fatal(errApp)
		return
	}

	if a.deck != nil && a.deck.Id == deckID {
		a.deck = nil
	}

	fmt.Println("Deck Deleted")
}

func (a *CliApp) handleDelete(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", delete_usage)
		return
	}

	target := args[0]
	switch target {
	case "card", "cards":
		a.handleDeleteCard(args[1:])
	case "deck", "decks":
		a.handleDeleteDeck(args[1:])
	default:
		fmt.Printf("Unknown delete target: %q\nUsage:\n%s\n", target, delete_usage)
	}
}
