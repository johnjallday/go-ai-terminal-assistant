package mathagent

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/openai/openai-go"
)

// MathAgent handles mathematical calculations and problems
type MathAgent struct{}

func New() *MathAgent {
	return &MathAgent{}
}

func (a *MathAgent) CanHandle(prompt string) bool {
	mathKeywords := []string{
		"calculate", "compute", "solve", "math", "equation", "formula",
		"add", "subtract", "multiply", "divide", "plus", "minus", "times",
		"square", "root", "power", "logarithm", "sin", "cos", "tan",
		"derivative", "integral", "geometry", "algebra", "trigonometry",
		"percentage", "percent", "fraction", "decimal", "probability",
	}

	mathSymbols := []string{
		"+", "-", "*", "/", "=", "^", "√", "∫", "∂", "%",
		"π", "∞", "≈", "≠", "≤", "≥", "<", ">",
	}

	promptLower := strings.ToLower(prompt)

	// Check for keywords
	for _, keyword := range mathKeywords {
		if strings.Contains(promptLower, keyword) {
			return true
		}
	}

	// Check for mathematical symbols (but be careful with single letters)
	for _, symbol := range mathSymbols {
		if strings.Contains(prompt, symbol) {
			return true
		}
	}

	// Check for mathematical constants like 'e' and 'pi' more carefully
	// Only match if they appear in mathematical context
	mathConstPattern := regexp.MustCompile(`(?i)\b(e|pi)\b\s*[=\^\*\+\-\/]|\b(log|ln|exp)\s*\(|[=\^\*\+\-\/]\s*\b(e|pi)\b`)
	if mathConstPattern.MatchString(prompt) {
		return true
	}

	// Check for number patterns that might indicate calculations
	numberPattern := regexp.MustCompile(`\d+\s*[+\-*/^]\s*\d+`)
	if numberPattern.MatchString(prompt) {
		return true
	}

	return false
}

func (a *MathAgent) Handle(prompt string, client *openai.Client, model string) (string, error) {
	mathContext := `You are a mathematical expert and calculator. The user is asking about mathematical problems, calculations, or concepts.

Guidelines:
1. For numerical calculations, provide step-by-step solutions
2. Show your work clearly with intermediate steps
3. For complex math concepts, provide clear explanations with examples
4. Use proper mathematical notation when helpful
5. If the problem involves specific formulas, mention them
6. For word problems, break down the problem into mathematical components
7. Always double-check your calculations
8. If approximate, clearly state it's an approximation

Provide accurate, educational, and well-structured mathematical responses.`

	ctx := context.Background()
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(mathContext),
		openai.UserMessage(prompt),
	}

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    model,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get math response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no math response from OpenAI")
	}

	return "🧮 [Math Agent] " + completion.Choices[0].Message.Content, nil
}

func (a *MathAgent) GetName() string {
	return "Math"
}

func (a *MathAgent) GetDescription() string {
	return "Specialized agent for mathematical calculations, problems, and concepts"
}
