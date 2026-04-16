package cli

import "testing"

func TestParseFromFlag_ParsesLongForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from", "/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_ParsesEqualsForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from=/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_MissingValueFails(t *testing.T) {
	t.Parallel()

	_, err := parseFromFlag([]string{"--from"})
	if err == nil {
		t.Fatalf("expected missing value error")
	}
}

func TestParseProviderAction_ListParsesWithoutName(t *testing.T) {
	t.Parallel()

	action, name, err := parseProviderAction([]string{"list"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if action != providerActionList {
		t.Fatalf("expected providerActionList, got %s", action)
	}

	if name != "" {
		t.Fatalf("expected empty provider name for list, got %q", name)
	}
}

func TestParseProviderAction_AddMissingNameFails(t *testing.T) {
	t.Parallel()

	_, _, err := parseProviderAction([]string{"add"})
	if err == nil {
		t.Fatalf("expected provider add parse error when name is missing")
	}
}
