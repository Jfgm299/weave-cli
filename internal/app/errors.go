package app

import (
	"errors"
	"fmt"

	"github.com/Jfgm299/weave-cli/internal/config"
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
	if errors.Is(err, config.ErrOutdatedSchema) {
		return fmt.Errorf("%w: %v. See %s (%s)", ErrInvalidConfig, err, DocsPathMigration, DocsURL(DocsPathMigration))
	}
	return fmt.Errorf("%w: %v. See %s (%s)", ErrInvalidConfig, err, DocsPathConfig, DocsURL(DocsPathConfig))
}
