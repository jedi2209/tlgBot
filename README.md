# Telegram Interactive Survey Bot in Go

[![CI](https://github.com/jedi2209/tlgBot/workflows/CI/badge.svg)](https://github.com/jedi2209/tlgBot/actions)
[![codecov](https://codecov.io/gh/jedi2209/tlgBot/branch/main/graph/badge.svg)](https://codecov.io/gh/jedi2209/tlgBot)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedi2209/tlgBot)](https://goreportcard.com/report/github.com/jedi2209/tlgBot)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Interactive Telegram bot for conducting surveys with rich features:

- Tree-based question logic with transitions
- Support for multiple messages and images
- Text input and inline buttons
- Automatic transitions with configurable delays
- User geolocation collection
- Message personalization (name substitution)
- Response summary
- External links
- Flexible question configuration system
- Support for custom questions outside the repository

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
│   ├── config.example.json # Configuration example
│   ├── questions.json      # Demo questions
│   └── questions.example.json # Questions example
├── assets/                 # Static resources (images)
├── go.mod                  # Go module
├── go.sum                  # Module dependencies
├── Makefile                # Build commands
├── MIGRATION_GUIDE.md      # Migration guide for custom questions
├── QUESTIONS_SETUP.md      # Guide for setting up custom questions
├── README.md               # This file
├── SECURITY.md             # Security policy
└── .gitignore              # Git exclusions
```

## Installation and Running

### Requirements

- Go 1.23 or higher
- Telegram bot token

### Configuration

You can configure the bot using either environment variables or a JSON configuration file.

#### Option 1: Environment Variables

Create environment variables:

```bash
export TELEGRAM_TOKEN="your_telegram_bot_token"
export QUESTIONS_FILE_PATH="path/to/your/questions.json"  # optional
export START_QUESTION_ID="start"                         # optional
export DELAY_MS="700"                                     # optional
```

#### Option 2: Configuration File

Create a configuration file based on `configs/config.example.json`:

```json
{
  "telegram_token": "YOUR_TELEGRAM_BOT_TOKEN",
  "questions_file_path": "configs/questions.json",
  "delay_ms": 700,
  "start_question_id": "start"
}
```

### Build and Run

```bash
# Install dependencies
make deps

# Build project (creates telegram-bot executable)
make build

# Run project (requires TELEGRAM_TOKEN environment variable)
make run

# Alternative: run with configuration file
./telegram-bot config.json

# Alternative: run compiled binary with environment variables
./telegram-bot

# Show all available commands
make help
```

## Questions Configuration

### Using Demo Questions

By default, the bot uses demo questions from `configs/questions.json`. This is perfect for trying out the bot.

### Using Custom Questions

For production use, you should create your own questions file outside the repository:

1. **Create your questions file:**

   ```bash
   cp configs/questions.example.json my-questions.json
   ```

2. **Configure the path:**

   ```bash
   export QUESTIONS_FILE_PATH="my-questions.json"
   ```

3. **Run the bot:**

   ```bash
   ./telegram-bot
   ```

For detailed instructions, see [QUESTIONS_SETUP.md](QUESTIONS_SETUP.md).

### Migration from Old Versions

If you have existing questions in the repository, see [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) for migration instructions.

## Usage

1. Configure your Telegram bot token (see Configuration section)
2. Set up your questions (use demo questions or create custom ones)
3. Run the bot with one of the methods described above

**Important:** The application requires a valid `TELEGRAM_TOKEN` to function.

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

### Local Testing

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

### Continuous Integration

This project uses GitHub Actions for automated testing and quality checks:

#### CI Pipeline Features

- **Multi-version testing**: Tests run on Go 1.23
- **Comprehensive test suite**: Unit tests, race condition detection, and coverage analysis
- **Code quality**: Automated linting with golangci-lint
- **Build validation**: Cross-platform builds (Linux and macOS)
- **Benchmarking**: Performance testing on main branch updates
- **Dependency management**: Automated dependency updates via Dependabot

#### Workflow Triggers

- **Push to main/develop**: Full test suite, build, and benchmarks
- **Pull requests**: All tests and checks (excluding benchmarks)
- **Weekly dependency updates**: Automated via Dependabot

#### Coverage Reporting

Coverage reports are automatically uploaded to Codecov on every CI run.

### Code Coverage

Current code coverage:

- **internal/models**: 100%
- **internal/services**: 100%  
- **internal/config**: 90.6%
- **Total coverage**: 31.6%

## Deployment

```bash
# Build for Linux
make build-linux

# Clean build artifacts
make clean
```

## Documentation

Documentation files in the project:

- [QUESTIONS_SETUP.md](QUESTIONS_SETUP.md) - Detailed guide for setting up custom questions
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Guide for migrating from old question format
- [SECURITY.md](SECURITY.md) - Security policy

## Dependencies

- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API

For current list of dependencies see `go.mod` file.

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `TELEGRAM_TOKEN` | - | Telegram bot token (required) |
| `QUESTIONS_FILE_PATH` | `configs/questions.json` | Path to questions file |
| `START_QUESTION_ID` | `start` | ID of the starting question |
| `DELAY_MS` | `700` | Default delay between messages (ms) |
