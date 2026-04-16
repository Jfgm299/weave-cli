package app

import (
	"errors"
	"fmt"
)

var (
	ErrNotInProjectRoot = errors.New("not in a project root")
	ErrInvalidConfig    = errors.New("invalid config")
)

func WrapInvalidConfig(err error) error {
	if err == nil {
		return ErrInvalidConfig
	}
	return fmt.Errorf("%w: %v", ErrInvalidConfig, err)
}
