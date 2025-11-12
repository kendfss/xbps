package main

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCommandDiscovery(t *testing.T) {
	tmpDir := t.TempDir()

	bin1 := filepath.Join(tmpDir, "bin1")
	bin2 := filepath.Join(tmpDir, "bin2")

	if err := os.MkdirAll(bin1, 0o755); err != nil {
		t.Fatalf("Failed to create bin1: %v", err)
	}
	if err := os.MkdirAll(bin2, 0o755); err != nil {
		t.Fatalf("Failed to create bin2: %v", err)
	}

	testCommands := []string{
		"xbps-install",
		"xbps-remove",
		"xbps-query",
		"xbps-reconfigure",
		"xbps-alternatives",
	}

	for _, cmd := range testCommands {
		cmdPath := filepath.Join(bin1, cmd)
		if err := os.WriteFile(cmdPath, []byte("#!/bin/sh\necho "+cmd), 0o755); err != nil {
			t.Fatalf("Failed to create %s: %v", cmd, err)
		}
	}

	extraCommands := []string{
		"xbps-install", // duplicate
		"xbps-search",  // new command
	}

	for _, cmd := range extraCommands {
		cmdPath := filepath.Join(bin2, cmd)
		if err := os.WriteFile(cmdPath, []byte("#!/bin/sh\necho "+cmd), 0o755); err != nil {
			t.Fatalf("Failed to create %s: %v", cmd, err)
		}
	}

	nonExecPath := filepath.Join(bin1, "xbps-config.txt")
	if err := os.WriteFile(nonExecPath, []byte("config file"), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	testPATH := bin1 + string(os.PathListSeparator) + bin2

	t.Run("discovers all commands", func(t *testing.T) {
		result, err := children("xbps", testPATH)
		if err != nil {
			t.Fatalf("children failed: %v", err)
		}

		expected := map[string]string{
			"install":      "xbps-install",
			"remove":       "xbps-remove",
			"query":        "xbps-query",
			"reconfigure":  "xbps-reconfigure",
			"alternatives": "xbps-alternatives",
			"search":       "xbps-search", // from bin2
		}

		if len(result) != len(expected) {
			t.Errorf("Expected %d commands, got %d", len(expected), len(result))
			t.Errorf("want: %s", slices.Sorted(maps.Keys(expected)))
			t.Errorf("have: %s", slices.Sorted(maps.Keys(result)))
		}

		for key, expectedFull := range expected {
			full, exists := result[key]
			if !exists {
				t.Errorf("Missing command %q", key)
				continue
			}
			if full != expectedFull {
				t.Errorf("For key %q, expected %q, got %q", key, expectedFull, full)
			}
		}
	})

	t.Run("deduplicates commands", func(t *testing.T) {
		result, err := children("xbps", testPATH)
		if err != nil {
			t.Fatalf("children failed: %v", err)
		}
		// xbps-install exists in both bin1 and bin2, should use bin1's version
		if result["install"] != "xbps-install" {
			t.Error("Deduplication failed or wrong precedence")
		}
	})

	t.Run("handles empty PATH", func(t *testing.T) {
		result, err := children("xbps", "")
		if err != nil {
			t.Fatalf("children failed with empty PATH: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result with empty PATH, got %v", result)
		}
	})

	t.Run("handles non-existent directories in PATH", func(t *testing.T) {
		pathWithBadDir := "/non/existent/dir" + string(os.PathListSeparator) + bin1
		result, err := children("xbps", pathWithBadDir)
		if err != nil {
			t.Fatalf("children failed with bad directory in PATH: %v", err)
		}
		// Should still find commands from the valid directory
		if len(result) == 0 {
			t.Error("Expected to find commands from valid directory despite bad directory")
		}
	})

	t.Run("ignores non-executable files", func(t *testing.T) {
		result, err := children("xbps", bin1)
		if err != nil {
			t.Fatalf("children failed: %v", err)
		}
		if _, exists := result["config.txt"]; exists {
			t.Error("Found non-executable .txt file in results")
		}
	})
}

func TestCommandDiscoveryEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")

	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	t.Run("handles commands with multiple hyphens", func(t *testing.T) {
		multiHyphenCmd := filepath.Join(binDir, "xbps-src-update")
		if err := os.WriteFile(multiHyphenCmd, []byte("#!/bin/sh"), 0o755); err != nil {
			t.Fatalf("Failed to create multi-hyphen command: %v", err)
		}

		result, err := children("xbps", binDir)
		if err != nil {
			t.Fatalf("children failed: %v", err)
		}

		if full, exists := result["src-update"]; !exists || full != "xbps-src-update" {
			t.Errorf("Unexpected result for multi-hyphen command: %v", result)
		}
	})

	t.Run("handles empty directories", func(t *testing.T) {
		emptyDir := t.TempDir()
		result, err := children("xbps", emptyDir)
		if err != nil {
			t.Fatalf("children failed with empty directory: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result from empty directory, got %v", result)
		}
	})
}
