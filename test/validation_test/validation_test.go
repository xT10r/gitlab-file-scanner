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

package validation_test

import (
	"testing"

	"gitlabFileScanner/internal/cli/commands"
	"gitlabFileScanner/internal/strutil"
)

// Тесты валидации URL (через ScanFlags.ToDomain).
func TestGitLabURLValidation(t *testing.T) {
	testCases := []struct {
		Description string
		URL         string
		ShouldPass  bool
	}{
		{Description: "Valid HTTPS URL", URL: "https://gitlab.com", ShouldPass: true},
		{Description: "Valid HTTP URL", URL: "http://localhost:8080", ShouldPass: true},
		{Description: "Empty URL", URL: "", ShouldPass: false},
		{Description: "URL without scheme", URL: "gitlab.com", ShouldPass: false},
		{Description: "URL with path", URL: "https://gitlab.com/api/v4", ShouldPass: true},
		{Description: "URL with query", URL: "https://gitlab.com?token=abc", ShouldPass: true},
		{Description: "Invalid URL with spaces", URL: "https://gitlab .com", ShouldPass: false},
		{Description: "Just scheme", URL: "https://", ShouldPass: false},
		{Description: "Malformed URL", URL: "://broken", ShouldPass: false},
		{Description: "URL with unicode", URL: "https://gitlab.com/тест", ShouldPass: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			flags := commands.ScanFlags{
				GitLabURL:     tc.URL,
				GitLabBranch:  "main",
				ExportPath:    "/tmp/output",
				FilesMask:     "*",
				ProjectsLimit: 100,
			}

			_, err := flags.ToDomain()
			passed := err == nil

			if passed != tc.ShouldPass {
				t.Errorf("URL '%s': expected pass=%v, got %v (err=%v)", tc.URL, tc.ShouldPass, passed, err)
			}
		})
	}
}

// Тесты валидации масок (через strutil.MaskToFileRegex)
func TestFileMaskValidation(t *testing.T) {
	testCases := []struct {
		Description string
		Mask        string
		ShouldPass  bool
	}{
		{Description: "Simple wildcard", Mask: "*", ShouldPass: true},
		{Description: "Extension mask", Mask: "*.go", ShouldPass: true},
		{Description: "Multiple masks", Mask: "*.go|*.py", ShouldPass: true},
		{Description: "Group mask", Mask: "*.(go|py)", ShouldPass: true},
		{Description: "Empty mask", Mask: "", ShouldPass: true},
		{Description: "Mask with spaces", Mask: "*.go | *.py", ShouldPass: true},
		{Description: "Complex nested groups", Mask: "*.(go|py)|*.(js|ts)", ShouldPass: true},
		{Description: "Mask with special regex chars", Mask: "file[1].go", ShouldPass: true}, // regex special chars become part of regex
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			_, err := strutil.MaskToFileRegex(tc.Mask)
			passed := err == nil

			if passed != tc.ShouldPass {
				t.Errorf("Mask '%s': expected pass=%v, got %v (err=%v)", tc.Mask, tc.ShouldPass, passed, err)
			}
		})
	}
}
