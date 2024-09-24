package notes

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Note represents a single note with metadata and content.
type Note struct {
	Title   string   `yaml:"title"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Content string   `yaml:"-"`
}

// FrontMatter represents the YAML front matter of a note.
type FrontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
}

// FileSystem interface for dependency injection in file operations.
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	AppendToFile(path string, data string) error
	MkdirAll(path string, perm os.FileMode) error
}

// OSFileSystem implements FileSystem using the OS package.
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
		log.Printf("Failed to open file %s: %v", path, err)
		return err
	}
	defer func() {
		if closeError := f.Close(); closeError != nil && err == nil {
			err = fmt.Errorf("failed to close file %s: %w", path, closeError)
		}
	}()
	_, err = f.WriteString(data)
	if err != nil {
		log.Printf("Failed to write to file %s: %v", path, err)
		return err
	}
	return nil
}

func (fs OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// ProcessNotes parses, validates, and saves notes from the provided data.
func ProcessNotes(data, markdownDir string, fs FileSystem) error {
	notes, err := parseNotes(data)
	if err != nil {
		log.Println("Failed to parse notes")
		return err
	}

	// Validate all notes before processing
	for _, note := range notes {
		if err := validateNote(note); err != nil {
			log.Printf("Failed to validate note for date: %s, title: %s\n", note.Date, note.Title)
			return err
		}
	}

	// Process and save each note
	for _, note := range notes {
		log.Printf("Processing note for date: %s, title: %s\n", note.Date, note.Title)
		filePath, err := buildMarkdownPath(note, markdownDir)
		if err != nil {
			return err
		}

		if err := fs.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Printf("Failed to create directories for file %s: %v\n", filePath, err)
			return err
		}

		// Format the note with YAML front matter
		fullNote, err := formatNoteContent(note)
		if err != nil {
			return err
		}

		if err := fs.AppendToFile(filePath, fullNote); err != nil {
			log.Printf("Failed to write note to file %s: %v\n", filePath, err)
			return err
		}
		log.Printf("Wrote note to file %s\n", filePath)
	}

	return nil
}

// parseNotes splits the input data into individual notes.
func parseNotes(data string) ([]Note, error) {
	var notes []Note

	entries := strings.Split(data, "---")
	for i := 1; i < len(entries); i += 2 {
		var note Note

		metadata := entries[i]
		content := ""
		if i+1 < len(entries) {
			content = strings.TrimSpace(entries[i+1])
		}

		if strings.TrimSpace(metadata) == "" && content == "" {
			continue
		}

		if err := yaml.Unmarshal([]byte(metadata), &note); err != nil {
			log.Println("Failed to parse YAML")
			return nil, err
		}

		note.Content = content
		notes = append(notes, note)
	}

	return notes, nil
}

// validateNote checks if the note has all required fields and valid data.
func validateNote(note Note) error {
	if note.Title == "" {
		return errors.New("missing title")
	}
	if note.Date == "" {
		return errors.New("missing date")
	}
	if _, err := time.Parse("2006-01-02", note.Date); err != nil {
		log.Printf("Invalid date: %s\n", note.Date)
		return err
	}
	return nil
}

// buildMarkdownPath creates the file path for a note based on its date.
func buildMarkdownPath(note Note, baseDir string) (string, error) {
	noteDate, err := time.Parse("2006-01-02", note.Date)
	if err != nil {
		log.Printf("Invalid date: %s\n", note.Date)
		return "", err
	}

	datePath := filepath.Join(baseDir, noteDate.Format("2006/01"))
	fileName := fmt.Sprintf("%02d.md", noteDate.Day())

	return filepath.Join(datePath, fileName), nil
}

// formatNoteContent formats the note's content with YAML front matter.
func formatNoteContent(note Note) (string, error) {
	frontMatter := FrontMatter{
		Title: note.Title,
		Date:  note.Date,
		Tags:  note.Tags,
	}

	yamlFrontMatterBytes, err := yaml.Marshal(frontMatter)
	if err != nil {
		log.Println("Failed to marshal YAML front matter")
		return "", err
	}

	yamlFrontMatter := string(yamlFrontMatterBytes)

	// Post-process to remove quotes around the date field
	yamlFrontMatter = removeQuotesFromDateField(yamlFrontMatter, note.Date)

	return fmt.Sprintf("---\n%s---\n%s\n\n", yamlFrontMatter, note.Content), nil
}

// removeQuotesFromDateField removes quotes around the date field in the YAML front matter.
func removeQuotesFromDateField(yamlContent string, dateValue string) string {
	re := regexp.MustCompile(`(?m)^date:.*$`)
	unquotedDate := fmt.Sprintf("date: %s", dateValue)
	return re.ReplaceAllString(yamlContent, unquotedDate)
}
