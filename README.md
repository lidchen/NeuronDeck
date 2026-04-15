# NeuronDeck

NeuronDeck is an interactive flashcard CLI written in Go. It stores users, decks, cards, and spaced-repetition state in Postgres, and can also call an LLM to generate cards from source text or generate hints during review.

## Features

- User login and account creation
- Deck management per user
- Card CRUD inside the currently opened deck
- Spaced repetition review flow with 1-5 ratings
- LLM-assisted card generation from source text
- LLM hint generation while reviewing cards

## Requirements

- Go 1.26.1
- PostgreSQL
- A DeepSeek-compatible chat API endpoint

## Configuration

The backend loads environment variables from a `.env` file one level above the backend module, so the usual layout is a root-level `.env` file next to `go.work`.

Required variables:

- `DEEPSEEK_API_KEY`
- `DB_DSN`
- `URL`

Optional debug variables:

- `DEBUG_MODE`  
- `DEBUG_AUTO_LOGIN`
- `DEBUG_USERNAME`
- `DEBUG_PASSWORD`
- `DEBUG_DECKNAME`

Example `.env`:

```env
DEEPSEEK_API_KEY=your-api-key
DB_DSN=postgres://postgres:postgres@localhost:5432/neuron_deck?sslmode=disable
URL=https://api.deepseek.com/chat/completions
DEBUG_MODE=1

# Optional helpers for local development
DEBUG_AUTO_LOGIN=1
DEBUG_USERNAME=demo
DEBUG_PASSWORD=demo
DEBUG_DECKNAME=core
```

## Database Schema

The app expects these tables:

- `users`
- `decks`
- `cards`
- `card_srs`

The schema in `note.md` shows the current structure and the trigger used to create SRS records automatically when a card is inserted.

## Run

From the `backend` directory:

```bash
go run .
```

The program starts an interactive CLI session. Type `help` to see the available commands.

## Commands

Authentication and navigation:

- `login <username>` or `login <username> <password>`
- `create user <name> <password>`
- `create deck <name>`
- `opendeck <deckname>`

Card workflow:

- `create card <front> <back>`
- `show cards`
- `show srs`
- `update card <id> <front|_> <back|_>`
- `delete card <id>`

Deck and user management:

- `show decks`
- `show user`
- `update deck <id> <name|_>`
- `update user me [name|_] [password|_]`
- `update user <id> [name|_] [password|_]`
- `delete deck <id>`

LLM and review:

- `llm gencard <source_text>`
- `review`

During review, press `Enter` to reveal the answer, `h` to request an LLM hint, and `q` to quit.

## Examples

```text
login alice secret
create deck golang
opendeck golang
create card "What is a goroutine?" "A lightweight concurrent function execution unit."
show cards
llm gencard "Go channels are typed conduits for communication between goroutines."
review
```

## Test

From the `backend` directory:

```bash
go test ./...
```

There is a focused test suite for the learning-phase spaced-repetition logic in `backend/test/learningphase_test.go`.
