package filefilter_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/filefilter"
)

func TestFilter_Apply_Success(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.go", "README.md", "utils.go", "test.py"},
		"*.go",
	)

	if len(result) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(result), result)
	}

	expected := map[string]bool{"main.go": true, "utils.go": true}
	for _, f := range result {
		if !expected[f] {
			t.Errorf("unexpected file: %s", f)
		}
	}
}

func TestFilter_Apply_MultipleMasks(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.go", "README.md", "test.py"},
		"*.go|*.md",
	)

	if len(result) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_GroupMask(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.go", "app.py", "index.js"},
		"*.(go|py)",
	)

	if len(result) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_NestedPaths(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"src/main.go", "src/utils/helper.go", "test/main_test.go"},
		"*.go",
	)

	if len(result) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_NoMatches(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.py", "test.js"},
		"*.go",
	)

	if len(result) != 0 {
		t.Errorf("expected 0 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_NilInput(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(nil, "*.go")

	if result != nil {
		t.Errorf("expected nil result for nil input, got %v", result)
	}
}

func TestFilter_Apply_EmptyInput(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply([]string{}, "*.go")

	if result != nil {
		t.Errorf("expected nil result for empty input, got %v", result)
	}
}

func TestFilter_Apply_EmptyMask(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply([]string{"main.go"}, "")

	// Empty mask after trimming = empty regex = matches only empty string
	if len(result) != 0 {
		t.Errorf("expected 0 files for empty mask, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_WildcardMask(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.go", "README.md", "test.py"},
		"*",
	)

	if len(result) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_CaseSensitive(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"main.GO", "Main.Go", "MAIN.go"},
		"*.go",
	)

	// Only MAIN.go matches (case sensitive)
	if len(result) != 1 {
		t.Fatalf("expected 1 file, got %d: %v", len(result), result)
	}
	if result[0] != "MAIN.go" {
		t.Errorf("expected 'MAIN.go', got %s", result[0])
	}
}

func TestFilter_Apply_Sorted(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"z.go", "a.go", "m.go"},
		"*.go",
	)

	for i := 1; i < len(result); i++ {
		if result[i-1] > result[i] {
			t.Errorf("result not sorted: %v", result)
			break
		}
	}
}

func TestFilter_Apply_UnicodePaths(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"файл.go", "目录/test.go"},
		"*.go",
	)

	if len(result) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(result), result)
	}
}

func TestFilter_Apply_SpecialCharacters(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	result := f.Apply(
		[]string{"src/main (1).go", "test [copy].go", "file.go"},
		"*.go",
	)

	if len(result) != 3 {
		t.Errorf("expected 3 files, got %d: %v", len(result), result)
	}
}
