//go:build windows

package main

import "testing"

func TestSameWindowsPathIgnoresCaseAndTrailingSlash(t *testing.T) {
	if !sameWindowsPath(`C:\Tools\MITC\`, `c:\tools\mitc`) {
		t.Fatal("equivalent Windows paths were not matched")
	}
	if sameWindowsPath(`C:\Tools\MITC`, `C:\Tools\Other`) {
		t.Fatal("different Windows paths were matched")
	}
}
