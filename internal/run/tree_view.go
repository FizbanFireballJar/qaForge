package run

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"qaforge/internal/gen"
)

// Lip Gloss styles
var (
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	tagStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	approvedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	needsWorkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("11"))

	needsReviewStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))

	rejectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	statsStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, false, false).
		Padding(1, 0)
)

// Render renderuje drzewo testów jako string
func (t *TestTree) Render() string {
	var b strings.Builder

	// Header
	title := fmt.Sprintf("📦 All Tests (%d total)", t.Stats.Total)
	b.WriteString(headerStyle.Render(title))
	b.WriteString("\n\n")

	// Jeśli brak testów
	if len(t.TestCases) == 0 {
		b.WriteString("⚠️  Brak test cases do wyświetlenia.\n")
		return b.String()
	}

	// Sortuj tagi alfabetycznie
	tags := make([]string, 0, len(t.ByTag))
	for tag := range t.ByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	// Renderuj każdą grupę tagów
	for _, tag := range tags {
		tests := t.ByTag[tag]

		// Tag header
		tagHeader := fmt.Sprintf("%s (%d tests)", tag, len(tests))
		b.WriteString(tagStyle.Render(tagHeader))
		b.WriteString("\n")

		// Sortuj testy w grupie po ID
		sort.Slice(tests, func(i, j int) bool {
			return tests[i].ID < tests[j].ID
		})

		// Renderuj każdy test w grupie
		for i, tc := range tests {
			isLast := i == len(tests)-1
			b.WriteString(renderTestCase(tc, isLast))
		}

		b.WriteString("\n")
	}

	// Statystyki (tylko jeśli ShowStats jest true)
	if t.ShowStats {
		b.WriteString(renderStats(t.Stats))
	}

	return b.String()
}

// renderTestCase renderuje pojedynczy test case w drzewie
func renderTestCase(tc gen.TestCase, isLast bool) string {
	// Tree characters
	prefix := "├─ "
	if isLast {
		prefix = "└─ "
	}

	// Status icon i kolor
	icon, style := getStatusIconAndStyle(tc.Status)

	// Stylizuj tylko ikonę, nie całą linię (żeby nie psuć formatowania)
	styledIcon := style.Render(icon)

	// Format: ├─ ✅ TC-001: Title (status)
	return fmt.Sprintf("%s%s %s: %s (%s)\n",
		prefix,
		styledIcon,
		tc.ID,
		tc.Title,
		tc.Status,
	)
}

// getStatusIconAndStyle zwraca ikonę i styl dla statusu
func getStatusIconAndStyle(status string) (string, lipgloss.Style) {
	switch status {
	case "approved":
		return "✅", approvedStyle
	case "needs-work":
		return "🟡", needsWorkStyle
	case "needs-review":
		return "⏳", needsReviewStyle
	case "rejected":
		return "🔴", rejectedStyle
	default:
		return "⚪", lipgloss.NewStyle()
	}
}

// renderStats renderuje statystyki
func renderStats(stats TreeStats) string {
	var b strings.Builder

	b.WriteString("─────────────────────────────────────────────────────\n")
	b.WriteString("Statistics:\n")

	// Oblicz procenty
	approvedPct := 0
	needsWorkPct := 0
	needsReviewPct := 0
	rejectedPct := 0

	if stats.Total > 0 {
		approvedPct = (stats.Approved * 100) / stats.Total
		needsWorkPct = (stats.NeedsWork * 100) / stats.Total
		needsReviewPct = (stats.NeedsReview * 100) / stats.Total
		rejectedPct = (stats.Rejected * 100) / stats.Total
	}

	// Renderuj każdą linię statystyk
	if stats.Approved > 0 {
		line := fmt.Sprintf("  ✅ Approved:      %d (%d%%)\n", stats.Approved, approvedPct)
		b.WriteString(approvedStyle.Render(line))
	}

	if stats.NeedsWork > 0 {
		line := fmt.Sprintf("  🟡 Needs Work:    %d (%d%%)\n", stats.NeedsWork, needsWorkPct)
		b.WriteString(needsWorkStyle.Render(line))
	}

	if stats.NeedsReview > 0 {
		line := fmt.Sprintf("  ⏳ Needs Review:  %d (%d%%)\n", stats.NeedsReview, needsReviewPct)
		b.WriteString(needsReviewStyle.Render(line))
	}

	if stats.Rejected > 0 {
		line := fmt.Sprintf("  🔴 Rejected:      %d (%d%%)\n", stats.Rejected, rejectedPct)
		b.WriteString(rejectedStyle.Render(line))
	}

	b.WriteString(fmt.Sprintf("  Total:           %d\n", stats.Total))

	return statsStyle.Render(b.String())
}
