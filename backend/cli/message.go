package cli

var create_usage string = ` 
create card <front> <back>   Create a new flashcard
create deck <name>           Create a new deck
create user <name>           Set the current user
`

var show_usage string = `
show cards                   List all cards
show decks                   List all decks
show user                    Show current user
`

var update_usage string = `
update card <id> <front|_> <back|_>         Update a card in the open deck
update deck <id> <name|_>                   Rename one of your decks
update user me [name|_] [password|_]        Update the current user
update user <id> [name|_] [password|_]      Update the current user after password check
`

var delete_usage string = `
delete card <id>             Delete a card by ID
delete deck <id>             Delete a deck by ID
`

var other_usage string = `
login <username>             User login
deck <deckname>              Open deck
help                         Show this help message
exit                         Quit the program
`
