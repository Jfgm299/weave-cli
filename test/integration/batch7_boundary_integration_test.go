package integration

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBatch7_Integration_AppLayerIsReusableOutsideCLI(t *testing.T) {
	t.Parallel()

	repoRoot := repoRootFromIntegrationCaller(t)
	appSources := readGoSources(t, filepath.Join(repoRoot, "internal", "app"))

	if strings.Contains(appSources, "internal/cli") {
		t.Fatalf("internal/app must not import internal/cli to keep domain reusable for future TUI adapters")
	}
}

func TestBatch7_Integration_GoModuleStaysFreeOfTUIRuntimeDeps(t *testing.T) {
	t.Parallel()

	repoRoot := repoRootFromIntegrationCaller(t)
	goModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	goMod := string(goModBytes)
	for _, disallowed := range []string{
		"github.com/charmbracelet/bubbletea",
		"github.com/charmbracelet/bubbles",
		"github.com/charmbracelet/lipgloss",
		"github.com/rivo/tview",
		"github.com/jroimartin/gocui",
	} {
		if strings.Contains(goMod, disallowed) {
			t.Fatalf("go.mod must not include TUI runtime dependency %q", disallowed)
		}
	}
}

func repoRootFromIntegrationCaller(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve caller path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}

func readGoSources(t *testing.T, root string) string {
	t.Helper()

	var b strings.Builder
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		b.WriteString("\n")
		b.Write(content)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}

	return b.String()
}
