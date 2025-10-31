package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPassFiles tests all files that should parse successfully
func TestPassFiles(t *testing.T) {
	testDir := "__test__"

	// First, build the binary
	buildCmd := exec.Command("go", "build", "-o", "jsonparser_test_bin", "main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("jsonparser_test_bin")

	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "pass") && strings.HasSuffix(file.Name(), ".json") {
			t.Run(file.Name(), func(t *testing.T) {
				filePath := filepath.Join(testDir, file.Name())

				cmd := exec.Command("./jsonparser_test_bin", filePath)
				output, err := cmd.CombinedOutput()

				// Exit code 0 means success
				if err != nil {
					t.Errorf("File %s should parse successfully but failed with error: %v\nOutput: %s",
						file.Name(), err, string(output))
				}
			})
		}
	}
}

// TestFailFiles tests all files that should fail to parse
func TestFailFiles(t *testing.T) {
	testDir := "__test__"

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "jsonparser_test_bin", "main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("jsonparser_test_bin")

	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "fail") && strings.HasSuffix(file.Name(), ".json") {
			t.Run(file.Name(), func(t *testing.T) {
				filePath := filepath.Join(testDir, file.Name())

				cmd := exec.Command("./jsonparser_test_bin", filePath)
				output, err := cmd.CombinedOutput()

				// Non-zero exit code means failure (expected for these files)
				if err == nil {
					t.Errorf("File %s should fail to parse but succeeded\nOutput: %s",
						file.Name(), string(output))
				} else {
					t.Logf("File %s correctly failed: %v", file.Name(), err)
				}
			})
		}
	}
}
