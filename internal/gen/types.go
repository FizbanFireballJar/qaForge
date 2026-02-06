package gen

import "time"

// TestCase reprezentuje pojedynczy przypadek testowy
// Każde pole ma:
// - nazwę (np. ID)
// - typ (np. string)
// - tag YAML (jak pole się nazywa w pliku YAML)
type TestCase struct {
	ID             string    `yaml:"id"`              // TC-001
	Title          string    `yaml:"title"`           // "Logowanie poprawnym hasłem"
	Tags           []string  `yaml:"tags"`            // ["logowanie", "P1"]
	Status         string    `yaml:"status"`          // "approved", "needs-review"
	Type           string    `yaml:"type"`            // "manual", "automated"
	Steps          []Step    `yaml:"steps"`           // Lista kroków
	ExpectedResult string    `yaml:"expected_result"` // Oczekiwany rezultat
	Created        time.Time `yaml:"created"`         // Kiedy utworzono
	Source         string    `yaml:"source"`          // Skąd pochodzi (opis/URL)
	Runs           []string  `yaml:"runs"`            // Historia test runów (puste na start)
}

// Step reprezentuje pojedynczy krok w teście
type Step struct {
	Action   string `yaml:"action"`   // "Otwórz stronę /login"
	Expected string `yaml:"expected"` // "Formularz się wyświetla"
}

// ClaudeRequest to struktura requestu do Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage to pojedyncza wiadomość w konwersacji
type ClaudeMessage struct {
	Role    string `json:"role"`    // "user" lub "assistant"
	Content string `json:"content"` // Treść wiadomości
}

// ClaudeResponse to odpowiedź z Claude API
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}
