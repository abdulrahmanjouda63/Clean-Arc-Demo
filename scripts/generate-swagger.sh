#!/bin/bash

# Swagger generation script for Go project
# This script generates Swagger documentation from Go annotations

echo "Generating Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "swag command not found. Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate docs
echo "Running swag init..."
swag init -g main.go -o ./docs

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo "📁 Documentation files created in ./docs/"
    echo "🌐 Access Swagger UI at: http://localhost:8080/swagger/index.html"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi

echo "🎉 Swagger setup complete!"
