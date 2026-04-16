package app

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDocsReferencePaths_ExistInRepository(t *testing.T) {
	t.Parallel()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve caller path")
	}

	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	required := []string{
		DocsPathDoctor,
		DocsPathConfig,
		DocsPathTransactions,
		DocsPathProviders,
		DocsPathMigration,
	}

	for _, rel := range required {
		abs := filepath.Join(repoRoot, rel)
		if _, err := os.Stat(abs); err != nil {
			t.Fatalf("missing docs reference target %s: %v", rel, err)
		}
	}
}
