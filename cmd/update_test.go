package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestUpdateCommand_DevVersionSkipsUpdate(t *testing.T) {
	oldVersion := Version
	t.Cleanup(func() { Version = oldVersion })

	Version = "dev"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"update"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Running development version. Skipping update.") {
		t.Errorf("Expected skip message for dev version, got %q", output)
	}
}

func TestUpdateCommand_PreviewFlagIsHidden(t *testing.T) {
	flag := updateCmd.Flags().Lookup("preview")
	if flag == nil {
		t.Fatal("Expected --preview flag to exist")
	}
	if !flag.Hidden {
		t.Error("Expected --preview flag to be hidden")
	}
}

func TestUpdateCommand_DevVersionSkipsEvenWithPreview(t *testing.T) {
	oldVersion := Version
	t.Cleanup(func() { Version = oldVersion })

	Version = "dev"

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"update", "--preview"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Running development version. Skipping update.") {
		t.Errorf("Expected skip message for dev version with --preview, got %q", output)
	}
}
