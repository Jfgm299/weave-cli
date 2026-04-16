package app

import (
	"errors"
	"fmt"
)

var (
	ErrNotInProjectRoot   = errors.New("not in a project root")
	ErrInvalidConfig      = errors.New("invalid config")
	ErrUnsafeMutationPath = errors.New("unsafe mutation path")
)

func WrapInvalidConfig(err error) error {
	if err == nil {
		return ErrInvalidConfig
	}
	return fmt.Errorf("%w: %v. See %s (%s)", ErrInvalidConfig, err, DocsPathConfig, DocsURL(DocsPathConfig))
}
