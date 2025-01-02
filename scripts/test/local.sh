#!/bin/bash
set -e

echo "🔍 Testing Local Development Environment"

# Check Docker services
echo "📦 Checking Docker services..."
if ! docker-compose ps --services --filter "status=running" | grep -q "web"; then
    echo "❌ Web service not running"
    echo "Docker compose status:"
    docker-compose ps
    exit 1
fi

if ! docker-compose ps --services --filter "status=running" | grep -q "db"; then
    echo "❌ Database service not running"
    echo "Docker compose status:"
    docker-compose ps
    exit 1
fi

# Wait for web service to be actually ready
echo "🔄 Waiting for web service to be ready..."
max_attempts=30
attempt=1
while ! curl -s http://localhost:3000 > /dev/null && [ $attempt -le $max_attempts ]; do
    echo "Attempt $attempt of $max_attempts: Web service not ready yet..."
    sleep 2
    attempt=$((attempt + 1))
done

if [ $attempt -gt $max_attempts ]; then
    echo "❌ Web service failed to become ready within 60 seconds"
    echo "Docker logs:"
    docker-compose logs web
    exit 1
fi

# Test database connection
echo "🔌 Testing database connection..."
if ! docker-compose exec -T db psql -U postgres -d lenslocked -c "\dt" > /dev/null; then
    echo "❌ Database connection failed"
    echo "Docker logs:"
    docker-compose logs db
    exit 1
fi

# Test web service
echo "🌐 Testing web service..."
if ! curl -s http://localhost:3000 > /dev/null; then
    echo "❌ Web service not responding"
    echo "Docker logs:"
    docker-compose logs web
    exit 1
fi

# Test static file serving
echo "📁 Testing static file serving..."
if ! curl -s http://localhost:3000/static/css/output.css > /dev/null; then
    echo "❌ Static file serving failed"
    echo "Docker logs:"
    docker-compose logs web
    exit 1
fi

echo "✅ Local environment tests passed!" 
