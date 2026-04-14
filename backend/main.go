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
)

var requiredEnvKeys = []string{"DEEPSEEK_API_KEY", "DB_DSN", "URL", "DEBUG_MODE"}

func main() {
	initEnv()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	RunCardGeneratorTest()
}

func RunCardGeneratorTest() {
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
	client := &http.Client{}
	testmessage := []llmstream.Message{
		{
			Role:    "system",
			Content: prompt,
		},
		{
			Role:    "user",
			Content: "token 是模型用来表示自然语言文本的基本单位，也是我们的计费单元，可以直观的理解为“字”或“词”；通常 1 个中文词语、1 个英文单词、1 个数字或 1 个符号计为 1 个 token。 一般情况下模型中 token 和字数的换算比例大致如下： 1 个英文字符 ≈ 0.3 个 token。 1 个中文字符 ≈ 0.6 个 token。 但因为不同模型的分词不同，所以换算比例也存在差异，每一次实际处理 token 数量以模型返回为准，您可以从返回结果的 usage 中查看。 ",
		},
	}
	cardResponse, err := llmstream.GenerateCard(client, &testmessage)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(cardResponse.SourceText)
	for _, c := range cardResponse.CardData {
		fmt.Printf("front: %s", *c.Front)
		fmt.Printf("back: %s", *c.Back)
	}
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
	cliapp := cli.NewCliApp(db)
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
