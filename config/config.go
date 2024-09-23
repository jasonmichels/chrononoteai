// config/config.go
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	BufferFile string `json:"buffer_file"`
	NotesDir   string `json:"notes_dir"`
	ConfigFile string // Path to the config file (not saved in JSON)
}

// InitializeWithArgs Modify Initialize to accept a FlagSet and arguments
func InitializeWithArgs(args []string) (*Config, error) {
	// Create a new FlagSet
	fs := flag.NewFlagSet("chrononoteai", flag.ContinueOnError)

	// Determine default config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	defaultConfigPath := filepath.Join(homeDir, ".config", "chrononoteai", "config.json")

	// Command-line flags
	configPath := fs.String("config", defaultConfigPath, "Path to the configuration file")
	bufferFile := fs.String("buffer", "", "Path to the buffer file")
	notesDir := fs.String("notes", "", "Path to the notes directory")

	// Parse the provided arguments
	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	// Load configuration
	cfg, err := LoadConfig(*configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Override with command-line arguments
	updated := false
	if *bufferFile != "" {
		cfg.BufferFile = *bufferFile
		updated = true
	}
	if *notesDir != "" {
		cfg.NotesDir = *notesDir
		updated = true
	}

	// Save updated configuration
	if updated {
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("error saving config: %w", err)
		}
	}

	// Log configuration information
	logConfiguration(cfg)

	return cfg, nil
}

// Initialize function calls InitializeWithArgs with os.Args[1:]
func Initialize() (*Config, error) {
	return InitializeWithArgs(os.Args[1:])
}

func logConfiguration(cfg *Config) {
	fmt.Println("Configuration:")
	fmt.Printf("  Config File: %s\n", cfg.ConfigFile)
	fmt.Printf("  Buffer File: %s\n", cfg.BufferFile)
	fmt.Printf("  Notes Dir:   %s\n", cfg.NotesDir)
	fmt.Println("You can modify these settings in the config file or via command-line flags.")
}

// LoadConfig loads the configuration from the given path or initializes it with defaults.
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		ConfigFile: configPath,
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Use default values if config file doesn't exist
		config.setDefaults()
		// Create config directory if it doesn't exist
		os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		// Save the default config
		if err := config.Save(); err != nil {
			return nil, err
		}
		return config, nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// Save writes the configuration to the config file.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write to config file
	if err := os.WriteFile(c.ConfigFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) setDefaults() {
	homeDir, _ := os.UserHomeDir()
	c.BufferFile = filepath.Join(homeDir, ".config", "chrononoteai", "note.md")
	c.NotesDir = filepath.Join(homeDir, ".config", "chrononoteai", "notes")
}
