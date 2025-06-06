# --- STAGE 1: Builder (Builds the Go binary) ---
FROM golang:1.24-alpine as builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build -o ./bot ./cmd/bot/main.go

# --- STAGE 2: Development Runtime (Includes 'air' for hot-reloading) ---
FROM builder as dev
RUN apk add --no-cache ffmpeg ca-certificates
ENV SSL_CERT_DIR=/etc/ssl/certs
RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

# --- STAGE 3: Production Runtime (Minimal, only the compiled binary) ---
FROM scratch as prod

WORKDIR /app

COPY --from=builder /app/bot .

CMD ["/app/bot"]
