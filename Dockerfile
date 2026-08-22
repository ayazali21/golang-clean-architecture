# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 → static binary, works in scratch/alpine without libc issues
RUN CGO_ENABLED=0 GOOS=linux go build -o /task-api ./cmd/api

FROM alpine:3.20

# certs needed for any outbound HTTPS calls; not strictly required yet but cheap to include
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /task-api .

EXPOSE 8080

CMD ["./task-api"]