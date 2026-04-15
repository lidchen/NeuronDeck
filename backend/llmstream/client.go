package llmstream

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lidchen/neuron_deck/backend/model"
)

func GenerateCard(c *http.Client, content *string) (*CardResponse, *model.AppError) {
	prompt := `
You are a spaced repetition flashcard generator. 
Given source text, extract the most important concepts and generate flashcards. 
Rules: 
- Front: a clear, specific question or cloze prompt 
- Back: concise answer (1-3 sentences max) 
- Prefer atomic cards (one fact per card) 
- Output ONLY valid JSON, no markdown fences. 
Output format:{\"cards\": [{\"front\": \"...\", \"back\": \"...\"}]}
`
	m := []Message{
		{
			Role:    "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: *content,
		},
	}
	cardResponse, err := generateCard(c, &m)
	if err != nil {
		return nil, model.ErrInternal(err)
	}
	return cardResponse, nil
}

// TODO:
// specify language
// specify max cards generated
// custom parser, parse each card once finished
func generateCard(client *http.Client, message *[]Message) (*CardResponse, error) {
	var c CardResponse
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	payload, err := genStreamJsonPayload(message)
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
	fmt.Println()
	decoder := json.NewDecoder(strings.NewReader(fullResponse.String()))

	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}

	if c.SourceText, err = GetSourceData(message); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *ConversationManager) RunInteractiveChat(client *http.Client) {
	for {
		fmt.Print(":")
		input := readLine()
		if isExitCommand(input) {
			fmt.Println("EXIT")
			return
		}

		c.AddUser(input)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		payload, err := genStreamPayload(&c.History)
		if err != nil {
			log.Fatal("stream chat failed:", err)
			return
		}
		chunks, err := streamChatCompletionChunks(ctx, client, payload)

		if err != nil {
			log.Fatal("stream chat failed:", err)
			return
		}

		var fullResponse strings.Builder
		for chunk := range chunks {
			if len(chunk.Choices) == 0 {
				continue
			}
			fmt.Print(chunk.Choices[0].Delta.Content)
			fullResponse.WriteString(chunk.Choices[0].Delta.Content)
		}
		c.AddAssistant(fullResponse.String())
		fmt.Println()
		cancel()
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func isExitCommand(input string) bool {
	return input == "" || input == "exit" || input == "q" || input == "quit"
}
