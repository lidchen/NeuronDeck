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
update card                  Update a card in the open deck
  <id> <front|_> <back|_>         					
update deck                  Rename one of your decks
  <id> <name|_>                   
update user me               Update the current user
  [name|_] [password|_]        
update user                  Update the user
  <id> [name|_] [password|_]      
`

var delete_usage string = `
delete card <id>             Delete a card by ID
delete deck <id>             Delete a deck by ID
`

var other_usage string = `
login <username>             User login
opendeck <deckname>          Open deck
help                         Show this help message
exit                         Quit the program
`
