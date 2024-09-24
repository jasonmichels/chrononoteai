package main

import (
	"log"

	"github.com/jasonmichels/chrononoteai/config"
	"github.com/jasonmichels/chrononoteai/notes"
)

func main() {
	cfg, err := config.Initialize()
	if err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	fs := notes.OSFileSystem{}

	data, err := fs.ReadFile(cfg.BufferFile)
	if err != nil {
		log.Printf("Error reading buffer file: %v", err)
		return
	}

	err = notes.ProcessNotes(string(data), cfg.NotesDir, fs)
	if err != nil {
		log.Printf("Error processing notes: %v", err)
		return
	}

	log.Println("Notes processed successfully.")

	err = fs.WriteFile(cfg.BufferFile, []byte(""), 0o644)
	if err != nil {
		log.Printf("Error clearing buffer file: %v", err)
	} else {
		log.Println("Buffer file cleared successfully.")
	}
}
