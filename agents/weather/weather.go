package weatheragent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/openai/openai-go"
)

// WeatherAgent handles weather-related queries (basic version)
type WeatherAgent struct{}

// EnhancedWeatherAgent handles weather-related queries with real API data
type EnhancedWeatherAgent struct{}

// WeatherAPIResponse represents the structure of a weather API response
type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		Humidity int     `json:"humidity"`
		WindKph  float64 `json:"wind_kph"`
	} `json:"current"`
}

// NewBasic creates a new basic weather agent
func NewBasic() *WeatherAgent {
	return &WeatherAgent{}
}

// NewEnhanced creates a new enhanced weather agent
func NewEnhanced() *EnhancedWeatherAgent {
	return &EnhancedWeatherAgent{}
}

// Basic WeatherAgent methods
func (a *WeatherAgent) CanHandle(prompt string) bool {
	return canHandleWeatherPrompt(prompt)
}

func (a *WeatherAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	weatherContext := `You are a weather expert assistant. The user is asking about weather-related topics. 
	Provide helpful, accurate weather information. If they're asking about a specific location's weather, 
	let them know you don't have access to real-time weather data, but you can provide general weather 
	information, climate patterns, or suggest reliable weather sources like weather.com, weather.gov, or local meteorological services.
	
	If they're asking about weather concepts, phenomena, or general weather questions, provide detailed and educational responses.`

	ctx := context.Background()
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(weatherContext),
		openai.UserMessage(prompt),
	}

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get weather response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no weather response from OpenAI")
	}

	return "🌤️ [Weather Agent] " + completion.Choices[0].Message.Content, nil
}

func (a *WeatherAgent) GetName() string {
	return "Weather"
}

func (a *WeatherAgent) GetDescription() string {
	return "Specialized agent for weather-related queries and information"
}

// Enhanced WeatherAgent methods
func (a *EnhancedWeatherAgent) CanHandle(prompt string) bool {
	return canHandleWeatherPrompt(prompt)
}

func (a *EnhancedWeatherAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	location := a.extractLocation(prompt)

	if location != "" {
		// Check if API key is available, if not, automatically set it up
		apiKey := os.Getenv("WEATHER_API_KEY")
		if apiKey == "" {
			fmt.Println("\n🌤️ To provide real-time weather data, I need a Weather API key.")
			apiKey = a.promptForAPIKey()
			if apiKey != "" {
				os.Setenv("WEATHER_API_KEY", apiKey)
			}
		}

		weatherData, err := a.fetchWeatherData(location)
		if err == nil {
			weatherInfo := fmt.Sprintf(`🌤️ [Enhanced Weather Agent] Current weather for %s, %s:
• Temperature: %.1f°C (%.1f°F)
• Condition: %s
• Humidity: %d%%
• Wind: %.1f km/h

This is real-time weather data from WeatherAPI.com.`,
				weatherData.Location.Name,
				weatherData.Location.Country,
				weatherData.Current.TempC,
				weatherData.Current.TempF,
				weatherData.Current.Condition.Text,
				weatherData.Current.Humidity,
				weatherData.Current.WindKph)

			return weatherInfo, nil
		}
	}

	// Fallback to AI-powered weather assistant
	weatherContext := fmt.Sprintf(`You are a weather expert assistant. The user is asking about weather-related topics%s. 
	Provide helpful, accurate weather information. Since you don't have access to real-time weather data, 
	provide general weather information, climate patterns, or suggest reliable weather sources like weather.com, 
	weather.gov, AccuWeather, or local meteorological services.
	
	If they're asking about weather concepts, phenomena, or general weather questions, provide detailed and educational responses.`,
		func() string {
			if location != "" {
				return fmt.Sprintf(" for %s", location)
			}
			return ""
		}())

	ctx := context.Background()
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(weatherContext),
		openai.UserMessage(prompt),
	}

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get weather response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no weather response from OpenAI")
	}

	return "🌤️ [Enhanced Weather Agent] " + completion.Choices[0].Message.Content, nil
}

func (a *EnhancedWeatherAgent) GetName() string {
	return "Weather"
}

func (a *EnhancedWeatherAgent) GetDescription() string {
	return "Advanced weather agent with real-time data (requires WEATHER_API_KEY) and AI fallback"
}

func (a *EnhancedWeatherAgent) extractLocation(prompt string) string {
	locationPatterns := []string{
		`(?i)weather\s+in\s+([a-zA-Z\s,]+?)(?:\s|$|\?)`,
		`(?i)temperature\s+in\s+([a-zA-Z\s,]+?)(?:\s|$|\?)`,
		`(?i)forecast\s+for\s+([a-zA-Z\s,]+?)(?:\s|$|\?)`,
		`(?i)what.*weather.*in\s+([a-zA-Z\s,]+?)(?:\s|$|\?)`,
	}

	for _, pattern := range locationPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(prompt)
		if len(matches) > 1 {
			location := strings.TrimSpace(matches[1])
			location = regexp.MustCompile(`(?i)\s+(today|tomorrow|now|currently)$`).ReplaceAllString(location, "")
			return strings.TrimSpace(location)
		}
	}
	return ""
}

func (a *EnhancedWeatherAgent) fetchWeatherData(location string) (*WeatherAPIResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("no weather API key available")
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weatherData WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

	return &weatherData, nil
}

// promptForAPIKey prompts the user to enter their WeatherAPI key
func (a *EnhancedWeatherAgent) promptForAPIKey() string {
	fmt.Println("\n🌤️ Weather API Key Required")
	fmt.Println("─────────────────────────────")
	fmt.Println("To get real-time weather data, you need a free API key from WeatherAPI.com")
	fmt.Println("1. Visit: https://www.weatherapi.com/signup.aspx")
	fmt.Println("2. Sign up for a free account")
	fmt.Println("3. Get your API key from the dashboard")
	fmt.Println("4. Enter it below (or press Enter to skip)")
	fmt.Print("\n🔑 Enter your WeatherAPI key: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		apiKey := strings.TrimSpace(scanner.Text())
		if apiKey != "" {
			// Save to .env file for persistence
			if err := a.saveAPIKeyToEnv(apiKey); err != nil {
				fmt.Printf("⚠️  API key set for this session, but couldn't save to .env file: %v\n", err)
				fmt.Println("💡 To make this permanent, add this to your shell profile:")
				fmt.Printf("   export WEATHER_API_KEY='%s'\n", apiKey)
			} else {
				fmt.Println("✅ API key saved to .env file! Real-time weather data is now available.")
				fmt.Println("💡 The API key will persist across sessions.")
			}
			return apiKey
		}
	}

	fmt.Println("⚠️  Continuing without real-time weather data...")
	return ""
}

// saveAPIKeyToEnv saves the Weather API key to .env file
func (a *EnhancedWeatherAgent) saveAPIKeyToEnv(apiKey string) error {
	envFilePath := ".env"

	// Read existing .env file if it exists
	var existingLines []string
	var weatherKeyExists bool

	if file, err := os.Open(envFilePath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "WEATHER_API_KEY=") {
				// Replace existing weather API key line
				existingLines = append(existingLines, fmt.Sprintf("WEATHER_API_KEY=%s", apiKey))
				weatherKeyExists = true
			} else {
				existingLines = append(existingLines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading .env file: %w", err)
		}
	}

	// If WEATHER_API_KEY doesn't exist, add it
	if !weatherKeyExists {
		existingLines = append(existingLines, fmt.Sprintf("WEATHER_API_KEY=%s", apiKey))
	}

	// Write back to .env file
	file, err := os.Create(envFilePath)
	if err != nil {
		return fmt.Errorf("error creating .env file: %w", err)
	}
	defer file.Close()

	for _, line := range existingLines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return fmt.Errorf("error writing to .env file: %w", err)
		}
	}

	return nil
}

// Shared helper function
func canHandleWeatherPrompt(prompt string) bool {
	weatherKeywords := []string{
		"weather", "temperature", "rain", "snow", "sunny", "cloudy", "forecast",
		"hot", "cold", "humid", "wind", "storm", "hurricane", "tornado",
		"celsius", "fahrenheit", "degrees", "precipitation", "humidity",
	}

	promptLower := strings.ToLower(prompt)
	for _, keyword := range weatherKeywords {
		if strings.Contains(promptLower, keyword) {
			return true
		}
	}
	return false
}
