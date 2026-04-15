package db

import (
	"database/sql"

	"github.com/lidchen/neuron_deck/backend/model"
)

func Login(database *sql.DB, username string, password string) (user *model.User, success bool, err *model.AppError) {
	valid_user, err := GetUserByUsername(database, username)
	if err != nil {
		if err.Code == model.CodeNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	if valid_user.Password != password {
		return nil, false, nil
	}
	return valid_user, true, nil
}
