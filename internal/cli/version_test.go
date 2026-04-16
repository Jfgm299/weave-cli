package cli

import "testing"

func TestVersion_IsSemver(t *testing.T) {
	t.Parallel()

	if Version() != "0.1.0" {
		t.Fatalf("expected baseline semver 0.1.0, got %q", Version())
	}
}
