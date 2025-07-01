FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -x -v -o penguin-server cmd/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/penguin-server .

EXPOSE 7777

COPY ./assets ./assets
CMD ["./penguin-server"]
