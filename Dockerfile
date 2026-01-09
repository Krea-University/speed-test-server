# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application from the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o speed-test-server ./cmd/speed-test-server

# Runtime stage
FROM alpine:3.20

# Install runtime dependencies and security updates
RUN apk --no-cache add ca-certificates tzdata && \
    apk --no-cache upgrade && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/speed-test-server .

# Copy web folder for HTML pages
COPY --from=builder /app/web ./web

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Install wget for healthcheck
USER root
RUN apk --no-cache add wget
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run the binary
CMD ["./speed-test-server"]
