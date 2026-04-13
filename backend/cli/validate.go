package cli

import "errors"

func (a *CliApp) validateUser() error {
	if a.user == nil {
		return errors.New("not login, please login first")
	}
	return nil
}

func (a *CliApp) validateDeck() error {
	if a.deck == nil {
		return errors.New("deck is not opened, please open deck first")
	}
	return nil
}
