package gen

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update reaguje na wiadomości (keyboard input, window resize, etc.)
func (m reviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Naciśnięcie klawisza
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Zmiana rozmiaru okna terminala
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKeyPress obsługuje naciśnięcia klawiszy
func (m reviewModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Jeśli edytujemy tagi, obsłuż input tagów
	if m.editingTags {
		return m.handleTagInput(msg)
	}

	// Normalne klawisze (nie edycja tagów)
	switch msg.String() {

	// Quit
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	// Approve (🟢)
	case "a":
		if tc := m.currentTestCase(); tc != nil {
			tc.Status = "approved"
		}
		return m.moveNext()

	// Needs Work (🟡)
	case "n":
		if tc := m.currentTestCase(); tc != nil {
			tc.Status = "needs-work"
		}
		return m.moveNext()

	// Reject (🔴)
	case "r":
		if tc := m.currentTestCase(); tc != nil {
			tc.Status = "rejected"
		}
		return m.moveNext()

	// Edit tags
	case "t":
		if tc := m.currentTestCase(); tc != nil {
			m.editingTags = true
			// Załaduj obecne tagi do bufora
			m.tagInput = joinTags(tc.Tags)
		}
		return m, nil

	// Next test case
	case "right", "l":
		return m.moveNext()

	// Previous test case
	case "left", "h":
		return m.movePrev()
	}

	return m, nil
}

// handleTagInput obsługuje input w trybie edycji tagów
func (m reviewModel) handleTagInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	// Enter — zapisz tagi
	case "enter":
		if tc := m.currentTestCase(); tc != nil {
			tc.Tags = parseTags(m.tagInput)
		}
		m.editingTags = false
		m.tagInput = ""
		return m, nil

	// Escape — anuluj edycję
	case "esc":
		m.editingTags = false
		m.tagInput = ""
		return m, nil

	// Backspace
	case "backspace":
		if len(m.tagInput) > 0 {
			m.tagInput = m.tagInput[:len(m.tagInput)-1]
		}
		return m, nil

	// Dodaj znak do bufora
	default:
		m.tagInput += msg.String()
		return m, nil
	}
}

// moveNext przechodzi do następnego test case
func (m reviewModel) moveNext() (tea.Model, tea.Cmd) {
	if m.hasNext() {
		m.currentIndex++
	}
	return m, nil
}

// movePrev przechodzi do poprzedniego test case
func (m reviewModel) movePrev() (tea.Model, tea.Cmd) {
	if m.hasPrev() {
		m.currentIndex--
	}
	return m, nil
}

// Helper functions

func joinTags(tags []string) string {
	result := ""
	for i, tag := range tags {
		if i > 0 {
			result += " "
		}
		result += tag
	}
	return result
}

func parseTags(input string) []string {
	// Split po spacjach, usuń puste
	tags := strings.Fields(input)
	return tags
}
