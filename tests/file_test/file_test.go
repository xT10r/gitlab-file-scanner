package file_test

import (
	"gitlabFileScanner/internal/file"
	"testing"
)

func TestFilterFilesByMask(t *testing.T) {
	testCases := []struct {
		Description   string
		FilePaths     []string
		Mask          string
		ExpectedCount int
		ExpectedFiles []string
	}{
		{
			Description:   "Filter *.go files",
			FilePaths:     []string{"main.go", "test.py", "utils.go", "README.md"},
			Mask:          "*.go",
			ExpectedCount: 2,
			ExpectedFiles: []string{"main.go", "utils.go"},
		},
		{
			Description:   "Filter multiple masks",
			FilePaths:     []string{"main.go", "test.py", "README.md"},
			Mask:          "*.go|*.md",
			ExpectedCount: 2,
			ExpectedFiles: []string{"main.go", "README.md"},
		},
		{
			Description:   "Empty file list",
			FilePaths:     []string{},
			Mask:          "*.go",
			ExpectedCount: 0,
			ExpectedFiles: []string{},
		},
		{
			Description:   "Wildcard mask matches all",
			FilePaths:     []string{"main.go", "test.py", "README.md"},
			Mask:          "*",
			ExpectedCount: 3,
			ExpectedFiles: []string{"main.go", "test.py", "README.md"},
		},
		{
			Description:   "Nested paths",
			FilePaths:     []string{"src/main.go", "src/utils/helper.go", "test/main_test.go"},
			Mask:          "*.go",
			ExpectedCount: 3,
			ExpectedFiles: []string{"src/main.go", "src/utils/helper.go", "test/main_test.go"},
		},
		{
			Description:   "No matches",
			FilePaths:     []string{"main.py", "test.js"},
			Mask:          "*.go",
			ExpectedCount: 0,
			ExpectedFiles: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result := file.FilterFilesByMask(tc.FilePaths, tc.Mask)

			if len(result) != tc.ExpectedCount {
				t.Errorf("Expected %d files, got %d: %v", tc.ExpectedCount, len(result), result)
			}

			// Проверяем что все ожидаемые файлы присутствуют
			resultMap := make(map[string]bool)
			for _, f := range result {
				resultMap[f] = true
			}
			for _, expected := range tc.ExpectedFiles {
				if !resultMap[expected] {
					t.Errorf("Expected file '%s' not found in result", expected)
				}
			}
		})
	}
}
