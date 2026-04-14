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
	"gitlabFileScanner/internal/strutil"
	"net/url"
	"strconv"
	"testing"
)

// Тесты валидации URL (копия логики из flags.go)
func TestGitLabURLValidation(t *testing.T) {
	testCases := []struct {
		Description string
		URL         string
		ShouldPass  bool
	}{
		{Description: "Valid HTTPS URL", URL: "https://gitlab.com", ShouldPass: true},
		{Description: "Valid HTTP URL", URL: "http://localhost:8080", ShouldPass: true},
		{Description: "Empty URL", URL: "", ShouldPass: true}, // url.Parse accepts empty strings
		{Description: "URL without scheme", URL: "gitlab.com", ShouldPass: true}, // url.Parse accepts it
		{Description: "URL with path", URL: "https://gitlab.com/api/v4", ShouldPass: true},
		{Description: "URL with query", URL: "https://gitlab.com?token=abc", ShouldPass: true},
		{Description: "Invalid URL with spaces", URL: "https://gitlab .com", ShouldPass: false},
		{Description: "Just scheme", URL: "https://", ShouldPass: true},
		{Description: "Malformed URL", URL: "://broken", ShouldPass: false},
		{Description: "URL with unicode", URL: "https://gitlab.com/тест", ShouldPass: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			_, err := url.Parse(tc.URL)
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

// Тесты парсинга чисел (копия логики stringToNumber из flags.go)
func TestStringToNumber(t *testing.T) {
	testCases := []struct {
		Description string
		Input       string
		Expected    int
		ShouldPass  bool
	}{
		{Description: "Positive number", Input: "42", Expected: 42, ShouldPass: true},
		{Description: "Zero", Input: "0", Expected: 0, ShouldPass: true},
		{Description: "Negative number", Input: "-5", Expected: -5, ShouldPass: true},
		{Description: "Large number", Input: "999999999", Expected: 999999999, ShouldPass: true},
		{Description: "Empty string", Input: "", Expected: 0, ShouldPass: false},
		{Description: "Non-numeric string", Input: "abc", Expected: 0, ShouldPass: false},
		{Description: "Float string", Input: "3.14", Expected: 0, ShouldPass: false},
		{Description: "Spaces around number", Input: "  42  ", Expected: 0, ShouldPass: false}, // strconv.Atoi rejects spaces
		{Description: "Plus sign", Input: "+42", Expected: 42, ShouldPass: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()

			result, err := strconv.Atoi(tc.Input)
			passed := err == nil

			if passed != tc.ShouldPass {
				t.Errorf("Input '%s': expected pass=%v, got %v (err=%v)", tc.Input, tc.ShouldPass, passed, err)
			}

			if passed && result != tc.Expected {
				t.Errorf("Input '%s': expected %d, got %d", tc.Input, tc.Expected, result)
			}
		})
	}
}
