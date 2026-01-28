FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o engine-cli ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/engine-cli .

ENTRYPOINT ["./engine-cli"]