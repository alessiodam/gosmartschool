package assets

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed chromedriver chromedriver.exe
var content embed.FS

func ExtractFile(name string) (string, error) {
	data, err := content.ReadFile(name)
	if err != nil {
		return "", err
	}

	dir, err := os.MkdirTemp("", "gosmartschool")
	if err != nil {
		return "", err
	}

	extractedFile := filepath.Join(dir, name)
	err = os.WriteFile(extractedFile, data, 0755)
	if err != nil {
		return "", err
	}

	return extractedFile, nil
}
