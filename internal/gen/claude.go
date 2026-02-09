package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ClaudeClient obsługuje komunikację z Claude API
type ClaudeClient struct {
	apiKey  string
	baseURL string
}

// NewClaudeClient tworzy nowego klienta
func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY nie jest ustawiony")
	}

	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1/messages",
	}, nil
}

// ClaudeRequest to struktura requestu do Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage to pojedyncza wiadomość w konwersacji
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse to odpowiedź z Claude API
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// GenerateTestCases wysyła prompt do Claude i dostaje test cases
func (c *ClaudeClient) GenerateTestCases(description string) ([]TestCase, error) {
	// 1. Przygotuj prompt (funkcja z provider_common.go)
	prompt := buildPrompt(description) // ← używa wspólnej funkcji

	// 2. Utwórz request
	request := ClaudeRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 3. Konwertuj struct → JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("błąd przy tworzeniu JSON: %w", err)
	}

	// 4. Utwórz HTTP request
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("błąd przy tworzeniu requestu: %w", err)
	}

	// 5. Dodaj headery
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// 6. Wyślij request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("błąd przy wysyłaniu requestu: %w", err)
	}
	defer resp.Body.Close()

	// 7. Sprawdź status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API zwrócił błąd %d: %s", resp.StatusCode, string(body))
	}

	// 8. Przeczytaj response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd przy odczycie response: %w", err)
	}

	// 9. JSON → Go struct
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("błąd przy parsowaniu JSON: %w", err)
	}

	// 10. Wyciągnij tekst z response
	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("Claude nie zwrócił żadnej treści")
	}
	generatedText := claudeResp.Content[0].Text

	// 11. Sparsuj wygenerowane test cases (funkcja z provider_common.go)
	testCases, err := parseTestCases(generatedText) // ← używa wspólnej funkcji
	if err != nil {
		return nil, fmt.Errorf("błąd przy parsowaniu test cases: %w", err)
	}

	return testCases, nil
}
