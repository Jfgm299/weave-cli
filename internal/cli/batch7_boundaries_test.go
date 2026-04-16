package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var disallowedTUIModules = []string{
	"github.com/charmbracelet/bubbletea",
	"github.com/charmbracelet/bubbles",
	"github.com/charmbracelet/lipgloss",
	"github.com/rivo/tview",
	"github.com/jroimartin/gocui",
}

func TestArchitectureBoundary_AppLayerHasNoCLINorTUIImports(t *testing.T) {
	t.Parallel()

	repoRoot := batch7RepoRootFromCaller(t)
	content := readGoSourceTree(t, filepath.Join(repoRoot, "internal", "app"))

	if strings.Contains(content, "github.com/Jfgm299/weave-cli/internal/cli") {
		t.Fatalf("app layer must not import internal/cli")
	}

	for _, module := range disallowedTUIModules {
		if strings.Contains(content, module) {
			t.Fatalf("app layer must not import TUI module %q", module)
		}
	}
}

func TestRuntimeBoundary_V1HasNoTUIRuntimeDependency(t *testing.T) {
	t.Parallel()

	repoRoot := batch7RepoRootFromCaller(t)

	goMod := readFile(t, filepath.Join(repoRoot, "go.mod"))
	for _, module := range disallowedTUIModules {
		if strings.Contains(goMod, module) {
			t.Fatalf("go.mod must not require TUI runtime dependency %q", module)
		}
	}

	runtimeTree := readFile(t,
		filepath.Join(repoRoot, "cmd", "weave", "main.go"),
	)
	runtimeTree += readGoSourceTree(t, filepath.Join(repoRoot, "internal"))
	for _, module := range disallowedTUIModules {
		if strings.Contains(runtimeTree, module) {
			t.Fatalf("runtime source must not import TUI dependency %q", module)
		}
	}
}

func TestInstallFlow_OneCommandScriptAndDocsAreAligned(t *testing.T) {
	t.Parallel()

	repoRoot := batch7RepoRootFromCaller(t)

	installScript := readFile(t, filepath.Join(repoRoot, "scripts", "install.sh"))
	if info, err := os.Stat(filepath.Join(repoRoot, "scripts", "install.sh")); err != nil {
		t.Fatalf("install script must exist: %v", err)
	} else if info.Mode()&0o111 == 0 {
		t.Fatalf("install script must be executable")
	}
	if !strings.Contains(installScript, "go install") {
		t.Fatalf("install script must install weave via go install")
	}
	if !strings.Contains(installScript, "./cmd/weave") {
		t.Fatalf("install script must target ./cmd/weave")
	}

	readme := readFile(t, filepath.Join(repoRoot, "README.md"))
	installDoc := readFile(t, filepath.Join(repoRoot, "docs", "reference", "install.md"))

	for _, required := range []string{"./scripts/install.sh", "weave --version", "weave forge", "weave doctor"} {
		if !strings.Contains(readme, required) {
			t.Fatalf("README must document %q", required)
		}
		if !strings.Contains(installDoc, required) {
			t.Fatalf("docs/reference/install.md must document %q", required)
		}
	}
}

func batch7RepoRootFromCaller(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to resolve caller path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}

func readFile(t *testing.T, paths ...string) string {
	t.Helper()

	if len(paths) == 0 {
		t.Fatalf("readFile requires at least one path")
	}

	if len(paths) == 1 {
		b, err := os.ReadFile(paths[0])
		if err != nil {
			t.Fatalf("read %s: %v", paths[0], err)
		}
		return string(b)
	}

	b := strings.Builder{}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		b.WriteString("\n")
		b.Write(content)
	}

	return b.String()
}

func readGoSourceTree(t *testing.T, root string) string {
	t.Helper()

	b := strings.Builder{}
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
