package gen

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Lip Gloss styles
var (
	// Główna ramka
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	// Tytuł test case
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	// Tagi
	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	// Kroki
	stepStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Expected result
	expectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Italic(true)

	// Help bar
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	// Approve/Needs/Reject
	approveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	needsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true)

	rejectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)
)

// View renderuje UI
func (m reviewModel) View() string {
	// Jeśli quitting, pokaż goodbye message
	if m.quitting {
		return "\n✅ Zapisano zmiany. Do zobaczenia!\n\n"
	}

	// Jeśli brak test cases
	if len(m.testCases) == 0 {
		return "\n⚠️  Brak test cases do review.\n\n"
	}

	tc := m.currentTestCase()
	if tc == nil {
		return "\n❌ Błąd: nie można załadować test case.\n\n"
	}

	// Tryb edycji tagów
	if m.editingTags {
		return m.renderTagEditor(tc)
	}

	// Normalny widok
	return m.renderTestCase(tc)
}

// renderTestCase renderuje widok test case
func (m reviewModel) renderTestCase(tc *TestCase) string {
	var b strings.Builder

	// Header z numerem
	progress := fmt.Sprintf("[%d/%d]", m.currentIndex+1, len(m.testCases))
	header := fmt.Sprintf("Test Case Review %s", progress)
	b.WriteString(titleStyle.Render(header))
	b.WriteString("\n\n")

	// ID i Title
	b.WriteString(fmt.Sprintf("%s: %s\n\n", tc.ID, tc.Title))

	// Tags
	if len(tc.Tags) > 0 {
		b.WriteString("Tags: ")
		for _, tag := range tc.Tags {
			b.WriteString(tagStyle.Render(tag))
			b.WriteString(" ")
		}
		b.WriteString("\n\n")
	}

	// Steps
	b.WriteString("Steps:\n")
	for i, step := range tc.Steps {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, step.Action))
		b.WriteString(stepStyle.Render(fmt.Sprintf("     → %s\n", step.Expected)))
		b.WriteString("\n")
	}

	// Expected result
	b.WriteString(expectedStyle.Render(fmt.Sprintf("Expected: %s\n", tc.ExpectedResult)))
	b.WriteString("\n")

	// Status
	statusStr := ""
	switch tc.Status {
	case "approved":
		statusStr = approveStyle.Render("🟢 APPROVED")
	case "needs-work":
		statusStr = needsStyle.Render("🟡 NEEDS WORK")
	case "rejected":
		statusStr = rejectStyle.Render("🔴 REJECTED")
	default:
		statusStr = "⚪ NEEDS REVIEW"
	}
	b.WriteString(fmt.Sprintf("Status: %s\n\n", statusStr))

	// Wrap w box
	content := boxStyle.Render(b.String())

	// Help bar
	help := m.renderHelp()

	return content + "\n" + help + "\n"
}

// renderHelp renderuje help bar z klawiszami
func (m reviewModel) renderHelp() string {
	helps := []string{
		approveStyle.Render("[a] 🟢 Approve"),
		needsStyle.Render("[n] 🟡 Needs Work"),
		rejectStyle.Render("[r] 🔴 Reject"),
		"[t] Edit Tags",
		"[→/←] Next/Prev",
		"[q] Quit & Save",
	}

	return helpStyle.Render(strings.Join(helps, "   "))
}

// renderTagEditor renderuje widok edycji tagów
func (m reviewModel) renderTagEditor(tc *TestCase) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Edit Tags"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("%s: %s\n\n", tc.ID, tc.Title))
	b.WriteString("Enter tags separated by spaces:\n\n")
	b.WriteString(tagStyle.Render(m.tagInput))
	b.WriteString("█") // Cursor
	b.WriteString("\n\n")

	help := helpStyle.Render("[Enter] Save   [Esc] Cancel")

	content := boxStyle.Render(b.String())
	return content + "\n" + help + "\n"
}
