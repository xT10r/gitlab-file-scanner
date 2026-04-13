package filesystem_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gitlabFileScanner/internal/domain"
	"gitlabFileScanner/internal/infrastructure/filesystem"
)

func TestWriter_Save_Success(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	data := &domain.FileList{
		Name:      "test-project",
		WebURL:    "https://gitlab.com/test/project",
		ID:        12345,
		Branch:    "main",
		FilePaths: []string{"main.go", "utils/helper.go", "README.md"},
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was created
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["name"] != data.Name {
		t.Errorf("expected name %q, got %q", data.Name, parsed["name"])
	}

	files, ok := parsed["files"].([]any)
	if !ok || len(files) != len(data.FilePaths) {
		t.Errorf("expected %d files, got %v", len(data.FilePaths), parsed["files"])
	}
}

func TestWriter_Save_CreatesNestedDirectory(t *testing.T) {
	t.Helper()

	base := t.TempDir()
	target := filepath.Join(base, "nested", "deep", "output")

	w := filesystem.NewWriter()
	data := &domain.FileList{
		Name:      "project",
		WebURL:    "https://gitlab.com/p",
		ID:        1,
		Branch:    "main",
		FilePaths: []string{"a.go"},
	}

	path, err := w.Save(target, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected file to be created in nested directory")
	}
}

func TestWriter_Save_NilFilePaths(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	data := &domain.FileList{
		Name:      "empty-project",
		WebURL:    "https://gitlab.com/empty",
		ID:        99,
		Branch:    "main",
		FilePaths: nil,
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// nil slice should serialize as null in JSON
	if parsed["files"] != nil {
		t.Errorf("expected null files for nil slice, got %v", parsed["files"])
	}
}

func TestWriter_Save_EmptyFilePaths(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	data := &domain.FileList{
		Name:      "empty-project",
		WebURL:    "https://gitlab.com/empty",
		ID:        99,
		Branch:    "main",
		FilePaths: []string{},
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	files, ok := parsed["files"].([]any)
	if !ok {
		t.Fatalf("expected files to be an array")
	}
	if len(files) != 0 {
		t.Errorf("expected empty array, got %v", parsed["files"])
	}
}

func TestWriter_Save_LargeFileList(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	paths := make([]string, 10000)
	for i := range paths {
		paths[i] = filepath.Join("deep", "nested", "path", "file%d.go")
	}

	data := &domain.FileList{
		Name:      "large-project",
		WebURL:    "https://gitlab.com/large",
		ID:        777,
		Branch:    "develop",
		FilePaths: paths,
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	if !json.Valid(content) {
		t.Fatal("produced invalid JSON for large dataset")
	}
}

func TestWriter_Save_UnicodePaths(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	data := &domain.FileList{
		Name:      "unicode-project",
		WebURL:    "https://gitlab.com/unicode",
		ID:        66,
		Branch:    "main",
		FilePaths: []string{"файл.go", "目录/test.py", "αβγ.js"},
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	files, ok := parsed["files"].([]any)
	if !ok {
		t.Fatalf("expected files to be an array")
	}
	if len(files) != 3 {
		t.Errorf("expected 3 paths, got %d", len(files))
	}
}

func TestWriter_Save_FileNameFormat(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	w := filesystem.NewWriter()

	data := &domain.FileList{
		Name:      "project",
		WebURL:    "https://gitlab.com/p",
		ID:        12345,
		Branch:    "main",
		FilePaths: []string{"a.go"},
	}

	path, err := w.Save(dir, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(dir, "12345.json")
	if path != expected {
		t.Errorf("expected path %q, got %q", expected, path)
	}
}
