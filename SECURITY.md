# Secure Project Setup

## Environment Variables

The application supports configuration via environment variables for secure storage of sensitive data.

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `TELEGRAM_TOKEN` | Telegram bot token | `1234567890:ABCDEFGHijklmnopQRSTUVWXYZ` |

### Optional Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `GOOGLE_CREDS` | Path to Google credentials file | `google-credentials.json` |
| `SHEET_ID` | Google Sheet ID | `YOUR_GOOGLE_SHEET_ID` |
| `DELAY_MS` | Delay between messages in ms | `700` |
| `START_QUESTION_ID` | First question ID | `start` |

## Methods of Setting Variables

### 1. Via .env file (recommended for development)

Create a `.env` file in the project root (automatically excluded from git):

```bash
TELEGRAM_TOKEN=your_bot_token
GOOGLE_CREDS=google-credentials.json
SHEET_ID=your_sheet_id
DELAY_MS=700
START_QUESTION_ID=start
```

### 2. Via command line

```bash
export TELEGRAM_TOKEN="your_bot_token"
export GOOGLE_CREDS="google-credentials.json"
export SHEET_ID="your_sheet_id"
export DELAY_MS=700
export START_QUESTION_ID="start"
```

### 3. When running the application

```bash
TELEGRAM_TOKEN="your_bot_token" ./tlgbot
```

## Getting Telegram Bot Token

1. Find @BotFather in Telegram
2. Send the `/newbot` command
3. Follow the instructions to create a bot
4. Copy the received token

## Fallback to Configuration File

If environment variables are not set, the application will attempt to load configuration from the `config.json` file. This is less secure and not recommended for production.

## Files Excluded from Git

- `.env*` - environment variable files
- `config.json.local` - local configurations
- `google-credentials.json` - Google credentials
- `config.json` - main configuration file (contains sensitive data)

## Production Security

For production, it's recommended to:

1. Use environment variables instead of files
2. Store secrets in specialized services (AWS Secrets Manager, Azure Key Vault, etc.)
3. Don't include sensitive data in code or configuration files
4. Regularly rotate tokens and keys
