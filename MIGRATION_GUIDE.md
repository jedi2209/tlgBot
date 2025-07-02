# Questions Migration Guide

## Quick migration to move your questions out of the repository

### What has changed?

✅ **Now you have:**

- Demo questions in the public repository (`configs/questions.json`)
- Ability to store your unique questions separately
- Configuration via `QUESTIONS_FILE_PATH` environment variable

### Migration steps

#### 1. Save your current questions

```bash
# Your original questions are already saved in:
ls configs/questions.example.json
```

#### 2. Create a file with your questions outside the repository

```bash
# Copy your questions to a safe place
cp configs/questions.example.json ~/my-bot-questions.json
```

#### 3. Set up the environment variable

```bash
# Add to your .bashrc, .zshrc or .env file:
export QUESTIONS_FILE_PATH="$HOME/my-bot-questions.json"
```

#### 4. Run the bot

```bash
# Now the bot will use your questions:
./telegram-bot
```

### Verification

Make sure everything works:

```bash
# Test with demo questions (default):
go build -o telegram-bot cmd/telegram-bot/main.go

# Test with your questions:
QUESTIONS_FILE_PATH="$HOME/my-bot-questions.json" ./telegram-bot
```

### What's happening in the repository?

- `configs/questions.json` - now contains demo questions
- `configs/questions.example.json` - contains your original questions as an example
- New files in `.gitignore` protect your private questions

### Docker

If using Docker:

```dockerfile
# In Dockerfile:
ENV QUESTIONS_FILE_PATH=/app/private/questions.json
COPY path/to/your/questions.json /app/private/questions.json
```

### Environment variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `QUESTIONS_FILE_PATH` | `configs/questions.json` | Path to questions file |
| `TELEGRAM_TOKEN` | - | Bot token (required) |

### Important! 🔒

- **Never** commit files with your unique questions to a public repository
- Use `.gitignore` to protect private files
- Keep backup copies of your questions

---

For detailed documentation, see [QUESTIONS_SETUP.md](QUESTIONS_SETUP.md)
