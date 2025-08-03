#!/bin/bash

echo "🧪 Testing ScholaTheologiae API with DDD Architecture"
echo "=================================================="

# Start the server in background
echo "🚀 Starting server..."
./bin/schola-theologiae-api &
SERVER_PID=$!

# Wait for server to start
sleep 2

echo "🏥 Testing health endpoint..."
curl -s http://localhost:8080/v1/health | jq '.' 2>/dev/null || echo "Health check failed"

echo ""
echo "📚 Testing Summa Theologiae parts..."
curl -s http://localhost:8080/v1/summa-theologiae | jq '.' 2>/dev/null || echo "Summa parts failed"

echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null

echo "✅ DDD Architecture test completed!"
