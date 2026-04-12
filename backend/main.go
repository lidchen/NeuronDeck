package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	TOKEN  string
	DB_DSN string
	url    string = "https://api.deepseek.com/chat/completions"
	method string = "POST"
)

func main() {
	initEnv()
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatalf("open db: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
		return
	}

}

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading env file")
		return
	}
	TOKEN = os.Getenv("DEEPSEEK_API_KEY")
	DB_DSN = os.Getenv("DB_DSN")
}
