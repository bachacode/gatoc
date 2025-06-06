# --- STAGE 1: Builder (Builds the Go binary) ---
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bot ./cmd/bot/main.go

# --- STAGE 2: Development Runtime (Includes 'air' for hot-reloading) ---
FROM builder as dev

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", ".air.toml"]

# --- STAGE 3: Production Runtime (Minimal, only the compiled binary) ---
FROM scratch as prod

WORKDIR /app

COPY --from=builder /app/bot .

CMD ["/app/bot"]
