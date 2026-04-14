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

func TestStringPtr(t *testing.T) {

	testCases := []struct {
		Description string
		Input       string
		Expected    interface{}
	}{
		{Description: "Test with string input", Input: "test", Expected: "test"},
		{Description: `Test with empty string input`, Input: "", Expected: ""},
	}

	for _, tc := range testCases {
		result := strutil.StringPtr(tc.Input)

		if *result != tc.Expected.(string) {
			t.Errorf("Test case '%s' failed: expected %v, got %v", tc.Description, tc.Expected, result)
		}
	}
}

func TestBoolPtr(t *testing.T) {
	testCases := []struct {
		Description string
		Input       bool
		Expected    interface{}
	}{
		{Description: "Test with true input", Input: true, Expected: true},
		{Description: "Test with false input", Input: false, Expected: false},
	}

	for _, tc := range testCases {
		result := strutil.BoolPtr(tc.Input)

		if *result != tc.Expected.(bool) {
			t.Errorf("Test case '%q' failed: expected %v, got %v", tc.Description, tc.Expected, *result)
		}
	}
}

func TestGetDurationString(t *testing.T) {
	testCases := []struct {
		Description string
		Duration    time.Duration
		Expected    string
	}{
		{
			Description: "Test with zero duration",
			Duration:    0,
			Expected:    "0s",
		},
		{
			Description: "Test with 1 hour duration",
			Duration:    time.Hour,
			Expected:    "1 час",
		},
		{
			Description: "Test with 1 day, 2 hours, 30 minutes and 15 seconds",
			Duration:    1*time.Hour*24 + 2*time.Hour + 30*time.Minute + 15*time.Second,
			Expected:    "1 день 2 часа 30 минут 15 секунд",
		},
		{
			Description: "Test with 0 days, 5 hours, 0 minutes and 0 seconds",
			Duration:    5 * time.Hour,
			Expected:    "5 часов",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()
			t.Helper()
			result := strutil.GetDurationString(tc.Duration)
			if result != tc.Expected {
				t.Errorf("Test case '%s' failed: expected %v, got %v", tc.Description, tc.Expected, result)
			}
		})
	}
}
