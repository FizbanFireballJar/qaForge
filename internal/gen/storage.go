package gen

import (
	"fmt"
	"os"
	"path/filepath"

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
