package cli

import "testing"

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
}
