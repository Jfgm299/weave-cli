package config

import "testing"

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
}

func TestValidator_Validate_NonSymlinkModeFailsInV1(t *testing.T) {
	t.Parallel()

	err := (Validator{}).Validate(Config{Version: 1, Sync: Sync{Mode: "copy"}})
	if err == nil {
		t.Fatalf("expected error for non-symlink mode")
	}
}
