//go:build !windows

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestPathCommandIsWindowsOnly(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"path", "add"}, strings.NewReader(""), &out, &errOut, "v1.0"); code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}
	if !strings.Contains(errOut.String(), "only supported on Windows") {
		t.Fatalf("unexpected error: %q", errOut.String())
	}
}
