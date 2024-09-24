package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

// string constant for chrononoteai
const dirName = "chrononoteai"

type Config struct {
	BufferFile string `json:"buffer_file"`
	NotesDir   string `json:"notes_dir"`
	ConfigFile string // Path to the config file (not saved in JSON)
}

// InitializeWithArgs Modify Initialize to accept a FlagSet and arguments
func InitializeWithArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet(dirName, flag.ContinueOnError)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to get user home directory")
		return nil, err
	}
	defaultConfigPath := filepath.Join(homeDir, ".config", "chrononoteai", "config.json")

	configPath := fs.String("config", defaultConfigPath, "Path to the configuration file")
	bufferFile := fs.String("buffer", "", "Path to the buffer file")
	notesDir := fs.String("notes", "", "Path to the notes directory")

	if err := fs.Parse(args); err != nil {
		log.Println("Failed to parse command-line arguments")
		return nil, err
	}

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Println("Failed to load config")
		return nil, err
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
			log.Println("Failed to save config")
			return nil, err
		}
	}

	err = cfg.CreateBufferFileIfNeeded()
	if err != nil {
		return nil, err
	}

	logConfiguration(cfg)

	return cfg, nil
}

// Initialize function calls InitializeWithArgs with os.Args[1:]
func Initialize() (*Config, error) {
	return InitializeWithArgs(os.Args[1:])
}

func logConfiguration(cfg *Config) {
	log.Println("Configuration:")
	log.Printf("  Config File: %s\n", cfg.ConfigFile)
	log.Printf("  Buffer File: %s\n", cfg.BufferFile)
	log.Printf("  Notes Dir:   %s\n", cfg.NotesDir)
	log.Println("You can modify these settings in the config file or via command-line flags.")
}

// LoadConfig loads the configuration from the given path or initializes it with defaults.
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		ConfigFile: configPath,
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Use default values if config file doesn't exist
		defaultErr := config.setDefaults()
		if defaultErr != nil {
			return nil, defaultErr
		}

		mkDirErr := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
		if mkDirErr != nil {
			return nil, mkDirErr
		}

		if err := config.Save(); err != nil {
			return nil, err
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("Failed to read config file")
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		log.Println("Failed to parse config file")
		return nil, err
	}

	return config, nil
}

// CreateBufferFileIfNeeded checks if buffer file exists and if not it creates it
func (c *Config) CreateBufferFileIfNeeded() error {
	if _, err := os.Stat(c.BufferFile); os.IsNotExist(err) {
		bufferFile, err := os.Create(c.BufferFile)
		if err != nil {
			log.Println("Failed to create buffer file")
			return err
		}

		err = bufferFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Save writes the configuration to the config file.
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println("Failed to serialize config")
		return err
	}

	if err := os.WriteFile(c.ConfigFile, data, 0o644); err != nil {
		log.Println("Failed to write config file")
		return err
	}

	return nil
}

func (c *Config) setDefaults() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	c.BufferFile = filepath.Join(homeDir, ".config", dirName, "note.md")
	c.NotesDir = filepath.Join(homeDir, ".config", dirName, "notes")
	return nil
}
