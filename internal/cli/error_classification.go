package cli

import (
	"errors"

	"github.com/Jfgm299/weave-cli/internal/app"
)

func isInvalidConfigError(err error) bool {
	return errors.Is(err, app.ErrInvalidConfig)
}

func isMissingDependencyError(err error) bool {
	return errors.Is(err, app.ErrMissingProviderBinaries)
}

func isScopeGuardError(err error) bool {
	return errors.Is(err, app.ErrUnsafeMutationPath)
}
