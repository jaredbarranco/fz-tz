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

# Copy binary from builder
COPY --from=builder /fz-tz .

# Port your app listens on (adjust if needed)
EXPOSE 8080

ENTRYPOINT ["./fz-tz"]
