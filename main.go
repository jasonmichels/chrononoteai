package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Note struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Content string   `yaml:"-"`
}

func main() {
	// Define the markdown directory
	markdownDir := "./notes" // this could be passed as a flag or config

	// Path to the buffer file
	bufferFile := filepath.Join(markdownDir, "chrononoteai.md")

	// Read the buffer file
	data, err := os.ReadFile(bufferFile)
	if err != nil {
		fmt.Println("Error reading buffer file:", err)
		return
	}

	// Split the notes by YAML front matter
	notes := parseNotes(string(data))

	// Process each note
	for _, note := range notes {
		log.Printf("Processing note for date: %s, title: %s", note.Date, note.Title)
		filePath := buildMarkdownPath(note, markdownDir)
		err := appendNoteToFile(filePath, note)
		if err != nil {
			log.Println("Error writing note to file:", err)
			return
		}
	}

	err = clearBufferFile(bufferFile)
	if err != nil {
		log.Printf("Error clearing buffer file: %v", err)
	} else {
		log.Println("Buffer file cleared successfully.")
	}
}

func parseNotes(data string) []Note {
	var notes []Note

	// Split the file by `---` to separate the YAML front matter sections
	entries := strings.Split(data, "---")
	for i := 1; i < len(entries); i += 2 { // Process each pair of YAML and content
		var note Note

		// Skip if the YAML part is empty
		if strings.TrimSpace(entries[i]) == "" {
			continue
		}

		// Unmarshal the YAML front matter
		if err := yaml.Unmarshal([]byte(entries[i]), &note); err != nil {
			log.Println("Error parsing YAML:", err)
			continue
		}

		// The content of the note is in the next segment
		if i+1 < len(entries) {
			note.Content = strings.TrimSpace(entries[i+1]) // Remove extra whitespace
		}

		notes = append(notes, note)
	}

	return notes
}

// Extract content after YAML front matter
func extractContent(entry string) string {
	split := strings.SplitN(entry, "---", 2)
	if len(split) > 1 {
		return strings.TrimSpace(split[1])
	}
	return ""
}

// Build the file path based on the note date
func buildMarkdownPath(note Note, baseDir string) string {
	// Parse the date
	noteDate, err := time.Parse("2006-01-02", note.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return ""
	}

	// Create the file path (e.g., /notes/2024/09/12.md)
	datePath := filepath.Join(baseDir, noteDate.Format("2006/01"))
	fileName := fmt.Sprintf("%02d.md", noteDate.Day())

	return filepath.Join(datePath, fileName)
}

// Append the note (with YAML front matter) to the corresponding markdown file
func appendNoteToFile(path string, note Note) error {
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	// Open the file (create if it doesn't exist)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Prepare YAML front matter
	yamlData, err := yaml.Marshal(note)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %v", err)
	}

	// Format the full note content with YAML front matter
	fullNote := fmt.Sprintf("---\n%s---\n%s\n\n", string(yamlData), note.Content)

	// Append the full note (YAML + content) to the file
	if _, err := f.WriteString(fullNote); err != nil {
		return err
	}

	return nil
}

// Clear the buffer file after processing
func clearBufferFile(path string) error {
	return os.WriteFile(path, []byte(""), 0o644)
}
