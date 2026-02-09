package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MistralClient obsługuje komunikację z Ollama (Mistral Nemo)
type MistralClient struct {
	baseURL string
	model   string
}

// NewMistralClient tworzy klienta Ollama
func NewMistralClient() *MistralClient {
	return &MistralClient{
		baseURL: "http://localhost:11434/api/generate",
		model:   "mistral-nemo",
	}
}

// OllamaRequest to struktura requestu do Ollama
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse to odpowiedź z Ollama
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// GenerateTestCases wysyła prompt do Mistral i dostaje test cases
func (c *MistralClient) GenerateTestCases(description string) ([]TestCase, error) {
	// 1. Przygotuj prompt (ta sama funkcja co w claude.go - z provider_common.go)
	prompt := buildPrompt(description)

	// 2. Utwórz request
	request := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false, // Nie używamy streamingu
	}

	// 3. Konwertuj struct → JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("błąd przy tworzeniu JSON: %w", err)
	}

	// 4. Wyślij HTTP POST
	resp, err := http.Post(c.baseURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("błąd połączenia z Ollama (czy działa 'ollama serve'?): %w", err)
	}
	defer resp.Body.Close()

	// 5. Sprawdź status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama zwrócił błąd %d: %s", resp.StatusCode, string(body))
	}

	// 6. Przeczytaj response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd przy odczycie response: %w", err)
	}

	// 7. JSON → Go struct
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("błąd przy parsowaniu JSON: %w", err)
	}

	// 8. Sparsuj wygenerowane test cases (ta sama funkcja - z provider_common.go)
	testCases, err := parseTestCases(ollamaResp.Response)
	if err != nil {
		return nil, fmt.Errorf("błąd przy parsowaniu test cases: %w", err)
	}

	return testCases, nil
}
