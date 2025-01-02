#!/bin/sh
set -e  # Exit on error
set -x  # Print commands as they're executed

echo "Starting database migration..."
echo "Migration directory: /app/migrations"

echo "Running migrations..."
/app/goose -dir /app/migrations postgres "$DATABASE_URL" up
if [ $? -eq 0 ]; then
    echo "Migrations completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi

echo "Checking final migration status..."
/app/goose -dir /app/migrations postgres "$DATABASE_URL" status

echo "Migration process complete!" 
