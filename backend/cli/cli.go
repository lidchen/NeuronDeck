package cli

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/model"
	"github.com/lidchen/neuron_deck/backend/srs"
)

type CliApp struct {
	db   *sql.DB
	deck *model.Deck
	user *model.User
	srs  *srs.SRSService
}

func NewCliApp(database *sql.DB) (*CliApp, *model.AppError) {
	var c CliApp = CliApp{}
	c.db = database
	debugAutoLogin := os.Getenv("DEBUG_AUTO_LOGIN")
	debugUsername := os.Getenv("DEBUG_USERNAME")
	debugPassword := os.Getenv("DEBUG_PASSWORD")
	debugDeckname := os.Getenv("DEBUG_DECKNAME")
	debugSrs := os.Getenv("DEBUG_SRS")
	if debugSrs == "1" {
		c.srs = srs.NewSRSService(&srs.MockClock{})
	} else {
		c.srs = srs.NewSRSService(&srs.RealClock{})
	}
	if debugAutoLogin == "1" {
		if debugUsername != "" && debugPassword != "" {
			u, success, err := db.Login(database, debugUsername, debugPassword)
			if err != nil {
				return nil, err
			}
			if !success {
				return nil, model.ErrNotFound(model.CodeNotFound, "auto login not success")
			}
			c.user = u
		}
		if debugDeckname != "" {
			d, err := db.GetDeckByDeckName(database, c.user.Id, debugDeckname)
			if err != nil {
				return nil, err
			}
			c.deck = d
		}
	}
	return &c, nil
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
	fmt.Printf("Commands:\n%s\n%s\n%s\n%s\n%s\n",
		create_usage, show_usage, update_usage, delete_usage, other_usage)
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
