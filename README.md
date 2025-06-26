# Telegram Interactive Survey Bot in Go

Interactive Telegram bot for conducting surveys with rich features:

- Tree-based question logic with transitions
- Support for multiple messages and images
- Text input and inline buttons
- Automatic transitions with configurable delays
- User geolocation collection
- Message personalization (name substitution)
- Response summary
- External links

## Project Structure

```text
tlgbot/
├── cmd/                    # Main applications
│   └── telegram-bot/       # Main bot application
│       └── main.go
├── internal/               # Internal packages (not exported)
│   ├── bot/                # Telegram API logic
│   ├── config/             # Configuration handling
│   ├── handlers/           # Request handlers
│   ├── models/             # Data models
│   └── services/           # Business logic and services
├── configs/                # Configuration files
│   ├── config.example.json
│   └── questions.json
├── assets/                 # Static resources (images)
├── docs/                   # Documentation
│   ├── README.md           # Detailed documentation
│   ├── SECURITY.md         # Security policy
│   └── TESTING.md          # Testing documentation
├── go.mod                  # Go module
├── go.sum                  # Module dependencies
├── Makefile                # Build commands
└── .gitignore              # Git exclusions
```

## Installation and Running

### Requirements

- Go 1.21 or higher
- Telegram bot token

### Environment Variables

Create environment variables or use a `.env` file:

```bash
export TELEGRAM_TOKEN="your_telegram_bot_token"
export GOOGLE_CREDS="google-credentials.json"  # optional
export SHEET_ID="your_google_sheet_id"         # optional
export DELAY_MS="700"                           # optional
export START_QUESTION_ID="start"               # optional
```

### Build and Run

```bash
# Install dependencies
make deps

# Build project (creates telegram-bot executable)
make build

# Run project (requires TELEGRAM_TOKEN environment variable)
make run

# Alternative run after build
./telegram-bot

# Show all available commands
make help
```

## Usage

1. Configure environment variables with bot token
2. Edit the `configs/questions.json` file according to your needs
3. Run the bot with one of the commands:
   - `make run` - to run from source code
   - `make build && ./telegram-bot` - to run compiled binary

**Important:** When running without the `TELEGRAM_TOKEN` environment variable set, the application will exit with an error.

## Development

### Package Structure

- **cmd/telegram-bot/** - application entry point
- **internal/models/** - data structures (Config, Question, Option)
- **internal/config/** - configuration and questions loading
- **internal/bot/** - Telegram API logic
- **internal/handlers/** - HTTP/Telegram request handlers
- **internal/services/** - business logic (state and question managers)

### Adding New Features

1. Add data models to `internal/models/`
2. Place business logic in corresponding `internal/` packages
3. Use functional programming approach

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with race condition check
make test-race

# Code coverage with report
make test-coverage

# Run benchmarks
make bench

# Clean testing artifacts
make test-clean
```

### Code Coverage

Current code coverage:

- **internal/models**: 100%
- **internal/services**: 100%
- **internal/config**: 90.6%
- **Total coverage**: 31.6%

For detailed testing documentation see [docs/TESTING.md](docs/TESTING.md)

## Deployment

```bash
# Build for Linux
make build-linux

# Clean build artifacts
make clean
```

## Documentation

More detailed documentation is available in the `docs/` folder:

- [README.md](docs/README.md) - Detailed feature description and configuration
- [TESTING.md](docs/TESTING.md) - Testing guide
- [SECURITY.md](docs/SECURITY.md) - Security policy

## Dependencies

- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API

For current list of dependencies see `go.mod` file.