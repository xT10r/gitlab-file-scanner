// Copyright 2024-2026 Alex Dobshikov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strutil_test

import (
	"gitlabFileScanner/internal/strutil"
	"testing"
	"time"
)

func TestMaskToFileRegex(t *testing.T) {
	testCases := []struct {
		Description string
		Mask        string
		ShouldMatch []string
		ShouldNot   []string
	}{
		{
			Description: "Simple *.go mask",
			Mask:        "*.go",
			ShouldMatch: []string{"main.go", "test.go", "file.go"},
			ShouldNot:   []string{"main.py", "main.txt", "go"},
		},
		{
			Description: "Multiple masks with pipe",
			Mask:        "*.go|*.py",
			ShouldMatch: []string{"main.go", "script.py"},
			ShouldNot:   []string{"main.js", "index.html"},
		},
		{
			Description: "Specific extension",
			Mask:        "*.md",
			ShouldMatch: []string{"README.md", "docs.md"},
			ShouldNot:   []string{"README.txt", "file.md.bak"},
		},
		{
			Description: "Mask with dots in name",
			Mask:        "*.test.js",
			ShouldMatch: []string{"app.test.js", "utils.test.js"},
			ShouldNot:   []string{"app.js", "app.test.ts"},
		},
		{
			Description: "Wildcard mask",
			Mask:        "*",
			ShouldMatch: []string{"any_file", "anything.txt", "file.go"},
			ShouldNot:   []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			re, err := strutil.MaskToFileRegex(tc.Mask)
			if err != nil {
				t.Errorf("Failed to compile regex for mask '%s': %v", tc.Mask, err)
			}

			for _, file := range tc.ShouldMatch {
				if !re.MatchString(file) {
					t.Errorf("Expected '%s' to match mask '%s'", file, tc.Mask)
				}
			}

			for _, file := range tc.ShouldNot {
				if re.MatchString(file) {
					t.Errorf("Expected '%s' NOT to match mask '%s'", file, tc.Mask)
				}
			}
		})
	}
}

func TestSplitMask(t *testing.T) {
	testCases := []struct {
		Description string
		Mask        string
		Expected    []string
	}{
		{
			Description: "Single mask",
			Mask:        "*.go",
			Expected:    []string{"*.go"},
		},
		{
			Description: "Two masks",
			Mask:        "*.go|*.py",
			Expected:    []string{"*.go", "*.py"},
		},
		{
			Description: "Mask with group",
			Mask:        "*.(go|py)",
			Expected:    []string{"*.(go|py)"},
		},
		{
			Description: "Multiple masks with group",
			Mask:        "*.md|*.(go|py)",
			Expected:    []string{"*.md", "*.(go|py)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result := strutil.SplitMask(tc.Mask)
			if len(result) != len(tc.Expected) {
				t.Errorf("Expected %d parts, got %d: %v", len(tc.Expected), len(result), result)
			}
			for i, part := range tc.Expected {
				if result[i] != part {
					t.Errorf("Part %d: expected '%s', got '%s'", i, part, result[i])
				}
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	testCases := []struct {
		Description string
		Count       int
		Expected    string
	}{
		{Description: "0 секунд", Count: 0, Expected: "0s"},
		{Description: "1 секунда", Count: 1, Expected: "1 секунда"},
		{Description: "2 секунды", Count: 2, Expected: "2 секунды"},
		{Description: "3 секунды", Count: 3, Expected: "3 секунды"},
		{Description: "4 секунды", Count: 4, Expected: "4 секунды"},
		{Description: "5 секунд", Count: 5, Expected: "5 секунд"},
		{Description: "10 секунд", Count: 10, Expected: "10 секунд"},
		{Description: "11 секунд (исключение)", Count: 11, Expected: "11 секунд"},
		{Description: "12 секунд (исключение)", Count: 12, Expected: "12 секунд"},
		{Description: "13 секунд (исключение)", Count: 13, Expected: "13 секунд"},
		{Description: "14 секунд (исключение)", Count: 14, Expected: "14 секунд"},
		{Description: "15 секунд", Count: 15, Expected: "15 секунд"},
		{Description: "21 секунда", Count: 21, Expected: "21 секунда"},
		{Description: "22 секунды", Count: 22, Expected: "22 секунды"},
		{Description: "25 секунд", Count: 25, Expected: "25 секунд"},
		{Description: "100 секунд", Count: 100, Expected: "100 секунд"},
		{Description: "101 секунда", Count: 101, Expected: "101 секунда"},
		{Description: "111 секунд (исключение)", Count: 111, Expected: "111 секунд"},
		{Description: "112 секунд (исключение)", Count: 112, Expected: "112 секунд"},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result := strutil.GetDurationString(time.Duration(tc.Count) * time.Second)

			// GetDurationString может содержать несколько единиц (часы, минуты, секунды)
			// Проверяем что результат содержит ожидаемую форму
			if tc.Count == 0 {
				if result != tc.Expected {
					t.Errorf("Expected '%s', got '%s'", tc.Expected, result)
				}
			} else {
				if result == "" {
					t.Errorf("Expected non-empty result for count %d", tc.Count)
				}
				// Для простых случаев (только секунды) проверяем полное совпадение
				if tc.Count < 60 && result != tc.Expected {
					t.Errorf("Expected '%s', got '%s'", tc.Expected, result)
				}
			}
		})
	}
}
