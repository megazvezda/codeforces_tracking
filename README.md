# Codeforces Tracking Bot

This is a Telegram bot that monitors Codeforces user submissions and sends notifications when a tracked user solves a problem.

## Features

-   Monitors multiple Codeforces handles.
-   Sends notifications to a Telegram chat for accepted submissions.
-   Tracks wrong attempts before a problem is solved.
-   Includes problem details, submission stats, and links in notifications.

## Getting Started

### Prerequisites

-   Go 1.16 or higher.
-   A Telegram bot token.
-   Your Telegram chat ID.

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/your-username/cfvgtuai.git
    cd cfvgtuai
    ```

2.  Create a `config.json` file in the `configs` directory by copying the example:
    ```sh
    cp configs/config.json.example configs/config.json
    ```

3.  Edit `configs/config.json` and add your Telegram bot token, chat ID, and the Codeforces handles you want to track.

    ```json
    {
      "telegram_token": "YOUR_TELEGRAM_BOT_TOKEN",
      "telegram_chat_id": 123456789,
      "codeforces_handles": [
        "handle1",
        "handle2"
      ],
      "polling_interval_seconds": 60
    }
    ```

### Running the Bot

```sh
go run cmd/codeforces-tracking/main.go
```
