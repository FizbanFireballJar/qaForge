package platform

import (
	"os"
	"path/filepath"
)

// GetConfigDir zwraca ścieżkę do katalogu konfiguracji QA Forge
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".qaforge"), nil
}

// GetTestsDir zwraca ścieżkę do katalogu testów
func GetTestsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "tests"), nil
}