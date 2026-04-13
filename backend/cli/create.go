package cli

import "fmt"

func handleCreateCard(args []string) {

}

func handleCreateDeck(args []string) {

}

func handleCreateUser(args []string) {

}

func handleCreate(args []string) {
	if len(args) == 0 {
		fmt.Printf("Expect at least one parameter\nUsage:\n%s\n", create_usage)
		return
	}

	target := (args)[0]
	switch target {
	case "card":
		handleCreateCard(args[1:])
	case "deck":
		handleCreateDeck(args[1:])
	case "user":
		handleCreateUser(args[1:])
	default:
		fmt.Printf("Unknown create target: %q\nUsage:\n%s\n", target, create_usage)
	}
}
