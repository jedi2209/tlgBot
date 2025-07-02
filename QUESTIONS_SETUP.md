# User Questions Setup

This document describes how to set up your unique questions separately from the public code.

## Overview

The system supports two types of question configuration:

- **Demo questions** - stored in `configs/questions.json` in the public repository
- **User questions** - stored in a separate file outside the repository

## Quick Setup

### 1. Create a file with your questions

Create a file with your unique questions, for example `my-questions.json`:

```bash
cp configs/questions.example.json my-questions.json
```

### 2. Configure the path to the file

There are two ways to specify the path to your questions file:

#### Method 1: Environment variable (recommended)

```bash
export QUESTIONS_FILE_PATH="/path/to/your/my-questions.json"
```

#### Method 2: Configuration file

Create `config.json` based on `configs/config.example.json` and specify:

```json
{
  "questions_file_path": "/path/to/your/my-questions.json"
}
```

### 3. Run the bot

```bash
# With environment variable
QUESTIONS_FILE_PATH="my-questions.json" ./telegram-bot

# Or with configuration file
./telegram-bot config.json
```

## Questions file structure

The questions file should contain a JSON array with question objects:

```json
[
  {
    "id": "start",
    "messages": [
      "Hello, {name}!",
      "Now I will ask you a few questions."
    ],
    "auto_advance": true,
    "auto_advance_delay_ms": 1500,
    "options": [
      {"text": "Continue", "next_id": "question_1"}
    ]
  },
  {
    "id": "question_1",
    "text": "What type of air conditioner do you have?",
    "options": [
      {"text": "Ducted", "next_id": "question_2"},
      {"text": "Wall-mounted", "next_id": "question_2"}
    ]
  }
]
```

## Question object fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique question identifier |
| `text` | string | Main question text |
| `messages` | array | Array of messages for step-by-step display |
| `images` | array | Paths to images |
| `options` | array | Answer options with transitions |
| `auto_advance` | boolean | Automatic transition to next question |
| `auto_advance_delay_ms` | number | Delay for auto-advance |
| `delay_ms` | number | Delay before showing question |
| `input_type` | string | Input type ("text") |
| `input_placeholder` | string | Placeholder for input field |
| `external_link` | string | External link |
| `external_text` | string | Text for external link |

## Security

- Never commit files with your unique questions to a public repository
- Store confidential questions in private files
- Use environment variables for paths to private files

## Configuration examples

### For development

```bash
export QUESTIONS_FILE_PATH="my-dev-questions.json"
```

### For production

```bash
export QUESTIONS_FILE_PATH="/opt/bot/private/questions.json"
```

### Docker

```dockerfile
ENV QUESTIONS_FILE_PATH=/app/private/questions.json
COPY private/questions.json /app/private/questions.json
```

## Migrating existing questions

If you already have a `configs/questions.json` file with your questions:

1. Copy it to a safe location:

   ```bash
   cp configs/questions.json ~/my-bot-questions.json
   ```

2. Set up the environment variable:

   ```bash
   export QUESTIONS_FILE_PATH="$HOME/my-bot-questions.json"
   ```

3. Now the repository will only contain demo questions

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `QUESTIONS_FILE_PATH` | `configs/questions.json` | Path to questions file |
| `TELEGRAM_TOKEN` | - | Telegram bot token (required) |
| `START_QUESTION_ID` | `start` | Start question ID |
| `DELAY_MS` | `700` | Default delay between messages |

## Troubleshooting

### Error "Failed to load questions"

- Check that the file exists at the specified path
- Make sure the JSON is valid
- Check file access permissions

### Error "question not found"

- Make sure the start question exists
- Check that all `next_id` references point to existing questions

### Error "duplicate question ID"

- Make sure all question IDs are unique
