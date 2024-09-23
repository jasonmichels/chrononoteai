package main

import (
	"fmt"
	"github.com/jasonmichels/chrononoteai/config"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Note structure to hold metadata and content
type Note struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Content string   `yaml:"-"`
}

// FileSystem interface for dependency injection
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	AppendToFile(path string, data string) error
	MkdirAll(path string, perm os.FileMode) error
}

// OSFileSystem implementation of FileSystem that uses actual OS calls
type OSFileSystem struct{}

func (fs OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs OSFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (fs OSFileSystem) AppendToFile(path string, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	return err
}

func (fs OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Main function
func main() {
	// lets try to refactor this to make it more readable
	// first, read configuration and get file path, or some sort of config object
	// second, load the file into memory
	// third, validate the file to make sure we won't have any issues with the notes
	// fourth, process each individual note, which right now is just saving to correct location
	// fifth, empty the buffer file and return a success
	cfg, err := config.Initialize()
	if err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	fs := OSFileSystem{}

	data, err := fs.ReadFile(cfg.BufferFile)
	if err != nil {
		log.Printf("Error reading buffer file: %v", err)
		return
	}

	err = processNotes(string(data), cfg.NotesDir, fs)
	if err != nil {
		log.Printf("Error processing notes: %v", err)
		return
	}

	err = fs.WriteFile(cfg.BufferFile, []byte(""), 0o644)
	if err != nil {
		log.Printf("Error clearing buffer file: %v", err)
	} else {
		log.Println("Buffer file cleared successfully.")
	}
}

// processNotes handles the logic of parsing notes and appending them to files
func processNotes(data, markdownDir string, fs FileSystem) error {
	notes := parseNotes(data)

	for _, note := range notes {
		log.Printf("Processing note for date: %s, title: %s", note.Date, note.Title)
		filePath := buildMarkdownPath(note, markdownDir)

		if err := fs.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		// Manually format YAML front matter to avoid quotes around the date and format tags properly
		tags := ""
		for _, tag := range note.Tags {
			tags += fmt.Sprintf("  - %s\n", tag)
		}

		yamlFrontMatter := fmt.Sprintf(`---
title: %s
date: %s
tags:
%s---`, note.Title, note.Date, tags)

		// Format the full note content with YAML front matter
		fullNote := fmt.Sprintf("%s\n%s\n\n", yamlFrontMatter, note.Content)

		if err := fs.AppendToFile(filePath, fullNote); err != nil {
			return err
		}
	}

	return nil
}

// parseNotes parses the notes from the buffer content
func parseNotes(data string) []Note {
	var notes []Note

	entries := strings.Split(data, "---")
	for i := 1; i < len(entries); i += 2 {
		var note Note

		if strings.TrimSpace(entries[i]) == "" {
			continue
		}

		if err := yaml.Unmarshal([]byte(entries[i]), &note); err != nil {
			log.Println("Error parsing YAML:", err)
			continue
		}

		if i+1 < len(entries) {
			note.Content = strings.TrimSpace(entries[i+1])
		}

		notes = append(notes, note)
	}

	return notes
}

// buildMarkdownPath builds the file path based on the note date
func buildMarkdownPath(note Note, baseDir string) string {
	noteDate, err := time.Parse("2006-01-02", note.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return ""
	}

	datePath := filepath.Join(baseDir, noteDate.Format("2006/01"))
	fileName := fmt.Sprintf("%02d.md", noteDate.Day())

	return filepath.Join(datePath, fileName)
}
