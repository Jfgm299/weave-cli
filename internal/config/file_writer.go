package config

import "os"

type FileWriter struct {
	Path string
}

func (w FileWriter) Write(cfg Config) error {
	if _, err := os.Stat(w.Path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	b, err := MarshalDeterministic(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(w.Path, append(b, '\n'), 0o644)
}
