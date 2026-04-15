package llmstream

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lidchen/neuron_deck/backend/model"
)

func GenHint(client *http.Client, c *model.Card) (*string, *model.AppError) {
	cContent := fmt.Sprintf("front:%s, back:%s", c.Front, c.Back)
	prompt := `
You are a flashcard review assistant.
Given a card's front and back, generate a short hint that helps the learner recall the answer.
Rules:
- Keep the hint brief, specific, and useful.
- Do not repeat the full answer or reveal it verbatim.
- Focus on a keyword, concept, association, or framing clue.
- Output plain text only, with no markdown, bullets, or extra commentary.
`
	m := []Message{
		{
			Role:    "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: cContent,
		},
	}
	hint, err := genHint(client, &m)
	if err != nil {
		return nil, model.ErrInternal(err)
	}
	return hint, nil
}

func genHint(client *http.Client, message *[]Message) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	payload, err := genStreamPayload(message)
	if err != nil {
		return nil, err
	}
	chunks, err := streamChatCompletionChunks(ctx, client, payload)
	if err != nil {
		return nil, err
	}
	var fullResponse strings.Builder
	for chunk := range chunks {
		if len(chunk.Choices) == 0 {
			continue
		}
		fmt.Print(chunk.Choices[0].Delta.Content)
		fullResponse.WriteString(chunk.Choices[0].Delta.Content)
	}
	hint := fullResponse.String()
	return &hint, nil
}
