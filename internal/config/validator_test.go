package config

import (
	"errors"
	"testing"
)

func TestValidator_Validate_ValidV1Config(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 1, Sync: Sync{Mode: "symlink"}})
	if err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}

func TestValidator_Validate_InvalidVersionFails(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 0, Sync: Sync{Mode: "symlink"}})
	if err == nil {
		t.Fatalf("expected error for invalid version")
	}
	if !errors.Is(err, ErrOutdatedSchema) {
		t.Fatalf("expected outdated schema error, got %v", err)
	}
}

func TestValidator_Validate_NonSymlinkModeFailsInV1(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 1, Sync: Sync{Mode: "copy"}})
	if err == nil {
		t.Fatalf("expected error for non-symlink mode")
	}
}

func TestValidator_Validate_ConflictPolicyAcceptsSupportedValues(t *testing.T) {
	t.Parallel()

	policies := []string{"prompt", "overwrite", "skip", "backup"}
	for _, policy := range policies {
		policy := policy
		t.Run(policy, func(t *testing.T) {
			t.Parallel()

			err := (Validator{}).Validate(Config{Version: 1, Sync: Sync{Mode: "symlink", ConflictPolicy: policy}})
			if err != nil {
				t.Fatalf("expected policy %s to be valid, got %v", policy, err)
			}
		})
	}
}

func TestValidator_Validate_ConflictPolicyRejectsUnknownValue(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 1, Sync: Sync{Mode: "symlink", ConflictPolicy: "ask"}})
	if err == nil {
		t.Fatalf("expected unsupported conflict_policy error")
	}
}

func TestValidator_Validate_UnsupportedFutureSchemaFails(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 99, Sync: Sync{Mode: "symlink"}})
	if err == nil {
		t.Fatalf("expected unsupported schema error")
	}
	if !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatalf("expected unsupported schema classification, got %v", err)
	}
}
