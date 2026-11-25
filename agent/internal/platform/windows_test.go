//go:build windows

package platform

import (
	"strings"
	"testing"
)

func TestGetOSInfo(t *testing.T) {
	info, err := GetOSInfo()
	if err != nil {
		t.Fatalf("GetOSInfo() failed: %v", err)
	}

	if strings.TrimSpace(info.Name) == "" {
		t.Fatal("expected OS name to be populated")
	}

	if strings.TrimSpace(info.Version) == "" {
		t.Fatal("expected OS version to be populated")
	}

	t.Logf("OS Name: %s", info.Name)
	t.Logf("OS Version: %s", info.Version)
}
