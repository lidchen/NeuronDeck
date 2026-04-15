package cli

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
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
		done, err := reviewNextCard(a)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Message)
		}
		if done {
			return
		}
	}
}

func reviewNextCard(a *CliApp) (done bool, apperr *model.AppError) {
	c, err := db.GetCardToReview(a.db, time.Now())
	if err != nil {
		return false, err
	}
	if c == nil {
		fmt.Println("no card to review")
		return true, nil
	}

	should_exit, q, err := a.cardReviewInterface(c)
	if should_exit || err != nil {
		return true, err
	}

	cSrs, err := db.GetCardSrs(a.db, c.Id)
	if err != nil {
		return true, err
	}
	err = a.srs.Review(cSrs, *q)
	if err != nil {
		return true, err
	}
	err = db.UpdateCardSrs(a.db, cSrs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Message)
	}
	return false, nil
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
		line := readLineWithPrompt("rank from 1-5: ")
		if line == "exit" || line == "quit" {
			should_exit = true
			return should_exit, nil, nil
		}
		q, convErr := strconv.Atoi(line)
		if convErr != nil || q < 1 || q > 5 {
			fmt.Println("please enter number between 1 and 5")
			continue
		}
		return should_exit, &q, nil
	}
}
