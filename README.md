# 🤖 yt-dlp Telegram Bot

> A lightweight Telegram bot written in **Go** that downloads videos and audio using **yt-dlp** and sends them directly back to users.

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?logo=docker)
![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?logo=telegram)
![License](https://img.shields.io/badge/License-MIT-green)

---

## ✨ Features

- 🎥 Download videos from hundreds of supported websites via **yt-dlp**
- 🎵 Download audio or video content
- 🤖 Telegram Bot API integration
- 🍪 Supports authenticated downloads using cookies
- ⚡ Fast and lightweight Go implementation
- 🐳 Docker & Docker Compose support
- 🧩 Configurable JavaScript runtime for yt-dlp

---

## 📁 Project Structure

```
.
├── cmd/                 # Application entrypoint
├── internal/
│   ├── config/          # Configuration loader
│   └── service/         # Bot logic & handlers
├── pkg/
│   ├── telegram/        # Telegram API client
│   └── ytdlp/           # yt-dlp wrapper
├── Dockerfile
├── compose.yaml
├── go.mod
├── yt-dlp
└── cookies.txt
```

---

## ⚙️ Configuration

The application uses environment variables.

| Variable | Description |
|----------|-------------|
| `TOKEN` | Telegram Bot Token |
| `BASE_URL` | Telegram Bot API URL |
| `YT_DLP_PATH` | Path to the `yt-dlp` executable |
| `COOKIES_PATH` | Path to cookies file |
| `JS_RUNTIME` | JavaScript runtime used by yt-dlp |
| `BROWSER` | Browser profile (optional, for cookie extraction) |

---

## 🚀 Running with Docker

```bash
docker compose up --build
```

---

## 🛠️ Running Locally

### 1. Clone the repository

```bash
git clone <repository-url>
cd yt-dlp-bot
```

### 2. Install dependencies

- Go 1.26+
- Node.js
- yt-dlp

### 3. Configure environment

```bash
export TOKEN=your_bot_token
export BASE_URL=http://localhost:8081
export YT_DLP_PATH=./yt-dlp
export COOKIES_PATH=./cookies.txt
export JS_RUNTIME=node
```

### 4. Run

```bash
go run ./cmd
```

---

## 🐳 Docker Compose

The provided `compose.yaml` starts:

- 📡 Telegram Bot API server
- 🤖 Bot service

Simply execute:

```bash
docker compose up -d
```

---

## 🔧 Technologies

- Go
- Telegram Bot API
- yt-dlp
- Docker
- Docker Compose
- Node.js

---

## 📌 Notes

- Some websites require authentication. Use a valid `cookies.txt` file.
- Make sure the `yt-dlp` binary is executable.
- A JavaScript runtime (such as Node.js) may be required for certain extractors.

---

## 📄 License

This project is distributed under the **MIT License**.

Feel free to fork, modify and contribute! 🚀

---

## ⭐ Contributing

Contributions, issues and feature requests are always welcome.

If you find this project useful, consider giving it a ⭐ on GitHub.
