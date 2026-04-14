package llmstream

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func streamChatCompletionChunks(ctx context.Context, client *http.Client, payload *strings.Reader) (<-chan ChunkResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, os.Getenv("URL"), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("DEEPSEEK_API_KEY"))

	out := make(chan ChunkResponse)

	go func() {
		defer close(out)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("stream request failed:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("stream request failed: status %s: %s", resp.Status, strings.TrimSpace(string(body)))
			return
		}

		if err := scanChunkStream(ctx, resp, out); err != nil {
			log.Println("stream scan failed:", err)
		}
	}()

	return out, nil
}
func scanChunkStream(ctx context.Context, resp *http.Response, out chan<- ChunkResponse) error {
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return nil
		}

		var chunk ChunkResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- chunk:
		}
	}

	return scanner.Err()
}

func genStreamPayload(m *[]Message) (*strings.Reader, error) {
	messagesJson, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(fmt.Sprintf(`{
		"messages": %s,
		"model": "deepseek-chat",
		"thinking": {
			"type": "disabled"
		},
		"frequency_penalty": 0,
		"max_tokens": 4096,
		"presence_penalty": 0,
		"response_format": {
			"type": "text"
		},
		"stop": null,
		"stream": true,
		"stream_options": null,
		"temperature": 1,
		"top_p": 1,
		"tools": null,
		"tool_choice": "none",
		"logprobs": false,
		"top_logprobs": null
	}`, string(messagesJson)))
	return payload, nil
}

func genStreamJsonPayload(m *[]Message) (*strings.Reader, error) {
	messageJson, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(fmt.Sprintf(`{
		"messages": %s,
		"model": "deepseek-chat",
		"thinking": {
			"type": "disabled"
		},
		"frequency_penalty": 0,
		"max_tokens": 2048,
		"presence_penalty": 0,
		"response_format": {
			"type": "json_object"
		},
		"stop": null,
		"stream": true,
		"stream_options": null,
		"temperature": 1,
		"top_p": 1,
		"tools": null,
		"tool_choice": "none",
		"logprobs": false,
		"top_logprobs": null
	}`, string(messageJson)))
	return payload, nil
}
