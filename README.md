# 🎬 yt-dlp-bot

A tiny Telegram bot that turns video links into files. Send a link, pick **Video** or **Audio**, and get it back in chat. ✨

Built with Go, [yt-dlp](https://github.com/yt-dlp/yt-dlp) and a sprinkle of Valkey.

## How it works

1. 🔗 Send the bot a link
2. 👇 Tap **🎬 Video** or **🎧 Audio**
3. 🎉 Receive your file

Videos come as `mp4`, audio as `m4a` when available.

## Quick start with Docker 🐳

Create a `.env` file:

```env
BOT_TOKEN=your-bot-token
TELEGRAM_API_ID=your-api-id
TELEGRAM_API_HASH=your-api-hash
TELEGRAM_HTTP_PORT=8081
```

Then run:

```sh
docker compose up -d --build
```

That's it! Compose starts the bot, a local Telegram Bot API server and Valkey for you.

> 💡 You can get `BOT_TOKEN` from [@BotFather](https://t.me/BotFather) and the API ID/hash from [my.telegram.org](https://my.telegram.org).

## Run locally 🛠️

You'll need Go, `yt-dlp`, a JS runtime (like `quickjs`) and a running Valkey/Redis.

```sh
./start.sh
```

It loads your `.env` and starts the bot.

## Configuration ⚙️

| Variable | Required | What it does |
| --- | :---: | --- |
| `BOT_TOKEN` | ✅ | Your Telegram bot token |
| `JS_RUNTIME` | ✅ | JS runtime for yt-dlp (e.g. `quickjs`) |
| `REDIS_ADDR` | ✅ | Valkey/Redis address (e.g. `localhost:6379`) |
| `REDIS_PASSWORD` | | Valkey/Redis password |
| `BOT_API_BASE_URL` | | Local Telegram Bot API server URL |
| `BOT_FILES_DIR` | | Where downloaded files are stored |
| `YT_DLP_BIN` | | Path to yt-dlp (default: `yt-dlp`) |
| `COOKIES` | | Path to a `cookies.txt` file |
| `BROWSER` | | Browser to grab cookies from |

## License 📄

MIT, have fun!