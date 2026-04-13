package cli

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	. "github.com/lidchen/neuron_deck/backend/db"
)

// Default val for test
var username = "test"
var passwd = "testpasswd"
var deckname = "test"

func RunCliApp(db *sql.DB) {
	fmt.Println("🃏 Flashcard CLI — type \"help\" or \"exit\"")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		dispatch(scanner.Text())
	}
}

func testCilApp(db *sql.DB) {
	if username == "" || passwd == "" {
		username = readLineWithPrompt("username: ")
		passwd = readLineWithPrompt("passwd: ")
	}
	u, err := GetUserByUsername(db, username)
	if err != nil {
		if err.Code == "NOT_FOUND" {
			u, err = CreateUser(db, username, passwd)
			if err != nil {
				log.Fatal("Error at create user: ", err)
				return
			} else {
				printWithWrap(fmt.Sprintf("Create user %s", username))
			}
		} else {
			log.Fatal("Error at get user by name: ", err)
			return
		}
	} else {
		printWithWrap(fmt.Sprintf("Hello, %s", username))
	}

	if deckname == "" {
		deckname = readLineWithPrompt("Open or create deck: ")
	}
	deck, err := GetDeckByDeckName(db, u.Id, deckname)
	if err != nil {
		if err.Code == "NOT_FOUND" {
			deck, err = CreateDeck(db, u.Id, deckname)
			if err != nil {
				log.Fatal("Error at creating deck: ", err)
				return
			} else {
				printWithWrap(fmt.Sprintf("Created deck %s", deckname))
			}
		} else {
			log.Fatal("Error at get deck by name: ", err)
			return
		}
	} else {
		printWithWrap(fmt.Sprintf("Open deck %s", deckname))
	}

	cards, err := GetCards(db, deck.Id)
	if err != nil {
		log.Fatal("Error at getting cards: ", err)
		return
	}
	for _, c := range *cards {
		fmt.Println(c)
	}
}

func printWithWrap(s string) {
	l := len(s) + 4
	char := "="
	wrap := strings.Repeat(char, l)
	fmt.Println(wrap)
	fmt.Printf("| %s |\n", s)
	fmt.Println(wrap)
}

func handleShow(args []string) {

}

func handleDelete(args []string) {

}

func printHelp() {
	fmt.Printf("Commands:\n%s\n%s\n%s\n", create_usage, show_usage, delete_usage)
	fmt.Println(`
help                         Show this help message
exit                         Quit the program
	`)
}

// ---- Dispatch ----
func dispatch(input string) {
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
		handleCreate(args)
	case "show":
		handleShow(args)
	case "delete":
		handleDelete(args)
	case "help":
		printHelp()
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %q. Type \"help\" for usage.\n", cmd)
	}
}

func readLineWithPrompt(s string) string {
	fmt.Print(s)
	return readLine()
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
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
