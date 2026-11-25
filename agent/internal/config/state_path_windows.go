//go:build windows

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func defaultStateFilePath() (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")

	if localAppData == "" {
		return "", fmt.Errorf("LOCALAPPDATA is not set")
	}

	return filepath.Join(
		localAppData,
		"SignalLens",
		"agent",
		"signallens-state.json",
	), nil
}
