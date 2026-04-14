package cli

import "fmt"

func (a *CliApp) handleReview(args []string) {
	// Check login
	// 0 arg: prompt for deckname if not opened
	// 1 / 1+ arg: except deckname, ignore others
	if a.user == nil {
		fmt.Println(ErrNoLogin.Error())
	}
}
