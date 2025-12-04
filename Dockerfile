# ---- Build Stage ----
FROM golang:1.25.1 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code
COPY . .

# Build the Go binary
RUN go build -ldflags="-s -w" -o app ./cmd/api

# ---- Run Stage ----
FROM alpine:latest

# Install certs for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the built binary from builder stage
COPY --from=builder /app/app .

RUN chmod +x /app/app

EXPOSE 8080

CMD ["./app"]
