FROM golang:1.26.3
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.mod .
COPY go.sum .

COPY yt-dlp .

RUN go build -o app cmd/main.go
CMD ["./app"]