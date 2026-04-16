package fsops

import (
	"context"
	"os"
	"path/filepath"
)

type Engine struct{}

func (Engine) Apply(_ context.Context, ops []Operation) error {
	for _, op := range ops {
		switch op.Type {
		case OpEnsureDir:
			if err := os.MkdirAll(op.Path, 0o755); err != nil {
				return err
			}
		case OpWriteFile:
			if err := os.MkdirAll(filepath.Dir(op.Path), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(op.Path, op.Content, 0o644); err != nil {
				return err
			}
		case OpCreateLink:
			if err := os.MkdirAll(filepath.Dir(op.Path), 0o755); err != nil {
				return err
			}
			if err := os.RemoveAll(op.Path); err != nil {
				return err
			}
			if err := os.Symlink(op.Target, op.Path); err != nil {
				return err
			}
		case OpRemovePath:
			if err := os.RemoveAll(op.Path); err != nil {
				return err
			}
		}
	}

	return nil
}
