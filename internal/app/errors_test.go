package app

import (
	"errors"
	"strings"
	"testing"
)

func TestWrapInvalidConfig_IncludesDocsReferences(t *testing.T) {
	t.Parallel()

	err := WrapInvalidConfig(errors.New("broken yaml"))
	msg := err.Error()

	if !strings.Contains(msg, DocsPathConfig) {
		t.Fatalf("expected docs path in error, got %q", msg)
	}

	if !strings.Contains(msg, DocsURL(DocsPathConfig)) {
		t.Fatalf("expected docs URL in error, got %q", msg)
	}
}

func TestUnsafeMutationPathError_IsClassifiable(t *testing.T) {
	t.Parallel()

	err := ErrUnsafeMutationPath
	if !errors.Is(err, ErrUnsafeMutationPath) {
		t.Fatalf("expected unsafe mutation path classification")
	}
}
