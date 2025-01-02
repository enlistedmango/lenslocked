#!/bin/bash

echo "🔍 Testing Local Development Environment"

# Check Docker services
echo "Checking Docker services..."
if ! docker-compose ps | grep "web.*running" > /dev/null; then
    echo "❌ Web service not running"
    exit 1
fi

if ! docker-compose ps | grep "db.*running" > /dev/null; then
    echo "❌ Database service not running"
    exit 1
fi

# Test database connection
echo "Testing database connection..."
if ! docker-compose exec -T db psql -U postgres -d lenslocked -c "\dt" > /dev/null; then
    echo "❌ Database connection failed"
    exit 1
fi

# Test web service
echo "Testing web service..."
if ! curl -s http://localhost:8080 > /dev/null; then
    echo "❌ Web service not responding"
    exit 1
fi

echo "✅ Local environment tests passed"

# If HEROKU_APP_NAME is set, test production
if [ ! -z "$HEROKU_APP_NAME" ]; then
    echo "🔍 Testing Production Environment"
    
    # Check Heroku app status
    if ! heroku ps --app $HEROKU_APP_NAME | grep "web.*up" > /dev/null; then
        echo "❌ Heroku web dyno not running"
        exit 1
    fi
    
    # Check Heroku database
    if ! heroku pg:info --app $HEROKU_APP_NAME | grep "Status.*available" > /dev/null; then
        echo "❌ Heroku database not available"
        exit 1
    fi
    
    echo "✅ Production environment tests passed"
fi

echo "🎉 All tests passed!" 
