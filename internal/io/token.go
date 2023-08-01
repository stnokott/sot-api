package io

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// tokenFilename contains name of file with RAT token from seaofthieves.com
const tokenFilename string = "rat.txt"

// ReadToken returns the RAT token from the token file in the project root
func ReadToken() (string, error) {
	rootFolder, err := getRootfolder()
	if err != nil {
		return "", err
	}
	tokenFilePath := filepath.Join(rootFolder, tokenFilename)
	tokenBytes, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("could not read token file: %w", err)
	}
	return strings.TrimSpace(string(tokenBytes)), nil
}
