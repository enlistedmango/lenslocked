#!/bin/bash
set -e

echo "🔬 Running Integration Tests"

# Function to make HTTP requests and check responses
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local description=$4

    echo "Testing $description..."
    status=$(curl -s -o /dev/null -w "%{http_code}" -X $method http://localhost:3000$endpoint)
    
    if [ "$status" = "$expected_status" ]; then
        echo "✅ $description passed"
        return 0
    else
        echo "❌ $description failed (got $status, expected $expected_status)"
        return 1
    fi
}

# Test public endpoints
test_endpoint "GET" "/" 200 "Home page"
test_endpoint "GET" "/contact" 200 "Contact page"
test_endpoint "GET" "/faq" 200 "FAQ page"
test_endpoint "GET" "/signin" 200 "Sign in page"
test_endpoint "GET" "/signup" 200 "Sign up page"

# Test protected endpoints (should redirect to signin)
test_endpoint "GET" "/galleries" 302 "Galleries page (unauthorized)"

echo "✅ Integration tests completed!" 
