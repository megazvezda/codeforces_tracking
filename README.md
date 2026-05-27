# Codeforces Telegram Bot

A lightweight Telegram bot that monitors Codeforces submissions and notifies you when tracked users solve problems.
Monitors multiple Codeforces handles, instant notifications for accepted submissions. Shows wrong submission count before success, displays time, memory, and compiler info.

## Setup

### 1. Create a Telegram Bot

1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Send `/newbot` and follow instructions
3. Copy the bot token

### 2. Get Your Chat ID

#### For Private Chat:
1. Start a chat with your new bot
2. Send any message to it
3. Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Find your `chat_id` in the JSON response (will be a positive number)

#### For Public Channel:
1. Create a public channel or use an existing one
2. Add your bot as an administrator with "Post Messages" permission
3. Send a message to the channel (or forward a message there)
4. Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
5. Find your channel's `chat_id` in the JSON response (starts with `-100`, e.g., `-1001234567890`)
6. Use this channel ID as `telegram_chat_id` in config

### 3. Configure the Bot

Edit `config.json`:

```json
{
  "telegram_bot_token": "YOUR_BOT_TOKEN_HERE",
  "telegram_chat_id": YOUR_CHAT_ID_HERE,
  "codeforces_handles": [
    {
      "handle": "tourist",
      "real_name": "Gennady Korotkevich"
    }
  ],
  "poll_interval_seconds": 120
}
```

- `telegram_bot_token`: Your bot token from BotFather
- `telegram_chat_id`: Your chat ID or channel ID (numeric, no quotes)
  - Private chat: positive number (e.g., `123456789`)
  - Public channel: negative number starting with -100 (e.g., `-1001234567890`)
- `codeforces_handles`: Array of handles to monitor
- `poll_interval_seconds`: How often to check (recommended: 120-300 seconds)

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Bot

```bash
go run main.go
```

Or build and run:

```bash
go build -o cfbot
./cfbot
```

## Notification Format

```
✅ tourist (Gennady Korotkevich) has solved task B) Binary Search after 2 wrong submissions🎉

🧪 Tests passed: 47
⏱ Time: 156 ms
💾 Memory: 1.25 MB
💻 Language: GNU C++20 (64)
🕐 Solved at: 2026-02-10 14:35:22 (UTC+2)
🔗 Submission #123456789
```

The bot automatically detects whether a problem is from a regular contest or gym and generates the correct URL format.

## Sources

The official Codeforces API:
- [Codeforces API Documentation](https://codeforces.com/apiHelp)
- [PublicAPI.dev - Codeforces API](https://publicapi.dev/codeforces-api)
