package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/lidchen/neuron_deck/backend/model"
	. "github.com/lidchen/neuron_deck/backend/model"
)

func mapPgError(err error, conflictCode model.ErrorCode, conflictMessage string) *AppError {
	if pgErr, ok := err.(*pq.Error); ok {
		switch pgErr.Code {
		case "23505":
			if conflictCode != "" {
				return ErrConflict(conflictCode, conflictMessage)
			}
			return ErrConflict("CONFLICT", conflictMessage)
		case "23503":
			return ErrBadRequest("INVALID_REFERENCE_ID", conflictMessage)
		}
	}
	return ErrInternal(err)
}

func CreateUser(db *sql.DB, username, passwordHash string) (*User, *AppError) {
	var u User
	err := db.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash, created_at", username, passwordHash,
	).Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, ErrConflict(model.CodeUserAlreadyExists, "user already exists")
	}
	return &u, nil
}

func GetUserByID(db *sql.DB, id int) (*User, *AppError) {
	row := db.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE id=$1", id,
	)
	var u User
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "user not found")
		}
		return nil, ErrInternal(err)
	}
	return &u, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, *AppError) {
	row := db.QueryRow(
		"SELECT id, username, password_hash, created_at FROM users WHERE username=$1", username,
	)
	var u User
	err := row.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "user not found")
		}
		return nil, ErrInternal(err)
	}
	return &u, nil
}

func ListUsers(db *sql.DB) ([]User, *AppError) {
	rows, err := db.Query("SELECT id, username, password_hash, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, ErrInternal(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt); err != nil {
			return nil, ErrInternal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrInternal(err)
	}

	return users, nil
}

func UpdateUserPassword(db *sql.DB, id int, passwordHash string) *AppError {
	res, err := db.Exec("UPDATE users SET password_hash=$1 WHERE id=$2", passwordHash, id)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound("NOT_FOUND", "user not found")
	}
	return nil
}

func UpdateUser(db *sql.DB, id int, username, passwordHash string) *AppError {
	res, err := db.Exec(
		"UPDATE users SET username=$1, password_hash=$2 WHERE id=$3",
		username, passwordHash, id,
	)
	if err != nil {
		return mapPgError(err, "USER_ALREADY_EXISTS", "user already exists")
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "user not found")
	}
	return nil
}

func DeleteUser(db *sql.DB, id int) *AppError {
	res, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "user not found")
	}
	return nil
}

func CreateDeck(db *sql.DB, userID int, name string) (*Deck, *AppError) {
	var d Deck
	err := db.QueryRow(
		"INSERT INTO decks (user_id, name) VALUES ($1, $2) RETURNING *", userID, name,
	).Scan(&d.Id, &d.UserId, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, ErrConflict(model.CodeDeckAlreadyExists, "deck already exists")
	}
	return &d, nil
}

func GetDeckByDeckName(db *sql.DB, userID int, name string) (*Deck, *AppError) {
	row := db.QueryRow(
		"SELECT id, user_id, name, created_at, updated_at FROM decks WHERE user_id=$1 AND name=$2", userID, name,
	)
	var d Deck
	err := row.Scan(&d.Id, &d.UserId, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "deck not found")
		}
		return nil, ErrInternal(err)
	}
	return &d, nil
}

func GetDeckByDeckId(db *sql.DB, userID int, id int) (*Deck, *AppError) {
	row := db.QueryRow(
		"SELECT id, user_id, name, created_at, updated_at FROM decks WHERE user_id=$1 AND id=$2", userID, id,
	)
	var d Deck
	err := row.Scan(&d.Id, &d.UserId, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "deck not found")
		}
		return nil, ErrInternal(err)
	}
	return &d, nil
}

func ListDecksByUserID(db *sql.DB, userID int) ([]Deck, *AppError) {
	rows, err := db.Query(
		"SELECT id, user_id, name, created_at, updated_at FROM decks WHERE user_id=$1 ORDER BY id",
		userID,
	)
	if err != nil {
		return nil, ErrInternal(err)
	}
	defer rows.Close()

	decks := []Deck{}
	for rows.Next() {
		var d Deck
		if err := rows.Scan(&d.Id, &d.UserId, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, ErrInternal(err)
		}
		decks = append(decks, d)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrInternal(err)
	}

	return decks, nil
}

func UpdateDeckName(db *sql.DB, userID, id int, name string) *AppError {
	res, err := db.Exec(
		"UPDATE decks SET name=$1, updated_at=NOW() WHERE user_id=$2 AND id=$3",
		name, userID, id,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return ErrConflict(model.CodeDeckAlreadyExists, "deck already exists")
		}
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "deck not found")
	}
	return nil
}

func DeleteDeck(db *sql.DB, userID, id int) *AppError {
	res, err := db.Exec("DELETE FROM decks WHERE user_id=$1 AND id=$2", userID, id)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "deck not found")
	}
	return nil
}

func CreateCard(db *sql.DB, deckID int, front, back string, sourceText *string, createdByAI bool) *AppError {
	_, err := db.Exec(
		"INSERT INTO cards (deck_id, front, back, source_text, created_by_ai) VALUES ($1, $2, $3, $4, $5)",
		deckID, front, back, sourceText, createdByAI,
	)
	if err != nil {
		return ErrInternal(err)
	}
	return nil
}

func GetCardByID(db *sql.DB, deckID, id int) (*Card, *AppError) {
	row := db.QueryRow(
		"SELECT id, deck_id, front, back, source_text, created_by_ai, created_at, updated_at FROM cards WHERE deck_id=$1 AND id=$2",
		deckID, id,
	)

	var c Card
	var sourceText sql.NullString
	err := row.Scan(
		&c.Id, &c.DeckId, &c.Front, &c.Back, &sourceText, &c.CreatedByAi, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "card not found")
		}
		return nil, ErrInternal(err)
	}
	if sourceText.Valid {
		c.SourceText = &sourceText.String
	}

	return &c, nil
}

func GetCards(db *sql.DB, deckID int) ([]Card, *AppError) {
	rows, err := db.Query(
		"SELECT id, deck_id, front, back, source_text, created_by_ai, created_at, updated_at FROM cards WHERE deck_id=$1 ORDER BY id",
		deckID,
	)
	if err != nil {
		return nil, ErrInternal(err)
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var c Card
		var sourceText sql.NullString
		err := rows.Scan(
			&c.Id, &c.DeckId, &c.Front, &c.Back, &sourceText, &c.CreatedByAi, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, ErrInternal(err)
		}
		if sourceText.Valid {
			c.SourceText = &sourceText.String
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrInternal(err)
	}
	return cards, nil
}

func UpdateCard(db *sql.DB, deckID, id int, front, back string, sourceText *string, createdByAI bool) *AppError {
	res, err := db.Exec(
		"UPDATE cards SET front=$1, back=$2, source_text=$3, created_by_ai=$4, updated_at=NOW() WHERE deck_id=$5 AND id=$6",
		front, back, sourceText, createdByAI, deckID, id,
	)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "card not found")
	}
	return nil
}

func DeleteCard(db *sql.DB, deckID, id int) *AppError {
	res, err := db.Exec("DELETE FROM cards WHERE deck_id=$1 AND id=$2", deckID, id)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "card not found")
	}
	return nil
}

func CreateCardSrs(db *sql.DB, cardID int, nextReviewAt time.Time) *AppError {
	_, err := db.Exec(
		"INSERT INTO card_srs (card_id, next_review_at) VALUES ($1, $2)",
		cardID, nextReviewAt,
	)
	if err != nil {
		return mapPgError(err, "CARD_SRS_ALREADY_EXISTS", "card srs already exists")
	}
	return nil
}

func GetCardSrs(db *sql.DB, cardID int) (*CardSrs, *AppError) {
	row := db.QueryRow(
		"SELECT card_id, interval, ease_factor, repetitions, next_review_at, last_reviewed_at FROM card_srs WHERE card_id=$1",
		cardID,
	)
	var c CardSrs
	var lastReviewAt sql.NullTime
	err := row.Scan(&c.CardId, &c.Interval, &c.EaseFactor, &c.Repetitions, &c.NextReviewAt, &lastReviewAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound(model.CodeNotFound, "card srs not found")
		}
		return nil, ErrInternal(err)
	}
	if lastReviewAt.Valid {
		c.LastReviewAt = lastReviewAt.Time
	}
	return &c, nil
}

func UpdateCardSrs(db *sql.DB, cSrs *CardSrs) *AppError {
	res, err := db.Exec(
		"UPDATE card_srs SET interval=$1, ease_factor=$2, repetitions=$3, next_review_at=$4, last_reviewed_at=$5 WHERE card_id=$6",
		cSrs.Interval, cSrs.EaseFactor, cSrs.Repetitions, cSrs.NextReviewAt, cSrs.LastReviewAt, cSrs.CardId,
	)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "card srs not found")
	}
	return nil
}

/*
func UpdateCardSrs(db *sql.DB, cardID, interval, repetitions int, easeFactor float64, nextReviewAt time.Time, lastReviewAt *time.Time) *AppError {
	res, err := db.Exec(
		"UPDATE card_srs SET interval=$1, ease_factor=$2, repetitions=$3, next_review_at=$4, last_reviewed_at=$5 WHERE card_id=$6",
		interval, easeFactor, repetitions, nextReviewAt, lastReviewAt, cardID,
	)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "card srs not found")
	}
	return nil
}
*/

func DeleteCardSrs(db *sql.DB, cardID int) *AppError {
	res, err := db.Exec("DELETE FROM card_srs WHERE card_id=$1", cardID)
	if err != nil {
		return ErrInternal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return ErrInternal(err)
	}
	if affected == 0 {
		return ErrNotFound(model.CodeNotFound, "card srs not found")
	}
	return nil
}

func GetDueCardSrs(db *sql.DB, before time.Time) ([]CardSrs, *AppError) {
	rows, err := db.Query(
		"SELECT card_id, interval, ease_factor, repetitions, next_review_at, last_reviewed_at FROM card_srs WHERE next_review_at <= $1 ORDER BY next_review_at",
		before,
	)
	if err != nil {
		return nil, ErrInternal(err)
	}
	defer rows.Close()

	items := []CardSrs{}
	for rows.Next() {
		var c CardSrs
		var lastReviewAt sql.NullTime
		if err := rows.Scan(&c.CardId, &c.Interval, &c.EaseFactor, &c.Repetitions, &c.NextReviewAt, &lastReviewAt); err != nil {
			return nil, ErrInternal(err)
		}
		if lastReviewAt.Valid {
			c.LastReviewAt = lastReviewAt.Time
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, ErrInternal(err)
	}

	return items, nil
}

func GetCardToReview(db *sql.DB, before time.Time) (*Card, *AppError) {
	// Get one card with oldest next_review_at before time
	// if don't exist return nil, nil
	query :=
		`
SELECT cards.id, cards.deck_id, cards.front, cards.back,
			cards.source_text, cards.created_at, cards.updated_at
FROM cards 
LEFT JOIN card_srs ON cards.id = card_srs.card_id 
WHERE next_review_at <= $1 
ORDER BY next_review_at
LIMIT 1
`
	var c Card
	var sourceText sql.NullString
	err := db.QueryRow(query, before).Scan(&c.Id, &c.DeckId, &c.Front, &c.Back, &sourceText, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // expected: nothing to review
	}
	if err != nil {
		return nil, ErrInternal(err)
	}
	if sourceText.Valid {
		c.SourceText = &sourceText.String
	}
	return &c, nil
}

func PingDB(db *sql.DB) *AppError {
	if err := db.Ping(); err != nil {
		return ErrInternal(fmt.Errorf("database ping failed: %w", err))
	}
	return nil
}
