package gen

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// buildPrompt tworzy prompt dla AI (Claude lub Mistral)
func buildPrompt(description string) string {
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

// parseTestCases parsuje JSON z test cases od AI
func parseTestCases(text string) ([]TestCase, error) {
	// Czasem AI dodaje markdown backticks, usuń je
	text = cleanJSON(text)

	// Struktura tymczasowa do parsowania JSON
	var rawCases []struct {
		Title          string   `json:"title"`
		Steps          []Step   `json:"steps"`
		ExpectedResult string   `json:"expected_result"`
		Tags           []string `json:"tags"`
	}

	if err := json.Unmarshal([]byte(text), &rawCases); err != nil {
		return nil, fmt.Errorf("niepoprawny JSON od AI: %w\n\nOtrzymany tekst:\n%s", err, text)
	}

	// Konwertuj na nasze TestCase ze wszystkimi polami
	testCases := make([]TestCase, 0, len(rawCases))
	for _, raw := range rawCases {
		tc := TestCase{
			ID:             "", // Puste! SaveTestCases() przypisze ID
			Title:          raw.Title,
			Tags:           raw.Tags,
			Status:         "needs-review",
			Type:           "manual",
			Steps:          raw.Steps,
			ExpectedResult: raw.ExpectedResult,
			Created:        time.Now(),
			Source:         "qaforge gen create",
			Runs:           []string{},
		}
		testCases = append(testCases, tc)
	}

	return testCases, nil
}

// cleanJSON usuwa markdown backticks i trailing commas które AI czasem dodaje
func cleanJSON(text string) string {
	// Usuń markdown backticks
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	// Usuń trailing commas przed ] i } (Mistral czasem je generuje)
	// Przykład: {"key": "value",} → {"key": "value"}
	// Przykład: ["item1", "item2",] → ["item1", "item2"]
	trailingCommaRegex := regexp.MustCompile(`,(\s*[\]}])`)
	text = trailingCommaRegex.ReplaceAllString(text, "$1")

	return text
}
