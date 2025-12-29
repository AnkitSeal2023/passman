# Multi-stage build: build static Go binary, then run on Alpine

# Build stage
FROM golang:tip-trixie AS builder
WORKDIR /app

# Install build deps
RUN apt-get update && apt-get install -y git

# Cache go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

# Runtime stage
FROM alpine:3.20

# Create non-root user and add certs
RUN adduser -D -u 10001 appuser && apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary
COPY --from=builder /app/server /app/server

# Copy static assets used by the server
COPY views/static /app/views/static

# Configure
ENV PORT=5000
EXPOSE 5000

USER appuser
CMD ["/app/server"]
