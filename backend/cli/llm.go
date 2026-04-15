package cli

import (
	"fmt"
	"log"

	"github.com/lidchen/neuron_deck/backend/db"
	"github.com/lidchen/neuron_deck/backend/llmstream"
)

func (a *CliApp) handleLLMGenCard(args []string) {
	// no args: prompt for source text
	// 1 / 1+ args: ignore

	if a.user == nil {
		fmt.Println(ErrNoLogin.Error())
		return
	}
	if a.deck == nil {
		fmt.Println(ErrNoDeckOpen.Error())
		return
	}
	source_text := readLineWithPrompt("source_text: ")

	if len(source_text) == 0 {
		fmt.Println("Source text should be non empty")
		return
	}

	cardResponse, err := llmstream.GenerateCard(a.client, &source_text)
	if err != nil {
		log.Fatal(err.Message)
		return
	}
	a.cardResponseInterface(cardResponse)
}

func (a *CliApp) cardResponseInterface(cardResponse *llmstream.CardResponse) {
	for _, cData := range cardResponse.CardData {
		fmt.Printf("front: %s\n", *cData.Front)
		fmt.Printf("back: %s\n", *cData.Back)
	outer:
		for {
			i := readLineWithPrompt("y: keep\tn:discard\ttype exit to quit\n:")
			switch i {
			case "y":
				{
					err := db.CreateCard(a.db, a.deck.Id, cData.Front, cData.Back, cardResponse.SourceText, true)
					if err != nil {
						log.Printf("Error at createcard: %s", err.Message)
					}
					break outer
				}
			case "n":
				break outer
			case "quit":
				return
			default:
				fmt.Printf("unknown command: %q, try again\n", i)
			}
		}
	}
}

func (a *CliApp) handleLLM(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", llm_usage)
		return
	}

	target := (args)[0]
	switch target {
	case "gencard":
		a.handleLLMGenCard(args[1:])
	default:
		fmt.Printf("Unknown llm target: %q\nUsage:\n%s\n", target, llm_usage)
	}
}
