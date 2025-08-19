package main

import (
	"strings"
	"testing"
)

func TestIntegrationCtrlCBehavior(t *testing.T) {
	tests := []struct {
		name     string
		commit   Commit
		args     []string
		expected []string // substrings that should be in the generated command
	}{
		{
			name: "basic type and scope",
			commit: Commit{
				Type:  "feat",
				Scope: "ui",
			},
			args:     []string{},
			expected: []string{"git", "commit", "-m", "feat(ui):"},
		},
		{
			name: "breaking change",
			commit: Commit{
				Type:             "feat",
				Scope:            "api", 
				IsBreakingChange: true,
			},
			args:     []string{},
			expected: []string{"git", "commit", "-m", "feat(api)!:"},
		},
		{
			name: "with complete message and body",
			commit: Commit{
				Type:    "fix",
				Scope:   "auth",
				Message: "fix(auth): resolve login issue",
				Body:    "Fixed authentication flow",
			},
			args:     []string{"--no-verify"},
			expected: []string{"git", "commit", "-m", "fix(auth): resolve login issue", "-m", "Fixed authentication flow", "--no-verify"},
		},
		{
			name: "no scope",
			commit: Commit{
				Type: "docs",
			},
			args:     []string{},
			expected: []string{"git", "commit", "-m", "docs:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up global state as handleInterrupt would see it
			originalCommit := currentCommit
			originalArgs := osArgs
			
			currentCommit = &tt.commit
			osArgs = tt.args
			
			defer func() {
				currentCommit = originalCommit
				osArgs = originalArgs
			}()

			// Test the commit message generation logic
			message := tt.commit.Message
			body := tt.commit.Body
			
			if message == "" && tt.commit.Type != "" {
				if tt.commit.IsBreakingChange {
					if tt.commit.Scope != "" {
						message = tt.commit.Type + "(" + tt.commit.Scope + ")!: "
					} else {
						message = tt.commit.Type + "!: "
					}
				} else {
					if tt.commit.Scope != "" {
						message = tt.commit.Type + "(" + tt.commit.Scope + "): "
					} else {
						message = tt.commit.Type + ": "
					}
				}
			}

			_, printableCommand := buildCommitCommand(message, body, osArgs)

			// Check that all expected substrings are in the command
			for _, expected := range tt.expected {
				if !strings.Contains(printableCommand, expected) {
					t.Errorf("Expected command to contain '%s', got: %s", expected, printableCommand)
				}
			}
		})
	}
}

func TestBuildCommitCommandWithAdditionalArgs(t *testing.T) {
	// Test that additional command line arguments are preserved
	message := "feat: add new feature"
	body := "Detailed description"
	args := []string{"--no-verify", "--author=test@example.com"}
	
	_, printableCommand := buildCommitCommand(message, body, args)
	
	expectedParts := []string{
		"git",
		"commit",
		"-m",
		"feat: add new feature",
		"-m", 
		"Detailed description",
		"--no-verify",
		"--author=test@example.com",
	}
	
	for _, part := range expectedParts {
		if !strings.Contains(printableCommand, part) {
			t.Errorf("Expected command to contain '%s', got: %s", part, printableCommand)
		}
	}
}