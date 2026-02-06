package gen

import (
	tea "github.com/charmbracelet/bubbletea"
)

// reviewModel to stan aplikacji review UI
type reviewModel struct {
	testCases    []TestCase // Wszystkie test cases do review
	currentIndex int        // Który test case oglądamy (0-based)
	width        int        // Szerokość terminala
	height       int        // Wysokość terminala
	quitting     bool       // Czy user nacisnął 'q'?
	editingTags  bool       // Czy jesteśmy w trybie edycji tagów?
	tagInput     string     // Buffer dla edycji tagów
}

// initialModel tworzy początkowy stan
func initialModel(testCases []TestCase) reviewModel {
	return reviewModel{
		testCases:    testCases,
		currentIndex: 0,
		width:        80, // Domyślnie 80 kolumn
		height:       24, // Domyślnie 24 linie
		quitting:     false,
		editingTags:  false,
		tagInput:     "",
	}
}

// Init to funkcja wywoływana przy starcie (opcjonalna)
// Zwraca Command do wykonania (lub nil jeśli nic)
func (m reviewModel) Init() tea.Cmd {
	return nil
}

// currentTestCase zwraca aktualnie wyświetlany test case
func (m reviewModel) currentTestCase() *TestCase {
	if len(m.testCases) == 0 {
		return nil
	}
	return &m.testCases[m.currentIndex]
}

// hasNext sprawdza czy jest następny test case
func (m reviewModel) hasNext() bool {
	return m.currentIndex < len(m.testCases)-1
}

// hasPrev sprawdza czy jest poprzedni test case
func (m reviewModel) hasPrev() bool {
	return m.currentIndex > 0
}
