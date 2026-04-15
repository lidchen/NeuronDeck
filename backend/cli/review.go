package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/llmstream"
	"github.com/lidchen/neuron_deck/backend/model"
	"golang.org/x/term"
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
	var shouldExit bool = false
	// Validate card
	if c == nil {
		return shouldExit, nil, model.ErrInternal(fmt.Errorf("Not validate card"))
	}
	// Print front
	fmt.Printf("Front: %s\n", c.Front)
	// Get input

	shouldExit = a.promptHint(c)
	if shouldExit {
		return shouldExit, nil, nil
	}

	// Print back
	fmt.Printf("Back: %s\n", c.Back)
	// Handle input, get q

	for {
		line := readLineWithPrompt("rank from 1-5: ")
		if line == "exit" || line == "quit" {
			shouldExit = true
			return shouldExit, nil, nil
		}
		q, convErr := strconv.Atoi(line)
		if convErr != nil || q < 1 || q > 5 {
			fmt.Println("please enter number between 1 and 5")
			continue
		}
		return shouldExit, &q, nil
	}
}

// TODO:
// block input when hint is not finished
func (a *CliApp) promptHint(c *model.Card) bool {
	inputHint := "Press <h> get llm hint, <Enter> to show answer, <q> to quit"
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)

	for {
		fmt.Println(inputHint)
		fmt.Print("\r")
		_, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}

		b := buf[0]

		switch b {
		case '\r', '\n':
			return false
		case 'h':
			fmt.Print("Hint: ")
			_, err := llmstream.GenHint(a.client, c)
			if err != nil {
				log.Printf("Error at GenHint: %s", err.Message)
				return false
			}
			fmt.Print("\n\r")
			continue
		case 'q':
			return true
		default:
			continue
		}
	}
}
