package text_test

import (
	"gitlabFileScanner/internal/text"
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
		result := text.StringPtr(tc.Input)

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
		result := text.BoolPtr(tc.Input)

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
			Expected:    "",
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
			result := text.GetDurationString(tc.Duration)
			if result != tc.Expected {
				t.Errorf("Test case '%s' failed: expected %v, got %v", tc.Description, tc.Expected, result)
			}
		})
	}
}
