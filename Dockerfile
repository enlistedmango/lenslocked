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

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
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

# Default port (will be overridden by Heroku)
ENV PORT=3000

# Expose port
EXPOSE ${PORT}

# Command to run
CMD ["./main"]
