# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/build/sqm ./cmd/sqm

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy the binary from builder
COPY --from=builder /app/build/sqm /usr/local/bin/sqm

# Create non-root user
RUN adduser -D -g '' squaremind
USER squaremind

# Set environment variables
ENV SQUAREMIND_ENV=production

# Expose port for potential API
EXPOSE 8080

# Default command
ENTRYPOINT ["sqm"]
CMD ["--help"]
