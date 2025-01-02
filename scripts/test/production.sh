#!/bin/bash
set -e

if [ -z "$RAILWAY_STATIC_URL" ]; then
    echo "❌ RAILWAY_STATIC_URL environment variable is required"
    exit 1
fi

echo "🔍 Testing Production Environment on Railway"

# Check application status
echo "📡 Checking application status..."
if ! curl -s "$RAILWAY_STATIC_URL" > /dev/null; then
    echo "❌ Web service not responding"
    exit 1
fi

# Check for recent errors
echo "📋 Checking recent logs for errors..."
railway logs | grep -i "error\|fatal\|panic"

echo "✅ Production environment tests passed!" 
