package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Remove any existing config file to test defaults
	os.Remove(configPath)

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if defaults are set
	homeDir, _ := os.UserHomeDir()
	expectedBufferFile := filepath.Join(homeDir, ".config", "chrononoteai", "note.md")
	expectedNotesDir := filepath.Join(homeDir, ".config", "chrononoteai", "notes")

	if cfg.BufferFile != expectedBufferFile {
		t.Errorf("Expected BufferFile %s, got %s", expectedBufferFile, cfg.BufferFile)
	}
	if cfg.NotesDir != expectedNotesDir {
		t.Errorf("Expected NotesDir %s, got %s", expectedNotesDir, cfg.NotesDir)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write a sample config file
	sampleConfig := `{
        "buffer_file": "/tmp/buffer.md",
        "notes_dir": "/tmp/notes"
    }`
	err := os.WriteFile(configPath, []byte(sampleConfig), 0o644)
	if err != nil {
		t.Fatalf("Failed to write sample config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if values are loaded correctly
	if cfg.BufferFile != "/tmp/buffer.md" {
		t.Errorf("Expected BufferFile /tmp/buffer.md, got %s", cfg.BufferFile)
	}
	if cfg.NotesDir != "/tmp/notes" {
		t.Errorf("Expected NotesDir /tmp/notes, got %s", cfg.NotesDir)
	}
}

func TestInitializeWithArgs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Remove any existing config file to test defaults
	os.Remove(configPath)

	// Prepare command-line arguments
	args := []string{
		"--config", configPath,
		"--buffer", "/test/buffer.md",
		"--notes", "/test/notes",
	}

	cfg, err := InitializeWithArgs(args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if values are overridden correctly
	if cfg.BufferFile != "/test/buffer.md" {
		t.Errorf("Expected BufferFile /test/buffer.md, got %s", cfg.BufferFile)
	}
	if cfg.NotesDir != "/test/notes" {
		t.Errorf("Expected NotesDir /test/notes, got %s", cfg.NotesDir)
	}

	// Check if the config file was saved
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Expected config file to be saved at %s", configPath)
	}
}

func TestInitializeWithExistingConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write a sample config file
	sampleConfig := `{
        "buffer_file": "/existing/buffer.md",
        "notes_dir": "/existing/notes"
    }`
	err := os.WriteFile(configPath, []byte(sampleConfig), 0o644)
	if err != nil {
		t.Fatalf("Failed to write sample config file: %v", err)
	}

	// Prepare command-line arguments (no overrides)
	args := []string{
		"--config", configPath,
	}

	cfg, err := InitializeWithArgs(args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if values are loaded from the config file
	if cfg.BufferFile != "/existing/buffer.md" {
		t.Errorf("Expected BufferFile /existing/buffer.md, got %s", cfg.BufferFile)
	}
	if cfg.NotesDir != "/existing/notes" {
		t.Errorf("Expected NotesDir /existing/notes, got %s", cfg.NotesDir)
	}
}
