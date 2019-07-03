package fs

import (
	"os"
	"testing"
)

type MockFile struct {
	path    string
	content string
}

func NewMockFile(path, content string) *MockFile {
	return &MockFile{path, content}
}

func (file *MockFile) Pathname() string {
	return file.path
}

func (file *MockFile) MustCreate(t *testing.T) {
	f, err := os.Create(file.path)
	if err != nil {
		t.Fatalf("error while create file: %s: %v", file.path, err)
	}
	_, err = f.WriteString(file.content)
	if err != nil {
		t.Fatalf("error while write file: %s: %v", file.path, err)
	}
}

func (file *MockFile) MustRemove(t *testing.T) {
	if !Exists(file.path) {
		return
	}

	err := os.Remove(file.path)
	if err != nil {
		t.Fatalf("error while create file: %s: %v", file.path, err)
	}
}

type MockDirectory string

func (dir MockDirectory) Pathname() string {
	return string(dir)
}

func (dir MockDirectory) MustCreate(t *testing.T) {
	path := string(dir)
	if !Exists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			t.Fatalf("error while create directory: %s", path)
		}
	}
}

func (dir MockDirectory) MustRemove(t *testing.T) {
	path := string(dir)
	if Exists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("error while create directory: %s", path)
		}
	}
}
