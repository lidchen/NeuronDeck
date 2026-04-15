package llmstream

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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
