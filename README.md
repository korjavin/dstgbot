# DeepSeek Telegram Bot

A Telegram bot that uses DeepSeek API to generate responses in group chats.

## Features

- Responds to mentions and replies
- Maintains conversation context (up to 10 messages)
- Handles text messages only (images/audio not supported)
- Retries failed API calls
- In-memory message caching for context (including bot's own messages)
- Detailed logging for debugging

## Prerequisites

- Go 1.23+
- Telegram bot token
- DeepSeek API key
- Telegram group ID

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/korjavin/dstgbot.git
   cd dstgbot
   ```

2. Set up environment variables:
   ```bash
   export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
   export DEEPSEEK_APIKEY="your-deepseek-api-key"
   export TG_GROUP_ID="your-group-id"
   export SYSTEM_MSG="Your system message for DeepSeek"
   ```

3. Build and run:
   ```bash
   go build -o dstgbot .
   ./dstgbot
   ```

## Podman Instructions

1. Build the container:
   ```bash
   podman build -t dstgbot .
   ```

2. Run the container:
   ```bash
   podman run -d \
     -e TELEGRAM_BOT_TOKEN="your-token" \
     -e DEEPSEEK_APIKEY="your-api-key" \
     -e TG_GROUP_ID="your-group-id" \
     -e SYSTEM_MSG="Your system message" \
     --name dstgbot \
     dstgbot
   ```

## GitHub Actions

The repository includes a GitHub Actions workflow that automatically builds and pushes a Docker image to Docker Hub when changes are pushed to the master branch.

To use this:

1. Set up Docker Hub secrets in your GitHub repository:
   - `DOCKER_HUB_USERNAME`
   - `DOCKER_HUB_TOKEN`

2. Push changes to the master branch to trigger the build.

## Configuration

Environment Variables:
- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token
- `DEEPSEEK_APIKEY`: Your DeepSeek API key
- `TG_GROUP_ID`: ID of the Telegram group to monitor
- `SYSTEM_MSG`: System message for DeepSeek (optional)

## License

MIT