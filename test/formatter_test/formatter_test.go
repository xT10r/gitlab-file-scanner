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

package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"gitlabFileScanner/internal/cli/formatter"
)

func TestFormatter_Print_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewJSON(&buf, &bytes.Buffer{})

	data := map[string]string{"key": "value"}
	f.Print(data)

	var parsed map[string]string
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["key"] != "value" {
		t.Errorf("expected key=value, got %v", parsed)
	}
}

func TestFormatter_PrintText(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewText(&buf, &bytes.Buffer{})

	f.Println("hello")

	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected 'hello' in output, got %q", buf.String())
	}
}

func TestFormatter_PrintTable_Text(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewText(&buf, &bytes.Buffer{})

	f.PrintTable(map[string]string{"version": "1.0.0", "commit": "abc"})

	out := buf.String()
	if !strings.Contains(out, "version:") || !strings.Contains(out, "1.0.0") {
		t.Errorf("expected table output, got %q", out)
	}
}

func TestFormatter_PrintTable_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewJSON(&buf, &bytes.Buffer{})

	f.PrintTable(map[string]string{"version": "1.0.0"})

	var parsed map[string]string
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["version"] != "1.0.0" {
		t.Errorf("expected version=1.0.0, got %v", parsed)
	}
}

func TestFormatter_Error_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewJSON(&bytes.Buffer{}, &buf)

	f.Error(testErr)

	if !strings.Contains(buf.String(), testErr.Error()) {
		t.Errorf("expected error in stderr, got %q", buf.String())
	}
}

func TestFormatter_Error_Text(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewText(&bytes.Buffer{}, &buf)

	f.Error(testErr)

	if !strings.Contains(buf.String(), "Error:") {
		t.Errorf("expected 'Error:' prefix, got %q", buf.String())
	}
}

func TestFormatter_Error_Nil(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewJSON(&bytes.Buffer{}, &buf)

	f.Error(nil)

	if buf.Len() > 0 {
		t.Errorf("expected no output for nil error, got %q", buf.String())
	}
}

func TestFormatter_Log_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewJSON(&bytes.Buffer{}, &buf)

	f.Log("scanning...")

	if !strings.Contains(buf.String(), "scanning...") {
		t.Errorf("expected log in stderr for JSON mode, got %q", buf.String())
	}
}

func TestFormatter_Log_Text(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	f := formatter.NewText(&buf, &bytes.Buffer{})

	f.Log("scanning...")

	if !strings.Contains(buf.String(), "scanning...") {
		t.Errorf("expected log in stdout for text mode, got %q", buf.String())
	}
}

type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

var testErr = &customErr{msg: "test error"}
