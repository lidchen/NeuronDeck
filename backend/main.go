package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lidchen/neuron_deck/backend/cli"
	"github.com/lidchen/neuron_deck/backend/llmstream"
	"github.com/lidchen/neuron_deck/backend/model"
)

var requiredEnvKeys = []string{"DEEPSEEK_API_KEY", "DB_DSN", "URL"}

func main() {
	initEnv()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	RunCliApp()
}

func RunCliApp() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatalf("open db: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
		return
	}
	var apperr *model.AppError
	cliapp, apperr := cli.NewCliApp(db)
	if apperr != nil {
		log.Fatalf("error at create cliapp: %s", apperr.Message)
	}
	cli.RunCliApp(cliapp)
}

func RunChat() {
	c := llmstream.NewConversation("your are a helpful assistant")
	client := &http.Client{}
	c.RunInteractiveChat(client)
}

func initEnv() {
	if err := initEnvFromPath("../.env"); err != nil {
		log.Fatal(err)
		return
	}
}

func initEnvFromPath(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("load env file: %w", err)
	}

	for _, key := range requiredEnvKeys {
		val, ok := os.LookupEnv(key)
		if !ok || strings.TrimSpace(val) == "" {
			return fmt.Errorf("missing required env var: %s", key)
		}
	}

	return nil
}
