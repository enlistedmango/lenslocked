#!/bin/bash
set -e

if [ -z "$HEROKU_APP_NAME" ]; then
    echo "❌ HEROKU_APP_NAME environment variable is required"
    exit 1
fi

echo "🔍 Testing Production Environment on $HEROKU_APP_NAME"

# Check application status
echo "📡 Checking application status..."
if ! heroku ps --app $HEROKU_APP_NAME | grep "web.*up" > /dev/null; then
    echo "❌ Web dyno not running"
    exit 1
fi

# Check database
echo "🔌 Testing database connection..."
if ! heroku pg:info --app $HEROKU_APP_NAME | grep "Status.*available" > /dev/null; then
    echo "❌ Database not available"
    exit 1
fi

# Test web endpoint
echo "🌐 Testing web endpoint..."
if ! curl -s https://$HEROKU_APP_NAME.herokuapp.com > /dev/null; then
    echo "❌ Web endpoint not responding"
    exit 1
fi

# Check for recent errors
echo "📋 Checking recent logs for errors..."
if heroku logs --app $HEROKU_APP_NAME --num 50 | grep -i "error\|fatal\|panic" > /dev/null; then
    echo "⚠️ Found errors in recent logs"
    heroku logs --app $HEROKU_APP_NAME --num 50 | grep -i "error\|fatal\|panic"
fi

echo "✅ Production environment tests passed!" 
