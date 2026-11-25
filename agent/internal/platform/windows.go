//go:build windows

package platform

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type OSInfo struct {
	Name    string
	Version string
}

func GetOSInfo() (OSInfo, error) {
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return OSInfo{}, fmt.Errorf("open Windows version registry key: %w", err)
	}
	defer key.Close()

	productName, _, err := key.GetStringValue("ProductName")
	if err != nil {
		return OSInfo{}, fmt.Errorf("read Windows product name: %w", err)
	}

	displayVersion, _, err := key.GetStringValue("DisplayVersion")
	if err != nil {
		return OSInfo{}, fmt.Errorf("read Windows display version: %w", err)
	}

	currentBuild, _, err := key.GetStringValue("CurrentBuild")
	if err != nil {
		return OSInfo{}, fmt.Errorf("read Windows current build: %w", err)
	}

	version := displayVersion

	if currentBuild != "" {
		version = fmt.Sprintf("%s (Build %s)", displayVersion, currentBuild)
	}

	return OSInfo{
		Name:    productName,
		Version: version,
	}, nil
}
