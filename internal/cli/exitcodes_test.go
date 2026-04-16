package cli

import (
	"errors"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
)

func TestExitCodes_AreStable(t *testing.T) {
	t.Parallel()

	if ExitOK != 0 {
		t.Fatalf("ExitOK changed: %d", ExitOK)
	}

	if ExitInvalidConfig != 2 {
		t.Fatalf("ExitInvalidConfig changed: %d", ExitInvalidConfig)
	}

	if ExitMissingDependency != 3 {
		t.Fatalf("ExitMissingDependency changed: %d", ExitMissingDependency)
	}

	if ExitRuntimeError != 4 {
		t.Fatalf("ExitRuntimeError changed: %d", ExitRuntimeError)
	}

	if ExitDoctorIssues != 5 {
		t.Fatalf("ExitDoctorIssues changed: %d", ExitDoctorIssues)
	}
}

func TestExitCodes_ErrorMappingIsDeterministic(t *testing.T) {
	t.Parallel()

	if got := exitCodeForError(nil); got != ExitOK {
		t.Fatalf("nil error should map to ExitOK, got %d", got)
	}

	if got := exitCodeForError(app.WrapInvalidConfig(errors.New("bad config"))); got != ExitInvalidConfig {
		t.Fatalf("invalid config should map to ExitInvalidConfig, got %d", got)
	}

	if got := exitCodeForError(app.ErrMissingProviderBinaries); got != ExitMissingDependency {
		t.Fatalf("missing dependency should map to ExitMissingDependency, got %d", got)
	}

	if got := exitCodeForError(app.ErrUnsafeMutationPath); got != ExitRuntimeError {
		t.Fatalf("unsafe path guard should map to ExitRuntimeError, got %d", got)
	}

	if got := exitCodeForError(errors.New("boom")); got != ExitRuntimeError {
		t.Fatalf("runtime errors should map to ExitRuntimeError, got %d", got)
	}
}
