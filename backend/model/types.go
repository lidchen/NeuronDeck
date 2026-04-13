package model

import "time"

type User struct {
	Id        int
	Username  string
	Password  string
	CreatedAt time.Time
}

type Deck struct {
	Id        int
	UserId    int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Card struct {
	Id          int
	DeckId      int
	Front       string
	Back        string
	SourceText  string
	CreatedByAi bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CardSrs struct {
	CardId       int
	Interval     int
	EaseFactor   float32
	Repetitions  int
	NextReviewAt time.Time
	LastReviewAt time.Time
}
