#!/bin/bash

echo "🚀 Starting Production Deployment Fix..."

# Step 1: Stop everything
echo "📦 Stopping containers..."
docker-compose down

# Step 2: Clean up old data
echo "🗑️  Removing old MySQL volume..."
docker volume rm speed-test-server_mysql_data || true

# Step 3: Verify .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your production settings!"
    exit 1
fi

# Step 4: Build fresh images
echo "🔨 Building fresh Docker images..."
docker-compose build --no-cache

# Step 5: Start services
echo "🚀 Starting services..."
docker-compose up -d

# Step 6: Wait for MySQL to be ready
echo "⏳ Waiting for MySQL to initialize..."
sleep 30

# Step 7: Check database connection
echo "🔍 Checking database connection..."
docker-compose logs app | grep -i "database" | tail -5

# Step 8: Verify web files
echo "🔍 Verifying web files in container..."
docker-compose exec app ls -la /app/web/

# Step 9: Test endpoints
echo "🧪 Testing endpoints..."
sleep 5
echo "Testing healthz endpoint:"
curl -I http://localhost:8080/healthz
echo ""
echo "Testing speedtest.html:"
curl -I http://localhost:8080/speedtest.html
echo ""
echo "Testing /new endpoint:"
curl -I http://localhost:8080/new

echo ""
echo "✅ Deployment complete!"
echo "📊 View logs: docker-compose logs -f"
echo "🌐 Access app: http://localhost:8080"
