# Telegram Bot with DeepSeek API - Plan

## Project Overview

This project aims to create a Telegram bot in Go that utilizes the DeepSeek API to generate responses. The bot will:

*   Read the Telegram bot token and DeepSeek API key from environment variables.
*   Connect to the Telegram API and listen for messages in a specified group.
*   Detect mentions of the bot or replies to the bot's messages.
*   Send the message text to the DeepSeek API for a response.
*   Send the DeepSeek API response back to the Telegram group.
*   Handle unsupported message types (images, voice, audio) gracefully.
*   Implement error handling and retry logic for API calls.
*   Maintain an in-memory cache of recent messages to provide context for conversations.
*   Include a comprehensive README file with instructions for setup and usage.
*   Provide instructions for building and running the bot using Podman.
*   Automate Docker image builds using GitHub Actions.

## Architecture Diagram

```mermaid
graph LR
    A[Start] --> B{Project Setup};
    B --> C{Telegram Bot Logic};
    C --> D{DeepSeek API Integration};
    D --> E{Message Handling};
    E --> F{Error Handling & Retries};
    F --> G{Caching};
    G --> H{README};
    H --> I{Podman Instructions};
    I --> J{GitHub Actions};
    J --> K[End];

    subgraph Project Setup
        B1[Initialize Go module]
        B2[Create project structure (api, telegram, cache)]
        B3[Define dependencies]
    end

    subgraph Telegram Bot Logic
        C1[Read Telegram token from env]
        C2[Connect to Telegram API]
        C3[Get bot name and login from API]
        C4[Listen for messages in group]
    end

    subgraph DeepSeek API Integration
        D1[Read DeepSeek API key from env]
        D2[Implement API request function (chat completion endpoint)]
        D3[Handle API responses]
    end

    subgraph Message Handling
        E1[Detect bot mentions/replies]
        E2[Extract message text]
        E3[Send message to DeepSeek API]
        E4[Send DeepSeek response to Telegram]
        E5[Handle unsupported message types]
    end

    subgraph Error Handling & Retries
        F1[Implement retry logic for API calls]
        F2[Log errors]
        F3[Handle Telegram API errors]
    end

    subgraph Caching
        G1[Implement in-memory cache (at least 300 messages)]
        G2[Store messages with reply_to ID]
        G3[Recursively retrieve context for replies using reply_to]
    end

    subgraph README
        H1[Describe project architecture]
        H2[Explain how to get API keys]
        H3[Provide usage instructions]
    end

   subgraph Podman Instructions
        I1[Write instructions for building and running with Podman]
    end

    subgraph GitHub Actions
        J1[Create Dockerfile]
        J2[Set up GitHub Actions workflow]
        J3[Build and push Docker image]
    end
```

## Detailed Breakdown

### 1. Project Setup

*   Initialize a Go module: `go mod init github.com/korjavin/dstgbot`
*   Create a project structure with a `main.go` file and packages: `api`, `telegram`, `cache`.
*   Define the necessary dependencies in `go.mod`:
    *   `github.com/go-telegram-bot-api/telegram-bot-api`
    *   `github.com/deepseek-ai/deepseek-go`

### 2. Telegram Bot Logic

*   Read the Telegram bot token from the environment using `os.Getenv("TELEGRAM_BOT_TOKEN")`.
*   Connect to the Telegram API using the `telegram-bot-api` library.
*   Get the bot's name and login from the Telegram API.
*   Listen for messages in the specified Telegram group using the group ID from the environment (`os.Getenv("TG_GROUP_ID")`).

### 3. DeepSeek API Integration

*   Read the DeepSeek API key from the environment using `os.Getenv("DEEPSEEK_APIKEY")`.
*   Implement a function to make requests to the DeepSeek API's chat completion endpoint, including setting the necessary headers and handling authentication.
*   Handle the API responses, including parsing the JSON and extracting the relevant information.

### 4. Message Handling

*   Detect if the bot is mentioned in a message or if someone is replying to the bot.
*   Extract the text from the message.
*   Send the message to DeepSeek API to generate a response.
*   Send the DeepSeek API response back to the Telegram group.
*   Handle unsupported message types (images, voice, audio) by sending a message indicating that they are not supported.

### 5. Error Handling & Retries

*   Implement retry logic for API calls to handle potential network issues or API downtime.
*   Log any errors that occur during the process.
*   Handle Telegram API errors, such as invalid tokens or connection issues.

### 6. Caching

*   Implement an in-memory cache to store at least 300 messages in the conversation.
*   Store messages with their `reply_to` ID.
*   When someone replies to the bot, recursively retrieve context for replies using the `reply_to` ID to build the conversation history.

### 7. README

*   Write a comprehensive README file that describes the project architecture, how to get the necessary API keys, and how to use the bot.

### 8. Podman Instructions

*   Write instructions for building and running the bot using Podman, including how to set the environment variables.

### 9. GitHub Actions

*   Create a Dockerfile to containerize the bot.
*   Set up a GitHub Actions workflow to automatically build and push the Docker image to a container registry whenever there are changes to the `master` branch.