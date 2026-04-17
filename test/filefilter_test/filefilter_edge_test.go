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

package filefilter_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/filefilter"
)

func TestFilter_Deduplication(t *testing.T) {
	t.Parallel()

	f := filefilter.New()

	// Маска с пересекающимися паттернами - один файл может совпасть дважды
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

	// *.go и *go - оба совпадут с main.go
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

	// [ ] должны экранироваться - не быть character class
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
