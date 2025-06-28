FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./backend/
COPY backend/main.go ./backend/

WORKDIR /app/backend

RUN go mod download

RUN go build -o /app/main .

EXPOSE 8080

CMD ["/app/main"]
