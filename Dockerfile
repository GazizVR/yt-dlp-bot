FROM golang:1.27.1
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.mod .
COPY go.sum .

COPY yt-dlp .

RUN apt update -y && apt install nodejs -y
RUN go build -o app cmd/main.go
CMD ["./app"]