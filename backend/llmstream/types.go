package llmstream

import (
	"errors"
)

type Delta struct {
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

type ChunkResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// Message represents a single turn in the conversation
type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

func GetSourceData(messages *[]Message) (*string, error) {
	for _, m := range *messages {
		if m.Role == "user" {
			return &m.Content, nil
		}
	}
	return nil, errors.New("cant found user role")
}

// ConversationManager holds the full history
type ConversationManager struct {
	History      []Message
	SystemPrompt string
}

type CardData struct {
	Front *string `json:"front"`
	Back  *string `json:"back"`
}

type CardResponse struct {
	CardData   []CardData `json:"cards"`
	SourceText *string
}
