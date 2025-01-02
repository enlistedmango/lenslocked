#!/bin/sh
set -e  # Exit on error
set -x  # Print commands as they're executed

echo "Starting database migration..."
echo "Migration directory: /app/migrations"
echo "Checking goose binary..."
ls -l /app/goose

echo "Checking migrations..."
ls -l /app/migrations

echo "Running migrations..."
/app/goose -dir /app/migrations postgres "$DATABASE_URL" up

echo "Checking migration status..."
/app/goose -dir /app/migrations postgres "$DATABASE_URL" status

echo "Migration complete!" 
