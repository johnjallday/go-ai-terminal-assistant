#!/bin/bash

# Test script for the agentic routing system
# This script validates that agents are properly routing based on different prompt types

echo "🧪 Testing Agentic Routing System"
echo "=================================="

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ OPENAI_API_KEY environment variable is not set"
    echo "💡 Please set your OpenAI API key: export OPENAI_API_KEY='your-key-here'"
    exit 1
fi

# Build the application
echo "🔨 Building application..."
if go build .; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "🎯 Agent Routing Test Cases:"
echo "----------------------------"

echo "1. Math Agent Test Cases:"
echo "   - 'calculate 25 * 17'"
echo "   - 'what is 15% of 200?'"
echo "   - 'solve for x: 2x + 5 = 15'"
echo "   - 'what's the square root of 144?'"

echo ""
echo "2. Weather Agent Test Cases:"
echo "   - 'what's the weather like today?'"
echo "   - 'is it going to rain tomorrow?'"
echo "   - 'temperature in New York'"
echo "   - 'weather forecast for London'"

echo ""
echo "3. Default Agent Test Cases:"
echo "   - 'hello, how are you?'"
echo "   - 'tell me a joke'"
echo "   - 'what's the capital of France?'"
echo "   - 'write a short poem'"

echo ""
echo "🚀 To test the application:"
echo "   1. Run: ./go-ai-terminal-assistant"
echo "   2. Try the test prompts above"
echo "   3. Use '/agents' to see all available agents"
echo "   4. Use '/model' to change AI models"

echo ""
echo "🔍 Expected Routing Behavior:"
echo "   - Math prompts → 🧮 [Math Agent]"
echo "   - Weather prompts → 🌤️ [Enhanced Weather Agent]"
echo "   - Other prompts → Default Agent (no special prefix)"

echo ""
echo "⚙️ Optional Weather API Setup:"
echo "   1. Get free API key from https://www.weatherapi.com/"
echo "   2. Set: export WEATHER_API_KEY='your-weather-key'"
echo "   3. Test with location-specific queries like 'weather in Tokyo'"
