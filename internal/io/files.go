// Package io provides utility function relating to file I/O
package io

import (
	"errors"
	"path/filepath"
	"runtime"
)

func getRootfolder() (string, error) {
	_, caller, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("could not get calling folder")
	}
	rootPath := filepath.Join(filepath.Dir(caller), "../..")
	return rootPath, nil
}
