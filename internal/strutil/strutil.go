// Copyright 2024 Alex Dobshikov
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

// Package strutil provides string manipulation and formatting utilities.
package strutil

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the given bool.
func BoolPtr(b bool) *bool {
	return &b
}

// MaskToFileRegex converts a glob-style file mask to a compiled regular expression.
//
// Supported patterns:
//   - * matches any characters
//   - . is escaped to match literal dots
//   - () groups (converted to non-capturing)
//   - | alternation (outside of groups)
func MaskToFileRegex(mask string) (*regexp.Regexp, error) {
	mask = strings.ReplaceAll(mask, " ", "")

	// Escape all regex meta-characters except *, |, (, )
	regex := strings.ReplaceAll(mask, `\`, `\\`)
	regex = strings.ReplaceAll(regex, `[`, `\[`)
	regex = strings.ReplaceAll(regex, `]`, `\]`)
	regex = strings.ReplaceAll(regex, `+`, `\+`)
	regex = strings.ReplaceAll(regex, `?`, `\?`)
	regex = strings.ReplaceAll(regex, `{`, `\{`)
	regex = strings.ReplaceAll(regex, `}`, `\}`)
	regex = strings.ReplaceAll(regex, `^`, `\^`)
	regex = strings.ReplaceAll(regex, `$`, `\$`)

	// Now apply glob patterns
	regex = strings.ReplaceAll(regex, ".", `\.`)
	regex = strings.ReplaceAll(regex, "*", ".*")
	regex = strings.ReplaceAll(regex, "(", "(?:")

	return regexp.Compile("^" + regex + "$")
}

// SplitMask splits a file mask by the '|' separator,
// respecting parentheses nesting.
func SplitMask(mask string) []string {
	var parts []string
	var part strings.Builder
	var level int

	for _, char := range mask {
		switch char {
		case '|':
			if level == 0 {
				parts = append(parts, part.String())
				part.Reset()
			} else {
				part.WriteRune(char)
			}
		case '(':
			level++
			part.WriteRune(char)
		case ')':
			level--
			part.WriteRune(char)
		default:
			part.WriteRune(char)
		}
	}

	return append(parts, part.String())
}

// GetDurationString formats a duration into a human-readable Russian string.
func GetDurationString(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, pluralize(days, "день", "дня", "дней"))
	}
	if hours > 0 {
		parts = append(parts, pluralize(hours, "час", "часа", "часов"))
	}
	if minutes > 0 {
		parts = append(parts, pluralize(minutes, "минута", "минуты", "минут"))
	}
	if seconds > 0 {
		parts = append(parts, pluralize(seconds, "секунда", "секунды", "секунд"))
	}

	if len(parts) == 0 {
		return "0s"
	}
	return strings.Join(parts, " ")
}

func pluralize(count int, form1, form2, form5 string) string {
	if count == 0 {
		return ""
	}
	abs := count
	if abs < 0 {
		abs = -abs
	}
	lastDigit := abs % 10
	lastTwoDigits := abs % 100

	switch {
	case lastTwoDigits >= 11 && lastTwoDigits <= 14:
		return fmt.Sprintf("%d %s", count, form5)
	case lastDigit == 1:
		return fmt.Sprintf("%d %s", count, form1)
	case lastDigit >= 2 && lastDigit <= 4:
		return fmt.Sprintf("%d %s", count, form2)
	default:
		return fmt.Sprintf("%d %s", count, form5)
	}
}

// SortFilePaths sorts file paths by their base name.
// Deprecated: use sort.Strings for alphabetical sorting.
func SortFilePaths(filePaths []string) {
	sort.SliceStable(filePaths, func(i, j int) bool {
		return filepath.Base(filePaths[i]) < filepath.Base(filePaths[j])
	})
}
