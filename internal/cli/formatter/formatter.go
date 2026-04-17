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

// Package formatter provides output formatting for CLI commands.
package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Formatter controls output format for CLI commands.
type Formatter struct {
	jsonMode bool
	stdout   io.Writer
	stderr   io.Writer
}

// New creates a Formatter.
func New(jsonMode bool) *Formatter {
	return &Formatter{
		jsonMode: jsonMode,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}
}

// NewJSON creates a Formatter with JSON output.
func NewJSON(stdout, stderr io.Writer) *Formatter {
	return &Formatter{
		jsonMode: true,
		stdout:   stdout,
		stderr:   stderr,
	}
}

// NewText creates a Formatter with text output.
func NewText(stdout, stderr io.Writer) *Formatter {
	return &Formatter{
		jsonMode: false,
		stdout:   stdout,
		stderr:   stderr,
	}
}

// Print outputs data in the configured format.
func (f *Formatter) Print(data any) {
	if f.jsonMode {
		enc := json.NewEncoder(f.stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			f.Error(fmt.Errorf("encoding output: %w", err))
		}
		return
	}
	// Text mode: print as key-value if it's a map
	if kv, ok := data.(map[string]string); ok {
		f.PrintTable(kv)
		return
	}
	fmt.Fprintf(f.stdout, "%+v\n", data)
}

// PrintTable outputs a simple key-value table in text mode, or JSON in JSON mode.
func (f *Formatter) PrintTable(kv map[string]string) {
	if f.jsonMode {
		f.Print(kv)
		return
	}
	for k, v := range kv {
		fmt.Fprintf(f.stdout, "%-12s %s\n", k+":", v)
	}
}

// Println prints a line in text mode only.
func (f *Formatter) Println(msg string) {
	if f.jsonMode {
		return
	}
	fmt.Fprintln(f.stdout, msg)
}

// Printf prints formatted text in text mode only.
func (f *Formatter) Printf(format string, args ...any) {
	if f.jsonMode {
		return
	}
	fmt.Fprintf(f.stdout, format, args...)
}

// Error prints an error message to stderr.
func (f *Formatter) Error(err error) {
	if err == nil {
		return
	}
	if f.jsonMode {
		_ = json.NewEncoder(f.stderr).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	fmt.Fprintf(f.stderr, "Error: %v\n", err)
}

// Err is an alias for Error.
func (f *Formatter) Err(err error) {
	f.Error(err)
}

// Log prints an informational message to stderr in JSON mode,
// or to stdout in text mode.
func (f *Formatter) Log(msg string) {
	if f.jsonMode {
		fmt.Fprintln(f.stderr, msg)
		return
	}
	fmt.Fprintln(f.stdout, msg)
}
