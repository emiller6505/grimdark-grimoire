package service

import (
	"os"
	"testing"
)

// getTestDataDir returns the test data directory path
func getTestDataDir(t *testing.T) string {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	return dataDir
}

