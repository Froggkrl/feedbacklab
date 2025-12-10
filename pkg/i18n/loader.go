package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// LoadTranslations loads all translation files from the specified directory into the bundle.
func LoadTranslations(bundle *goi18n.Bundle, localesDir string) error {
	if bundle == nil {
		return fmt.Errorf("bundle is nil")
	}

	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		return fmt.Errorf("locales directory does not exist: %s", localesDir)
	}

	entries, err := os.ReadDir(localesDir)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".json") {
			continue
		}

		filePath := filepath.Join(localesDir, filename)
		_, err := bundle.LoadMessageFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to load translation file %s: %w", filename, err)
		}
	}

	return nil
}
