package cmd

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	oldVersion := Version
	t.Cleanup(func() { Version = oldVersion })

	Version = "v0.1.0-test"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "hello-cli version v0.1.0-test " + runtime.GOOS + "/" + runtime.GOARCH + "\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestVersionCommand_Verbose(t *testing.T) {
	oldVersion, oldCommit, oldBuildTime := Version, Commit, BuildTime
	t.Cleanup(func() {
		Version, Commit, BuildTime = oldVersion, oldCommit, oldBuildTime
	})

	Version = "v0.1.0-test"
	Commit = "abc1234"
	BuildTime = "2026-03-31T12:00:00Z"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version", "--verbose"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()
	checks := []string{
		"hello-cli version: v0.1.0-test",
		"commit: abc1234",
		"build time: 2026-03-31T12:00:00Z",
		"platform: " + runtime.GOOS + "/" + runtime.GOARCH,
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected output to contain %q, but got %q", check, output)
		}
	}
}
