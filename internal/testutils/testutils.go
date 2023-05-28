package testutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Reads the file at `path` then returns its content
func TestAccExampleFile(t *testing.T, path string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("unable to get current file")
	}

	example, err := os.ReadFile(
		filepath.Join(
			filepath.Dir(currentFile),
			"..",
			"..",
			"examples",
			path,
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	return string(example)
}
