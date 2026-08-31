FROM golang:1.26.5
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.mod .
COPY go.sum .

COPY yt-dlp .
COPY cookies.txt .

RUN apt update -y && apt install nodejs -y
RUN go build -o app cmd/main.go
CMD ["./app"]