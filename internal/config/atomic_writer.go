package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type AtomicFileWriter struct {
	Path     string
	renameFn func(oldpath, newpath string) error
}

func (w AtomicFileWriter) Write(cfg Config) error {
	b, err := MarshalDeterministic(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(w.Path), 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(w.Path), ".weave-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	cleanup := func() {
		_ = os.Remove(tmpPath)
	}

	if _, err := tmp.Write(append(b, '\n')); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}

	if err := tmp.Chmod(0o644); err != nil {
		_ = tmp.Close()
		cleanup()
		return err
	}

	if err := tmp.Close(); err != nil {
		cleanup()
		return err
	}

	rename := w.renameFn
	if rename == nil {
		rename = os.Rename
	}

	if err := rename(tmpPath, w.Path); err != nil {
		cleanup()
		return fmt.Errorf("atomic rename failed: %w", err)
	}

	return nil
}
