package cli

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/lidchen/neuron_deck/backend/model"
)

type CliApp struct {
	db   *sql.DB
	deck *model.Deck
	user *model.User
}

func NewCliApp(db *sql.DB) *CliApp {
	return &CliApp{db: db}
}

func RunCliApp(a *CliApp) {
	fmt.Println("🃏 Flashcard CLI — type \"help\" or \"exit\"")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		a.dispatch(scanner.Text())
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

func printHelp() {
	fmt.Printf("Commands:\n%s\n%s\n%s\n%s\n", create_usage, show_usage, delete_usage, other_usage)
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

/*

// Default val for test
var username = "test"
var passwd = "testpasswd"
var deckname = "test"

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
*/
