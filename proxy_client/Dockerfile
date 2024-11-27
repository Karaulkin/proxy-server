FROM golang:1.23.2 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o client ./cmd/client/main.go

FROM ubuntu:22.04

WORKDIR /app
COPY --from=builder /app/client ./
CMD ./client -urls="https://www.youtube.com/watch?v=JqfbfraMjW8,https://www.youtube.com/watch?v=ziotOZhwz9Y"
