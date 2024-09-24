package notes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type MockFileSystem struct {
	Files map[string]string
	Dirs  map[string]bool
}

func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files: make(map[string]string),
		Dirs:  make(map[string]bool),
	}
}

func (fs *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if data, exists := fs.Files[path]; exists {
		return []byte(data), nil
	}
	return nil, os.ErrNotExist
}

func (fs *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	fs.Files[path] = string(data)
	return nil
}

func (fs *MockFileSystem) AppendToFile(path string, data string) error {
	fs.Files[path] += data
	return nil
}

func (fs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	fs.Dirs[path] = true
	return nil
}

func TestFormatNoteContent_PostProcessing(t *testing.T) {
	note := Note{
		Title:   "Test Note",
		Date:    "2023-10-01",
		Tags:    []string{"testing", "golang"},
		Content: "This is a test note content.",
	}

	fullNote, err := formatNoteContent(note)
	if err != nil {
		t.Fatalf("formatNoteContent failed: %v", err)
	}

	expectedContent := `---
title: Test Note
date: 2023-10-01
tags:
    - testing
    - golang
---
This is a test note content.

`

	if fullNote != expectedContent {
		t.Errorf("Full note content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, fullNote)
	}
}

func TestProcessNotes_ValidNotes(t *testing.T) {
	data := `---
title: Test Note
date: 2023-10-01
tags:
    - testing
    - golang
---
This is a test note content.
`

	fs := NewMockFileSystem()
	err := ProcessNotes(data, "/notes", fs)
	if err != nil {
		t.Fatalf("ProcessNotes failed: %v", err)
	}

	expectedPath := filepath.Join("/notes", "2023/10", "01.md")
	if _, exists := fs.Files[expectedPath]; !exists {
		t.Errorf("Expected file %s to be created", expectedPath)
	}

	expectedContent := `---
title: Test Note
date: 2023-10-01
tags:
    - testing
    - golang
---
This is a test note content.

`

	if fs.Files[expectedPath] != expectedContent {
		t.Errorf("File content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, fs.Files[expectedPath])
	}
}

func TestProcessNotes_InvalidNotes(t *testing.T) {
	data := `---
title: 
date: 2023-10-01
tags:
  - testing
---
Content without a title.
`

	fs := NewMockFileSystem()
	err := ProcessNotes(data, "/notes", fs)
	if err == nil {
		t.Fatal("Expected error due to missing title, but got none")
	}

	if !strings.Contains(err.Error(), "missing title") {
		t.Errorf("Expected error about missing title, got: %v", err)
	}
}

func TestValidateNote(t *testing.T) {
	validNote := Note{
		Title: "Valid Note",
		Date:  "2023-10-01",
	}

	if err := validateNote(validNote); err != nil {
		t.Errorf("Expected valid note, got error: %v", err)
	}

	invalidNote := Note{
		Title: "",
		Date:  "2023-10-01",
	}

	if err := validateNote(invalidNote); err == nil {
		t.Error("Expected error due to missing title, got none")
	}
}

func TestParseNotes(t *testing.T) {
	data := `---
title: First Note
date: 2023-10-01
tags:
  - test
---
Content of the first note.
---
title: Second Note
date: 2023-10-02
tags:
  - test
---
Content of the second note.
`

	notes, err := parseNotes(data)
	if err != nil {
		t.Fatalf("parseNotes failed: %v", err)
	}

	if len(notes) != 2 {
		t.Fatalf("Expected 2 notes, got %d", len(notes))
	}

	if notes[0].Title != "First Note" {
		t.Errorf("Expected first note title 'First Note', got '%s'", notes[0].Title)
	}

	if notes[1].Title != "Second Note" {
		t.Errorf("Expected second note title 'Second Note', got '%s'", notes[1].Title)
	}
}
