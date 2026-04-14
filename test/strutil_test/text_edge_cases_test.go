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

func TestMaskToFileRegexEdgeCases(t *testing.T) {
	testCases := []struct {
		Description string
		Mask        string
		ShouldMatch []string
		ShouldNot   []string
		ExpectError bool
	}{
		{
			Description: "Empty mask matches nothing",
			Mask:        "",
			ShouldMatch: []string{""},
			ShouldNot:   []string{"file.go"},
			ExpectError: false,
		},
		{
			Description: "Mask with spaces gets trimmed",
			Mask:        "*.go | *.py",
			ShouldMatch: []string{"main.go", "script.py"},
			ShouldNot:   []string{"main.js"},
			ExpectError: false,
		},
		{
			Description: "Multiple wildcards",
			Mask:        "*test*",
			ShouldMatch: []string{"mytest.go", "test.go", "testing_file.go"},
			ShouldNot:   []string{"main.go"},
			ExpectError: false,
		},
		{
			Description: "Only wildcard star",
			Mask:        "*",
			ShouldMatch: []string{"", "a", "anything.txt", "path/to/file.go"},
			ShouldNot:   []string{},
			ExpectError: false,
		},
		{
			Description: "Double star matches path",
			Mask:        "**",
			ShouldMatch: []string{"a", "anything", "path/to/file"},
			ShouldNot:   []string{},
			ExpectError: false,
		},
		{
			Description: "Complex nested groups",
			Mask:        "*.(go|py|js)",
			ShouldMatch: []string{"main.go", "app.py", "index.js"},
			ShouldNot:   []string{"main.rs", "file.txt"},
			ExpectError: false,
		},
		{
			Description: "Dot in filename",
			Mask:        "*.tar.gz",
			ShouldMatch: []string{"archive.tar.gz", "backup.tar.gz"},
			ShouldNot:   []string{"file.tar.bz2", "file.gz"},
			ExpectError: false,
		},
		{
			Description: "Exact filename",
			Mask:        "Makefile",
			ShouldMatch: []string{"Makefile"},
			ShouldNot:   []string{"makefile", "Makefile.txt", "MyMakefile"},
			ExpectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			re, err := strutil.MaskToFileRegex(tc.Mask)
			if tc.ExpectError {
				if err == nil {
					t.Errorf("Expected error for mask '%s', but got none", tc.Mask)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for mask '%s': %v", tc.Mask, err)
				return
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

func TestSplitMaskEdgeCases(t *testing.T) {
	testCases := []struct {
		Description string
		Mask        string
		Expected    []string
	}{
		{
			Description: "Empty string",
			Mask:        "",
			Expected:    []string{""},
		},
		{
			Description: "Only pipes",
			Mask:        "||",
			Expected:    []string{"", "", ""},
		},
		{
			Description: "Pipe at start",
			Mask:        "|*.go",
			Expected:    []string{"", "*.go"},
		},
		{
			Description: "Pipe at end",
			Mask:        "*.go|",
			Expected:    []string{"*.go", ""},
		},
		{
			Description: "Nested groups",
			Mask:        "*.(go|py)|*.(js|ts)",
			Expected:    []string{"*.(go|py)", "*.(js|ts)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result := strutil.SplitMask(tc.Mask)
			if len(result) != len(tc.Expected) {
				t.Errorf("Expected %d parts, got %d: %v", len(tc.Expected), len(result), result)
				return
			}
			for i, part := range tc.Expected {
				if result[i] != part {
					t.Errorf("Part %d: expected '%s', got '%s'", i, part, result[i])
				}
			}
		})
	}
}

func TestGetDurationStringEdgeCases(t *testing.T) {
	testCases := []struct {
		Description string
		Duration    time.Duration
		NonEmpty    bool
	}{
		{Description: "Negative duration", Duration: -5 * time.Second, NonEmpty: true},
		{Description: "Very large duration (1 year)", Duration: 365 * 24 * time.Hour, NonEmpty: true},
		{Description: "Exactly 1 day", Duration: 24 * time.Hour, NonEmpty: true},
		{Description: "Exactly 1 hour", Duration: time.Hour, NonEmpty: true},
		{Description: "Exactly 1 minute", Duration: time.Minute, NonEmpty: true},
		{Description: "1 second", Duration: time.Second, NonEmpty: true},
		{Description: "Sub-second duration", Duration: 500 * time.Millisecond, NonEmpty: true},
		{Description: "Mixed: 1 hour 1 minute 1 second", Duration: time.Hour + time.Minute + time.Second, NonEmpty: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result := strutil.GetDurationString(tc.Duration)

			if tc.NonEmpty && result == "" {
				t.Error("Expected non-empty string, got empty")
			}
		})
	}
}
