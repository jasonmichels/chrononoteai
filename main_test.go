package main

import (
	"os"
	"testing"
)

// MockFileSystem for unit testing
type MockFileSystem struct {
	ReadData     []byte
	WriteData    []byte
	AppendedData string
	Err          error
}

func (mfs *MockFileSystem) ReadFile(path string) ([]byte, error) {
	return mfs.ReadData, mfs.Err
}

func (mfs *MockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	mfs.WriteData = data
	return mfs.Err
}

func (mfs *MockFileSystem) AppendToFile(path string, data string) error {
	mfs.AppendedData = data
	return mfs.Err
}

func (mfs *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return mfs.Err
}

// Test for processNotes function
func TestProcessNotes(t *testing.T) {
	mockFS := &MockFileSystem{}
	data := `
---
title: Test Note
date: 2024-09-12
tags:
  - test
---
This is the content of the test note.
`
	err := processNotes(data, "./notes", mockFS)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedData := `---
title: Test Note
date: 2024-09-12
tags:
  - test
---
This is the content of the test note.

`
	if mockFS.AppendedData != expectedData {
		t.Fatalf("expected %s, got %s", expectedData, mockFS.AppendedData)
	}
}

func TestProcessNotesWithEmptyYAML(t *testing.T) {
	mockFS := &MockFileSystem{}
	data := `---`

	err := processNotes(data, "./notes", mockFS)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mockFS.AppendedData != "" {
		t.Fatalf("expected no appended data, got %s", mockFS.AppendedData)
	}
}

func TestBuildMarkdownPath(t *testing.T) {
	note := Note{
		Date: "2024-09-12",
	}
	expectedPath := "notes/2024/09/12.md"

	path := buildMarkdownPath(note, "./notes")
	if path != expectedPath {
		t.Fatalf("expected %s, got %s", expectedPath, path)
	}
}
