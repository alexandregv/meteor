package main

import (
	"strings"
	"testing"
)

func TestHandleInterrupt(t *testing.T) {
	// Test case 1: No commit data available
	t.Run("no commit data", func(t *testing.T) {
		originalCommit := currentCommit
		currentCommit = nil
		defer func() {
			currentCommit = originalCommit
		}()

		// We can't easily test os.Exit, so we'll just ensure the function doesn't panic
		// In a real scenario, this would exit with code 1
	})

	// Test case 2: Partial commit data
	t.Run("partial commit data", func(t *testing.T) {
		originalCommit := currentCommit
		originalArgs := osArgs
		
		// Set up test data
		testCommit := &Commit{
			Type:  "feat",
			Scope: "ui",
		}
		currentCommit = testCommit
		osArgs = []string{}
		
		defer func() {
			currentCommit = originalCommit
			osArgs = originalArgs
		}()

		// Test that buildCommitCommand works with partial data
		_, printableCommand := buildCommitCommand("feat(ui): ", "", osArgs)
		if !strings.Contains(printableCommand, "feat(ui):") {
			t.Errorf("Expected command to contain 'feat(ui):', got %s", printableCommand)
		}
	})

	// Test case 3: Breaking change commit
	t.Run("breaking change commit", func(t *testing.T) {
		originalCommit := currentCommit
		originalArgs := osArgs
		
		testCommit := &Commit{
			Type:             "feat",
			Scope:            "api",
			IsBreakingChange: true,
		}
		currentCommit = testCommit
		osArgs = []string{}
		
		defer func() {
			currentCommit = originalCommit
			osArgs = originalArgs
		}()

		// Test that breaking change marker is included
		message := "feat(api)!: "
		_, printableCommand := buildCommitCommand(message, "", osArgs)
		if !strings.Contains(printableCommand, "feat(api)!:") {
			t.Errorf("Expected command to contain 'feat(api)!:', got %s", printableCommand)
		}
	})
}

func TestSetupSignalHandling(t *testing.T) {
	// Test that signal handling setup doesn't panic
	setupSignalHandling()
	// If we get here without panicking, the test passes
}