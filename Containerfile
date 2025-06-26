# ---- Build stage ----
FROM golang:1.24.3 AS builder

# Enable Go modules
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy Go mod files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the codebase
COPY . .

# Build the binary
RUN go build -o /fz-tz ./cmd/main.go

# ---- Runtime stage ----
FROM debian:bullseye-slim

WORKDIR /app

# Install CA certificates for TLS verification
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /fz-tz .

EXPOSE 3333

ENTRYPOINT ["./fz-tz"]
