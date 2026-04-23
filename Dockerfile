# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/api .
COPY config.yaml .

EXPOSE 4000

CMD ["./api"]