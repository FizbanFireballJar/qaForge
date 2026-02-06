package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ClaudeClient obsługuje komunikację z Claude API
type ClaudeClient struct {
	apiKey  string
	baseURL string
}

// NewClaudeClient tworzy nowego klienta
// Konstruktor w Go to zwykła funkcja zwracająca strukturę
func NewClaudeClient() (*ClaudeClient, error) {
	// Pobieramy API key ze zmiennej środowiskowej
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY nie jest ustawiony")
	}

	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1/messages",
	}, nil
}

// GenerateTestCases wysyła prompt do Claude i dostaje test cases
func (c *ClaudeClient) GenerateTestCases(description string) ([]TestCase, error) {
	// 1. Przygotuj prompt
	prompt := c.buildPrompt(description)

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
	// json.Marshal() zamienia Go struct na []byte z JSON
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
	defer resp.Body.Close() // Zawsze zamykaj Body!

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

	// 11. Sparsuj wygenerowane test cases
	testCases, err := c.parseTestCases(generatedText)
	if err != nil {
		return nil, fmt.Errorf("błąd przy parsowaniu test cases: %w", err)
	}

	return testCases, nil
}

// buildPrompt tworzy prompt dla Claude
func (c *ClaudeClient) buildPrompt(description string) string {
	return fmt.Sprintf(`Jesteś ekspertem QA. Na podstawie poniższego opisu funkcjonalności wygeneruj przypadki testowe.

OPIS FUNKCJONALNOŚCI:
%s

Wygeneruj przypadki testowe w formacie JSON:
[
  {
    "title": "Tytuł testu",
    "steps": [
      {"action": "Krok 1", "expected": "Oczekiwany rezultat 1"},
      {"action": "Krok 2", "expected": "Oczekiwany rezultat 2"}
    ],
    "expected_result": "Końcowy rezultat",
    "tags": ["tag1", "tag2"]
  }
]

Wygeneruj 3-7 przypadków testowych obejmujących happy path, edge cases i negative scenarios.
Zwróć TYLKO JSON, bez żadnego dodatkowego tekstu.`, description)
}

// parseTestCases parsuje JSON z test cases od Claude
func (c *ClaudeClient) parseTestCases(text string) ([]TestCase, error) {
	// Czasem Claude dodaje markdown backticks, usuń je
	text = cleanJSON(text)

	// Struktura tymczasowa do parsowania JSON od Claude
	var rawCases []struct {
		Title          string   `json:"title"`
		Steps          []Step   `json:"steps"`
		ExpectedResult string   `json:"expected_result"`
		Tags           []string `json:"tags"`
	}

	if err := json.Unmarshal([]byte(text), &rawCases); err != nil {
		return nil, fmt.Errorf("niepoprawny JSON od Claude: %w", err)
	}

	// Konwertuj na nasze TestCase ze wszystkimi polami
	testCases := make([]TestCase, 0, len(rawCases))
	for i, raw := range rawCases {
		tc := TestCase{
			ID:             fmt.Sprintf("TC-%03d", i+1), // TC-001, TC-002, ...
			Title:          raw.Title,
			Tags:           raw.Tags,
			Status:         "needs-review", // Domyślnie do review
			Type:           "manual",
			Steps:          raw.Steps,
			ExpectedResult: raw.ExpectedResult,
			Created:        time.Now(),
			Source:         "qaforge gen create",
			Runs:           []string{}, // Pusta lista
		}
		testCases = append(testCases, tc)
	}

	return testCases, nil
}

// cleanJSON usuwa markdown backticks które Claude czasem dodaje
func cleanJSON(text string) string {
	// Usuń ```json ... ```
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	return strings.TrimSpace(text)
}
