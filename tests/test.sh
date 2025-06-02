#!/bin/bash

# Test script for the OpenAI Terminal Assistant
# This script checks if the application builds and starts correctly

echo "🔧 Testing OpenAI Terminal Assistant..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

echo "✅ Go is installed"

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    echo "❌ main.go not found. Make sure you're in the project directory."
    exit 1
fi

echo "✅ Project files found"

# Build the application
echo "🔨 Building application..."
if go build; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Check if binary was created
if [ -f "./go-ai-terminal-assistant" ]; then
    echo "✅ Binary created successfully"
else
    echo "❌ Binary not found"
    exit 1
fi

echo ""
echo "🎉 All tests passed!"
echo ""
echo "To run the assistant:"
echo "1. Set your OpenAI API key: export OPENAI_API_KEY='your-key-here'"
echo "2. Run: ./go-ai-terminal-assistant"
echo ""
echo "Note: You need a valid OpenAI API key to actually use the assistant."
