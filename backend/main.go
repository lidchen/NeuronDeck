package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/lidchen/neuron_deck/backend/llmstream"
)

func main() {
	initEnv()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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
	c := llmstream.NewConversation("your are a helpful assistant")
	client := &http.Client{}
	c.RunInteractiveChat(client)
}

func initEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading env file")
		return
	}
	_ = os.Getenv("DEEPSEEK_API_KEY")
	_ = os.Getenv("DB_DSN")
	_ = os.Getenv("URL")
}
