package run

import (
	"fmt"
	"strings"

	"qaforge/internal/gen"
)

// TreeOptions to opcje dla drzewa testów
type TreeOptions struct {
	FilterTag    string // Filtruj po tagu (opcjonalnie)
	FilterStatus string // Filtruj po statusie (opcjonalnie)
	ShowStats    bool   // Czy pokazywać statystyki
}

// TestTree reprezentuje drzewo testów
type TestTree struct {
	TestCases []gen.TestCase
	ByTag     map[string][]gen.TestCase
	Stats     TreeStats
	ShowStats bool // Czy pokazywać statystyki
}

// TreeStats to statystyki testów
type TreeStats struct {
	Total       int
	Approved    int
	NeedsWork   int
	NeedsReview int
	Rejected    int
}

// BuildTree buduje drzewo testów z opcjami filtrowania
func BuildTree(options TreeOptions) (*TestTree, error) {
	// 1. Załaduj wszystkie test cases
	allTests, err := gen.LoadAllTestCases()
	if err != nil {
		return nil, fmt.Errorf("błąd przy ładowaniu testów: %w", err)
	}

	// 2. Filtruj test cases
	filtered := filterTestCases(allTests, options)

	// 3. Grupuj po tagach
	byTag := groupByTag(filtered)

	// 4. Oblicz statystyki
	stats := calculateStats(filtered)

	return &TestTree{
		TestCases: filtered,
		ByTag:     byTag,
		Stats:     stats,
		ShowStats: options.ShowStats,
	}, nil
}

// filterTestCases filtruje test cases według opcji
func filterTestCases(tests []gen.TestCase, options TreeOptions) []gen.TestCase {
	filtered := []gen.TestCase{}

	for _, tc := range tests {
		// Filtruj po statusie
		if options.FilterStatus != "" && tc.Status != options.FilterStatus {
			continue
		}

		// Filtruj po tagu
		if options.FilterTag != "" {
			hasTag := false
			for _, tag := range tc.Tags {
				// Usuń # z początku tagu jeśli jest
				cleanTag := strings.TrimPrefix(tag, "#")
				filterTag := strings.TrimPrefix(options.FilterTag, "#")

				if cleanTag == filterTag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		filtered = append(filtered, tc)
	}

	return filtered
}

// groupByTag grupuje test cases po tagach
func groupByTag(tests []gen.TestCase) map[string][]gen.TestCase {
	byTag := make(map[string][]gen.TestCase)

	for _, tc := range tests {
		// Każdy test case może mieć wiele tagów
		for _, tag := range tc.Tags {
			// Normalizuj tag (dodaj # jeśli nie ma)
			if !strings.HasPrefix(tag, "#") {
				tag = "#" + tag
			}

			byTag[tag] = append(byTag[tag], tc)
		}
	}

	return byTag
}

// calculateStats oblicza statystyki
func calculateStats(tests []gen.TestCase) TreeStats {
	stats := TreeStats{
		Total: len(tests),
	}

	for _, tc := range tests {
		switch tc.Status {
		case "approved":
			stats.Approved++
		case "needs-work":
			stats.NeedsWork++
		case "needs-review":
			stats.NeedsReview++
		case "rejected":
			stats.Rejected++
		}
	}

	return stats
}
