package gen

import tea "github.com/charmbracelet/bubbletea"

// ReviewModel to publiczna wersja reviewModel
type ReviewModel struct {
	reviewModel
}

// InitialReviewModel tworzy początkowy model
func InitialReviewModel(testCases []TestCase) ReviewModel {
	return ReviewModel{
		reviewModel: initialModel(testCases),
	}
}

// TestCases zwraca test cases z modelu
func (m ReviewModel) TestCases() []TestCase {
	return m.testCases
}

// Implementacja tea.Model interface
func (m ReviewModel) Init() tea.Cmd {
	return m.reviewModel.Init()
}

func (m ReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := m.reviewModel.Update(msg)
	m.reviewModel = updatedModel.(reviewModel)
	return m, cmd
}

func (m ReviewModel) View() string {
	return m.reviewModel.View()
}
