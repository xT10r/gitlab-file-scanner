package filefilter_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/filefilter"
)

func TestFilter_Deduplication(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	// Маска с пересекающимися паттернами — один файл может совпасть дважды
	result := f.Apply(
		[]string{"main.go", "README.md"},
		"*.go|*.go", // дубликат маски
	)

	if len(result) != 1 {
		t.Fatalf("expected 1 file (deduplicated), got %d: %v", len(result), result)
	}
}

func TestFilter_Deduplication_OverlappingMasks(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	// *.go и *go — оба совпадут с main.go
	result := f.Apply(
		[]string{"main.go", "test.go", "README.md"},
		"*.go|*go",
	)

	if len(result) != 2 {
		t.Fatalf("expected 2 files (deduplicated), got %d: %v", len(result), result)
	}
}

func TestFilter_RegexSpecialChars(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	// [ ] должны экранироваться — не быть character class
	result := f.Apply(
		[]string{"test[0].go", "t.go", "e.go", "s.go"},
		"test[0].go",
	)

	// Должен совпасть только test[0].go, не t/e/s.go
	if len(result) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(result), result)
	}
	if result[0] != "test[0].go" {
		t.Errorf("expected 'test[0].go', got %s", result[0])
	}
}

func TestFilter_PlusQuestionMark(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	// + и ? должны экранироваться
	result := f.Apply(
		[]string{"file+name.go", "file?name.go", "filename.go"},
		"file+name.go",
	)

	if len(result) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(result), result)
	}
	if result[0] != "file+name.go" {
		t.Errorf("expected 'file+name.go', got %s", result[0])
	}
}
