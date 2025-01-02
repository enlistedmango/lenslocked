# Build stage
FROM golang:1.23-rc AS builder

# Set working directory
WORKDIR /app

# Install PostgreSQL client and goose
RUN apt-get update && apt-get install -y postgresql-client && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and scripts
COPY . .
COPY scripts/init-db.sh /init-db.sh
RUN chmod +x /init-db.sh

# For development, we'll use the builder stage
# For production, we'll use the following:

# Build the application
FROM alpine:latest AS production

# Install ca-certificates and PostgreSQL client
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /app

# Copy the binary and assets from builder
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 3000

# Command to run
CMD ["./main"]
