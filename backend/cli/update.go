package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
)

func (a *CliApp) handleUpdateCard(args []string) {
	// updated value: '_' means no update this
	// <id>(must) <newfront>/'_' <newback>/'_'
	// 0 arg: prompt for id, updated value
	// 1 arg: expect this as id, prompt for updated value
	// 2 args: error
	// 3 / 3+ args: expect <id> <front> <back> format, ignore others
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}
	if err := a.validateDeck(); err != nil {
		fmt.Println(err)
		return
	}

	var cardIDStr, front, back string
	switch len(args) {
	case 0:
		cardIDStr = readLineWithPrompt("card id: ")
		front = readLineWithPrompt("front (_ keep): ")
		back = readLineWithPrompt("back (_ keep): ")
	case 1:
		cardIDStr = args[0]
		front = readLineWithPrompt("front (_ keep): ")
		back = readLineWithPrompt("back (_ keep): ")
	case 2:
		fmt.Printf("Expect 3 args for update card\nUsage:\n%s\n", update_usage)
		return
	default:
		cardIDStr = args[0]
		front = args[1]
		back = args[2]
	}

	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		fmt.Println("invalid card id")
		return
	}

	card, errApp := db.GetCardByID(a.db, a.deck.Id, cardID)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("card not found")
			return
		}
		log.Fatal(errApp)
		return
	}

	if front == "_" {
		front = card.Front
	}
	if back == "_" {
		back = card.Back
	}
	if front == card.Front && back == card.Back {
		fmt.Println("No changes")
		return
	}

	sourceText := card.SourceText
	errApp = db.UpdateCard(a.db, a.deck.Id, cardID, front, back, &sourceText, card.CreatedByAi)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("card not found")
			return
		}
		log.Fatal(errApp)
		return
	}

	fmt.Println("Card Updated")
}

func (a *CliApp) handleUpdateDeck(args []string) {
	// 0 arg: prompt for id, updated deckname
	// 1 arg: prompt for updated deckname
	// 2 / 2+ args: expect id, updated deckname, ignore others
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}

	var deckIDStr, newName string
	switch len(args) {
	case 0:
		deckIDStr = readLineWithPrompt("deck id: ")
		newName = readLineWithPrompt("deckname (_ keep): ")
	case 1:
		deckIDStr = args[0]
		newName = readLineWithPrompt("deckname (_ keep): ")
	default:
		deckIDStr = args[0]
		newName = args[1]
	}

	deckID, err := strconv.Atoi(deckIDStr)
	if err != nil {
		fmt.Println("invalid deck id")
		return
	}

	deck, errApp := db.GetDeckByDeckId(a.db, a.user.Id, deckID)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("deck not found")
			return
		}
		log.Fatal(errApp)
		return
	}

	if newName == "_" {
		newName = deck.Name
	}
	if newName == deck.Name {
		fmt.Println("No changes")
		return
	}

	errApp = db.UpdateDeckName(a.db, a.user.Id, deckID, newName)
	if errApp != nil {
		if errApp.Code == model.CodeNotFound {
			fmt.Println("deck not found")
			return
		}
		if errApp.Code == model.CodeDeckAlreadyExists {
			fmt.Println("deck already exists")
			return
		}
		log.Fatal(errApp)
		return
	}

	if a.deck != nil && a.deck.Id == deckID {
		a.deck.Name = newName
	}

	fmt.Println("Deck Updated")
}

func (a *CliApp) handleUpdateUser(args []string) {
	// updated value" '_' means no update this
	// args begin with me:
	// check if login
	// 0 remain args: prompt for updated value
	// <newname>/'_' <newpassword">/'_'

	// args not begin with me:
	// validate: prompt for currentpassword, check if validate
	// 0 arg: prompt for id, do validate, prompt for updated value
	// 1 arg: expect id, validate prompt for update value
	// 2 args: error
	// 3 / 3+ args: expect id, updated_value, do validate
	if len(args) == 0 {
		args = []string{"me"}
	}

	if args[0] == "me" {
		if err := a.validateUser(); err != nil {
			fmt.Println(err)
			return
		}

		var newUsername, newPassword string
		switch len(args) {
		case 1:
			newUsername = readLineWithPrompt("username (_ keep): ")
			newPassword = readLineWithPrompt("password (_ keep): ")
		case 2:
			newUsername = args[1]
			newPassword = readLineWithPrompt("password (_ keep): ")
		default:
			newUsername = args[1]
			newPassword = args[2]
		}

		if newUsername == "_" {
			newUsername = a.user.Username
		}
		if newPassword == "_" {
			newPassword = a.user.Password
		}
		if newUsername == a.user.Username && newPassword == a.user.Password {
			fmt.Println("No changes")
			return
		}

		errApp := db.UpdateUser(a.db, a.user.Id, newUsername, newPassword)
		if errApp != nil {
			if errApp.Code == model.CodeUserAlreadyExists {
				fmt.Println("user already exists")
				return
			}
			log.Fatal(errApp)
			return
		}

		a.user.Username = newUsername
		a.user.Password = newPassword
		fmt.Println("User Updated")
		return
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid user id")
		return
	}

	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}
	if a.user.Id != userID {
		fmt.Println("forbidden")
		return
	}

	currentPassword := readLineWithPrompt("current password: ")
	if currentPassword != a.user.Password {
		fmt.Println("password mismatch")
		return
	}

	var newUsername, newPassword string
	switch len(args) {
	case 1:
		newUsername = readLineWithPrompt("username (_ keep): ")
		newPassword = readLineWithPrompt("password (_ keep): ")
	case 2:
		newUsername = args[1]
		newPassword = readLineWithPrompt("password (_ keep): ")
	default:
		newUsername = args[1]
		newPassword = args[2]
	}

	if newUsername == "_" {
		newUsername = a.user.Username
	}
	if newPassword == "_" {
		newPassword = a.user.Password
	}
	if newUsername == a.user.Username && newPassword == a.user.Password {
		fmt.Println("No changes")
		return
	}

	errApp := db.UpdateUser(a.db, a.user.Id, newUsername, newPassword)
	if errApp != nil {
		if errApp.Code == model.CodeUserAlreadyExists {
			fmt.Println("user already exists")
			return
		}
		log.Fatal(errApp)
		return
	}

	a.user.Username = newUsername
	a.user.Password = newPassword
	fmt.Println("User Updated")
}

func (a *CliApp) handleUpdate(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", update_usage)
		return
	}
	target := args[0]
	switch target {
	case "card":
		a.handleUpdateCard(args[1:])
	case "deck":
		a.handleUpdateDeck(args[1:])
	case "user":
		a.handleUpdateUser(args[1:])
	default:
		fmt.Printf("Unknown update target: %q\nUsage:\n%s\n", target, update_usage)
	}
}
