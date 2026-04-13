package cli

import (
	"fmt"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
)

func (a *CliApp) handleOpenDeck(args []string) {
	// no arg: prompt deck name and open
	// 1 arg: open deck
	// 1+ args: ignore other args
	var deckname string
	if len(args) == 0 {
		deckname = readLineWithPrompt("deckname: ")
	} else {
		deckname = args[0]
	}
	if err := a.validateUser(); err != nil {
		fmt.Println(err)
		return
	}
	var err *model.AppError
	a.deck, err = db.GetDeckByDeckName(a.db, a.user.Id, deckname)
	if err != nil {
		if err.Code == "NOT_FOUND" {
			fmt.Println("deck not found, please create one first")
			return
		}
		fmt.Println(err)
		return
	}
}
