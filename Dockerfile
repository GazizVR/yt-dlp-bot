FROM golang:1.27.1
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.mod .
COPY go.sum .

RUN wget -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/download/2026.10.21/yt-dlp_linux
RUN chmod +x /usr/local/bin/yt-dlp

RUN apt update -y && apt install ffmpeg nodejs -y
RUN go build -o app cmd/main.go
CMD ["./app"]