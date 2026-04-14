package cli

import (
	"fmt"
	"os"
	"strings"
)

// ---- Dispatch ----
func (a *CliApp) dispatch(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Simple tokenizer: respects quoted strings like "hello world"
	tokens := tokenize(input)
	if len(tokens) == 0 {
		return
	}

	cmd, args := tokens[0], tokens[1:]

	switch cmd {
	case "create":
		a.handleCreate(args)
	case "show":
		a.handleShow(args)
	case "update":
		a.handleUpdate(args)
	case "delete":
		a.handleDelete(args)
	case "login":
		a.handleLogin(args)
	case "opendeck":
		a.handleOpenDeck(args)
	case "review":
		a.handleReview(args)
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %q. Type \"help\" for usage.\n", cmd)
	}
}

// tokenize splits by whitespace but keeps quoted strings together.
// e.g. `create card "what is Go" "a language"` → ["create","card","what is Go","a language"]
func tokenize(s string) []string {
	var tokens []string
	var current strings.Builder
	inQuote := false

	for _, ch := range s {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ' ' && !inQuote:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}
