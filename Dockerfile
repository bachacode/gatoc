# syntax=docker/dockerfile:1

# Build stage - where we install tools and potentially build the final binary
FROM golang:1.24 as builder

WORKDIR /app

# Install 'air' in the builder stage's Go path
RUN go install github.com/air-verse/air@latest

# Copy Go modules and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your app (for building if you want a separate production binary, or for 'air' to see it)
COPY . ./

# Build the binary (this is for potential production use or initial compile, 'air' will re-build)
# Keep this as it's good practice for a multi-stage build, even if 'air' does its own thing in dev
RUN CGO_ENABLED=0 go build -o ./bot ./cmd/bot/main.go  # Adjust path to your main package


# Runtime stage - This needs to be a Go image for 'air' to function correctly
FROM golang:1.24

WORKDIR /app

# Copy the 'air' binary from the builder stage
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy your entire application source code into the runtime container
# This is essential for 'air' to watch and recompile your Go files
COPY . .

# Ensure Go modules are downloaded in this stage too, so 'air' can build successfully
RUN go mod download

# Copy .env file
COPY .env .

# Run 'air' as the command
CMD ["air", "-c", ".air.toml"]