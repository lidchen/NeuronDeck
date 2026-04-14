package cli

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
	"github.com/lidchen/neuron_deck/backend/srs"
)

func (a *CliApp) handleReview(args []string) {
	// Check login
	// 0 arg: prompt for deckname if not opened
	// 1 / 1+ arg: except deckname, ignore others
	if a.user == nil {
		fmt.Println(ErrNoLogin.Error())
		return
	}
	if a.deck == nil {
		fmt.Println(ErrNoDeckOpen.Error())
		return
	}

	// Review mode logic:
	// show due card, handle input(q), update cardsrs
	// type exit / quit to exit review mode
	for {
		c, err := db.GetCardToReview(a.db, time.Now())
		if err != nil {
			log.Fatal(err.Message)
			return
		}
		if c == nil {
			fmt.Println("no card to review")
			return
		}
		should_exit, q, err := a.cardReviewInterface(c)
		if should_exit {
			break
		}
		if err != nil {
			log.Fatal(err.Message)
			return
		}
		c_srs, err := db.GetCardSrs(a.db, c.Id)
		if err != nil {
			log.Fatal(err.Message)
			return
		}
		err = srs.Review(c_srs, *q)
		if err != nil {
			log.Fatal(err.Message)
		}
	}
}

func (a *CliApp) cardReviewInterface(c *model.Card) (bool, *int, *model.AppError) {
	var should_exit bool = false
	// Validate card
	if c == nil {
		return should_exit, nil, model.ErrInternal(fmt.Errorf("Not validate card"))
	}

	// Print front
	fmt.Printf("Front: %s\n", c.Front)
	// Get input
	readLineWithPrompt("Press <Enter> to show answer")
	// Print back
	fmt.Printf("Back: %s\n", c.Back)
	// Handle input, get q

	for {
		line := readLineWithPrompt("rank from 1-5:")
		if line == "exit" || line == "quit" {
			should_exit = true
			return should_exit, nil, nil
		}
		q, err := strconv.Atoi(line)
		if err != nil {
			return should_exit, nil, model.ErrInternal(err)
		}
		if q < 1 || q > 5 {
			fmt.Println("q should be 1-5")
		} else {
			return should_exit, &q, nil
		}
	}
}
