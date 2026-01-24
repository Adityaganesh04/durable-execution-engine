# Stage 1: Build the application
# We use the 1.25 image to match your local development environment
FROM golang:1.25-alpine AS builder

# Install build tools (GCC/Musl) required for SQLite CGO support
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# CGO_ENABLED=1 is mandatory for go-sqlite3
# -ldflags="-w -s" strips debug info to make the binary smaller
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o engine-cli ./cmd/main.go

# Stage 2: Create the minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/engine-cli .

# Command to run the application
ENTRYPOINT ["./engine-cli"]