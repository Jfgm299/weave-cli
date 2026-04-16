package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func ensureOpsWithinRoot(root string, ops []fsops.Operation) error {
	for _, op := range ops {
		if op.Path == "" {
			continue
		}
		if err := ensurePathWithinRoot(root, op.Path); err != nil {
			return err
		}
	}
	return nil
}

func ensurePathWithinRoot(root string, p string) error {
	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("%w: cannot resolve project root %q: %v", ErrUnsafeMutationPath, root, err)
	}
	cleanPath, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("%w: cannot resolve mutation path %q: %v", ErrUnsafeMutationPath, p, err)
	}

	if cleanPath == cleanRoot {
		return nil
	}

	prefix := cleanRoot + string(filepath.Separator)
	if strings.HasPrefix(cleanPath, prefix) {
		return nil
	}

	return fmt.Errorf("%w: %s is outside project root %s. Keep all mutation targets under the repository root", ErrUnsafeMutationPath, cleanPath, cleanRoot)
}
