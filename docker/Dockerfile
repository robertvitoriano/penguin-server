FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o penguin-server cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/penguin-server .

EXPOSE 7777

CMD ["./penguin-server"]
