FROM golang:1.23.2 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server ./cmd/server/main.go

FROM ubuntu:22.04
WORKDIR /app
COPY ./.env .env
COPY --from=builder /app/server .
COPY ./config ./config
COPY ./storage ./storage
ENV CONFIG_PATH=/app/config/local.yaml
CMD ./server
