package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeWithArgs_Defaults(t *testing.T) {
	// Suppress log output during testing
	log.SetOutput(os.Stdout)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Prepare command-line arguments with default values
	args := []string{
		"--config", configPath,
	}

	cfg, err := InitializeWithArgs(args)
	if err != nil {
		t.Fatalf("InitializeWithArgs failed: %v", err)
	}

	// Verify that the config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Verify that the buffer file was created
	if _, err := os.Stat(cfg.BufferFile); os.IsNotExist(err) {
		t.Fatalf("Buffer file was not created at %s", cfg.BufferFile)
	}

	// Check default values
	homeDir, _ := os.UserHomeDir()
	expectedBufferFile := filepath.Join(homeDir, ".config", dirName, "note.md")
	expectedNotesDir := filepath.Join(homeDir, ".config", dirName, "notes")

	if cfg.BufferFile != expectedBufferFile {
		t.Errorf("Expected BufferFile %s, got %s", expectedBufferFile, cfg.BufferFile)
	}
	if cfg.NotesDir != expectedNotesDir {
		t.Errorf("Expected NotesDir %s, got %s", expectedNotesDir, cfg.NotesDir)
	}
}

func TestInitializeWithArgs_Overrides(t *testing.T) {
	// Suppress log output during testing
	log.SetOutput(os.Stdout)

	// Create temporary paths
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	bufferFilePath := filepath.Join(tempDir, "buffer.md")
	notesDirPath := filepath.Join(tempDir, "notes")

	// Prepare command-line arguments with overrides
	args := []string{
		"--config", configPath,
		"--buffer", bufferFilePath,
		"--notes", notesDirPath,
	}

	cfg, err := InitializeWithArgs(args)
	if err != nil {
		t.Fatalf("InitializeWithArgs failed: %v", err)
	}

	// Verify that the config file and buffer file were created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}
	if _, err := os.Stat(cfg.BufferFile); os.IsNotExist(err) {
		t.Fatalf("Buffer file was not created at %s", cfg.BufferFile)
	}

	// Check overridden values
	if cfg.BufferFile != bufferFilePath {
		t.Errorf("Expected BufferFile %s, got %s", bufferFilePath, cfg.BufferFile)
	}
	if cfg.NotesDir != notesDirPath {
		t.Errorf("Expected NotesDir %s, got %s", notesDirPath, cfg.NotesDir)
	}

	// Read the config file to ensure values were saved
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	configFromFile := &Config{}
	if err := json.Unmarshal(data, configFromFile); err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	if configFromFile.BufferFile != bufferFilePath {
		t.Errorf("Config file BufferFile mismatch: expected %s, got %s", bufferFilePath, configFromFile.BufferFile)
	}
	if configFromFile.NotesDir != notesDirPath {
		t.Errorf("Config file NotesDir mismatch: expected %s, got %s", notesDirPath, configFromFile.NotesDir)
	}
}

func TestLoadConfig_NewConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify that the config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}

	// Check default values
	homeDir, _ := os.UserHomeDir()
	expectedBufferFile := filepath.Join(homeDir, ".config", dirName, "note.md")
	expectedNotesDir := filepath.Join(homeDir, ".config", dirName, "notes")

	if cfg.BufferFile != expectedBufferFile {
		t.Errorf("Expected BufferFile %s, got %s", expectedBufferFile, cfg.BufferFile)
	}
	if cfg.NotesDir != expectedNotesDir {
		t.Errorf("Expected NotesDir %s, got %s", expectedNotesDir, cfg.NotesDir)
	}
}

func TestLoadConfig_ExistingConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Write a sample config file
	sampleConfig := `{
		"buffer_file": "/tmp/test_buffer.md",
		"notes_dir": "/tmp/test_notes"
	}`
	if err := os.WriteFile(configPath, []byte(sampleConfig), 0644); err != nil {
		t.Fatalf("Failed to write sample config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Verify loaded values
	if cfg.BufferFile != "/tmp/test_buffer.md" {
		t.Errorf("Expected BufferFile /tmp/test_buffer.md, got %s", cfg.BufferFile)
	}
	if cfg.NotesDir != "/tmp/test_notes" {
		t.Errorf("Expected NotesDir /tmp/test_notes, got %s", cfg.NotesDir)
	}
}

func TestCreateBufferFileIfNeeded(t *testing.T) {
	// Suppress log output during testing
	log.SetOutput(os.Stdout)

	// Create a temporary buffer file path
	tempDir := t.TempDir()
	bufferFilePath := filepath.Join(tempDir, "buffer.md")

	cfg := &Config{
		BufferFile: bufferFilePath,
	}

	// Ensure the buffer file does not exist
	if _, err := os.Stat(bufferFilePath); !os.IsNotExist(err) {
		os.Remove(bufferFilePath)
	}

	// Call CreateBufferFileIfNeeded
	if err := cfg.CreateBufferFileIfNeeded(); err != nil {
		t.Fatalf("CreateBufferFileIfNeeded failed: %v", err)
	}

	// Verify that the buffer file was created
	if _, err := os.Stat(bufferFilePath); os.IsNotExist(err) {
		t.Fatalf("Buffer file was not created at %s", bufferFilePath)
	}
}

func TestSave(t *testing.T) {
	// Suppress log output during testing
	log.SetOutput(os.Stdout)

	// Create a temporary config file path
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &Config{
		BufferFile: "/tmp/test_buffer.md",
		NotesDir:   "/tmp/test_notes",
		ConfigFile: configPath,
	}

	// Call Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify that the config file was written
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not saved at %s", configPath)
	}

	// Read and verify the config file contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	configFromFile := &Config{}
	if err := json.Unmarshal(data, configFromFile); err != nil {
		t.Fatalf("Failed to parse config file: %v", err)
	}

	if configFromFile.BufferFile != cfg.BufferFile {
		t.Errorf("Expected BufferFile %s, got %s", cfg.BufferFile, configFromFile.BufferFile)
	}
	if configFromFile.NotesDir != cfg.NotesDir {
		t.Errorf("Expected NotesDir %s, got %s", cfg.NotesDir, configFromFile.NotesDir)
	}
}

func TestSetDefaults(t *testing.T) {
	cfg := &Config{}

	err := cfg.setDefaults()
	if err != nil {
		t.Fatalf("setDefaults failed: %v", err)
	}

	// Check default values
	homeDir, _ := os.UserHomeDir()
	expectedBufferFile := filepath.Join(homeDir, ".config", dirName, "note.md")
	expectedNotesDir := filepath.Join(homeDir, ".config", dirName, "notes")

	if cfg.BufferFile != expectedBufferFile {
		t.Errorf("Expected BufferFile %s, got %s", expectedBufferFile, cfg.BufferFile)
	}
	if cfg.NotesDir != expectedNotesDir {
		t.Errorf("Expected NotesDir %s, got %s", expectedNotesDir, cfg.NotesDir)
	}
}
