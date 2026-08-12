package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintUsesRequestedValues(t *testing.T) {
	var out, errOut bytes.Buffer
	code := Run([]string{"--print", "-y", "2026", "-u", "Tomoya Ogawa"}, strings.NewReader(""), &out, &errOut, "v1.0")
	if code != 0 || errOut.Len() != 0 {
		t.Fatalf("Run() = %d, stderr = %q", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Copyright (c) 2026 Tomoya Ogawa") {
		t.Fatal("requested values were not printed")
	}
}

func TestDefaultFileAndOverwriteConfirmation(t *testing.T) {
	dir := t.TempDir()
	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MITC_CONFIG_HOME", filepath.Join(dir, "home"))
	if err := os.Mkdir(filepath.Join(dir, "home"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"-y", "2026"}, strings.NewReader(""), &out, &errOut, "v1.0"); code != 0 {
		t.Fatalf("initial run = %d: %s", code, errOut.String())
	}
	license, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(license), "Copyright (c) 2026 John Doe") {
		t.Fatal("default user was not used")
	}
	out.Reset()
	if code := Run([]string{"-y", "2026"}, strings.NewReader("n\n"), &out, &errOut, "v1.0"); code != 0 || !strings.Contains(out.String(), "Cancelled.") {
		t.Fatalf("declined overwrite = %d, %q", code, out.String())
	}
}

func TestSetUserIsUsedLater(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MITC_CONFIG_HOME", home)
	var out, errOut bytes.Buffer
	if code := Run([]string{"--set-user", "Saved User"}, strings.NewReader(""), &out, &errOut, "v1.0"); code != 0 {
		t.Fatalf("set user = %d: %s", code, errOut.String())
	}
	out.Reset()
	if code := Run([]string{"-p", "-y", "2026"}, strings.NewReader(""), &out, &errOut, "v1.0"); code != 0 {
		t.Fatalf("print = %d: %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "Copyright (c) 2026 Saved User") {
		t.Fatal("saved user was not loaded")
	}
}

func TestRejectsInvalidCombinations(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"-p", "-f", "LICENSE"}, strings.NewReader(""), &out, &errOut, "v1.0"); code != 2 {
		t.Fatalf("Run() = %d, want 2", code)
	}
}
