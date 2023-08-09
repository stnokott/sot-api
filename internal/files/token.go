package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const tokenFilename string = "rat.txt"

// WriteToken writes the RAT token to file
func WriteToken(t string) error {
	tokenFilePath, err := getTokenFilepath()
	if err != nil {
		return err
	}
	return os.WriteFile(tokenFilePath, []byte(t), 0644)
}

// ReadToken returns the RAT token from the token file in the project root.
func ReadToken() (string, error) {
	tokenFilePath, err := getTokenFilepath()
	if err != nil {
		return "", err
	}
	tokenBytes, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("could not read token file: %w", err)
	}
	return strings.TrimSpace(string(tokenBytes)), nil
}

func getTokenFilepath() (string, error) {
	rootFolder, err := getRootfolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(rootFolder, tokenFilename), nil
}
