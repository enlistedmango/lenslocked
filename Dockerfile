# Build stage
FROM golang:1.23-rc AS builder

# Set working directory
WORKDIR /app

# Install PostgreSQL client and goose
RUN apt-get update && apt-get install -y postgresql-client && \
    CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and scripts
COPY . .
COPY scripts/init-db.sh /init-db.sh
COPY scripts/migrate.sh /app/scripts/migrate.sh
RUN chmod +x /init-db.sh /app/scripts/migrate.sh

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
FROM alpine:latest AS production

# Install ca-certificates, PostgreSQL client, bash and other dependencies
RUN apk --no-cache add \
    ca-certificates \
    postgresql-client \
    bash \
    libc6-compat \
    libgcc \
    libstdc++

WORKDIR /app

# Copy goose from builder and make it executable
COPY --from=builder /go/bin/goose /app/goose
RUN chmod +x /app/goose

# Copy the binary and assets from builder
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts/migrate.sh /app/scripts/migrate.sh
RUN chmod +x /app/scripts/migrate.sh

# Default port (will be overridden by Heroku)
ENV PORT=3000

# Expose port
EXPOSE ${PORT}

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT}/health || exit 1

# Command to run
CMD ["./main"]
