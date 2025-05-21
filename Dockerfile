# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24 as builder

WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your app
COPY . ./

# Build the binary
RUN CGO_ENABLED=0 go build -o ./bot ./cmd/bot/main.go  # Adjust path to your main package

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/bot /app/bot

# Copy .env file if you need it
COPY .env .

# Run it
CMD ["./bot"]
