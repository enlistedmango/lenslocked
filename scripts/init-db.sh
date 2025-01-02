#!/bin/bash

# Wait for postgres to be ready
echo "Waiting for postgres..."
while ! pg_isready -h db -p 5432 -U postgres > /dev/null 2>&1; do
    sleep 1
done

echo "Running migrations..."
goose -dir migrations postgres "host=db port=5432 user=postgres password=postgres dbname=lenslocked sslmode=disable" up

echo "Database setup complete!" 
