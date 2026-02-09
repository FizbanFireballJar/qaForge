package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"qaforge/internal/platform"

	"gopkg.in/yaml.v3"
)

// SaveTestCase zapisuje test case jako plik YAML
func SaveTestCase(tc TestCase) error {
	// 1. Uzyskaj katalog testów (np. ~/.qaforge/tests)
	testsDir, err := platform.GetTestsDir()
	if err != nil {
		return fmt.Errorf("nie można uzyskać katalogu testów: %w", err)
	}

	// 2. Utwórz katalog jeśli nie istnieje
	// os.MkdirAll tworzy wszystkie katalogi w ścieżce
	// 0755 to uprawnienia (rwxr-xr-x)
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		return fmt.Errorf("nie można utworzyć katalogu: %w", err)
	}

	// 3. Przygotuj ścieżkę do pliku
	// filepath.Join łączy ścieżki używając odpowiedniego separatora (/ lub \)
	filename := fmt.Sprintf("%s.yaml", tc.ID) // TC-001.yaml
	filePath := filepath.Join(testsDir, filename)

	// 4. Konwertuj TestCase → YAML ([]byte)
	data, err := yaml.Marshal(&tc)
	if err != nil {
		return fmt.Errorf("błąd przy konwersji do YAML: %w", err)
	}

	// 5. Zapisz do pliku
	// 0644 to uprawnienia pliku (rw-r--r--)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("błąd przy zapisie pliku: %w", err)
	}

	return nil
}

// SaveTestCases zapisuje listę test cases
func SaveTestCases(testCases []TestCase) error {
	// 1. Znajdź następny wolny ID
	nextID := getNextTestCaseID()

	// 2. Przypisz ID zaczynając od nextID
	for i := range testCases {
		testCases[i].ID = fmt.Sprintf("TC-%03d", nextID+i)
	}

	// 3. Zapisz wszystkie test cases
	for _, tc := range testCases {
		if err := SaveTestCase(tc); err != nil {
			return err
		}
	}

	return nil
}

// LoadTestCasesWithStatus ładuje test cases z danym statusem
func LoadTestCasesWithStatus(status string) ([]TestCase, error) {
	testsDir, err := platform.GetTestsDir()
	if err != nil {
		return nil, err
	}

	// Przeczytaj wszystkie pliki .yaml
	files, err := filepath.Glob(filepath.Join(testsDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	testCases := []TestCase{}
	for _, file := range files {
		tc, err := LoadTestCase(file)
		if err != nil {
			continue // Skip invalid files
		}

		// Filtruj po statusie
		if tc.Status == status {
			testCases = append(testCases, tc)
		}
	}

	return testCases, nil
}

// LoadTestCase ładuje pojedynczy test case z pliku
func LoadTestCase(filePath string) (TestCase, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TestCase{}, err
	}

	var tc TestCase
	if err := yaml.Unmarshal(data, &tc); err != nil {
		return TestCase{}, err
	}

	return tc, nil
}

// SaveReviewedTestCases zapisuje zmodyfikowane test cases
func SaveReviewedTestCases(testCases []TestCase) error {
	for _, tc := range testCases {
		if err := SaveTestCase(tc); err != nil {
			return err
		}
	}
	return nil
}

// GetTestsDir zwraca ścieżkę do katalogu testów
func GetTestsDir() (string, error) {
	return platform.GetTestsDir()
}

// LoadAllTestCases ładuje wszystkie test cases z katalogu
func LoadAllTestCases() ([]TestCase, error) {
	testsDir, err := platform.GetTestsDir()
	if err != nil {
		return nil, err
	}

	// Przeczytaj wszystkie pliki .yaml
	files, err := filepath.Glob(filepath.Join(testsDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	testCases := []TestCase{}
	for _, file := range files {
		tc, err := LoadTestCase(file)
		if err != nil {
			// Skip invalid files, ale loguj warning
			fmt.Fprintf(os.Stderr, "⚠️  Pominięto niepoprawny plik: %s\n", file)
			continue
		}
		testCases = append(testCases, tc)
	}

	// Sortuj po ID
	sort.Slice(testCases, func(i, j int) bool {
		return testCases[i].ID < testCases[j].ID
	})

	return testCases, nil
}

// getNextTestCaseID zwraca następny wolny numer ID
func getNextTestCaseID() int {
	// Wczytaj wszystkie istniejące test cases
	allTests, err := LoadAllTestCases()
	if err != nil || len(allTests) == 0 {
		return 1 // Jeśli brak testów, zacznij od 1
	}

	// Znajdź najwyższy numer ID
	maxID := 0
	for _, tc := range allTests {
		// Wyciągnij numer z "TC-001" → 1
		var num int
		fmt.Sscanf(tc.ID, "TC-%d", &num)
		if num > maxID {
			maxID = num
		}
	}

	// Zwróć następny wolny numer
	return maxID + 1
}
