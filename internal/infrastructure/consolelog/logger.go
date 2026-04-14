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

// Package consolelog implements domain.Logger with output to stdout.
package consolelog

import "fmt"

// Log level prefixes.
const (
	LevelDebug = "DEBUG: "
	LevelInfo  = "INFO: "
	LevelWarn  = "WARN: "
	LevelError = "ERROR: "
)

// Logger writes to stdout.
type Logger struct{}

// New creates a console logger.
func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...any) {
	prefix := LevelInfo
	if len(args) > 0 {
		fmt.Printf(prefix+msg+"\n", args...)
	} else {
		fmt.Print(prefix + msg + "\n")
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...any) {
	prefix := LevelDebug
	if len(args) > 0 {
		fmt.Printf(prefix+msg+"\n", args...)
	} else {
		fmt.Print(prefix + msg + "\n")
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Printf(LevelError+msg+"\n", args...)
	} else {
		fmt.Print(LevelError + msg + "\n")
	}
}

func (l *Logger) Warn(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Printf(LevelWarn+msg+"\n", args...)
	} else {
		fmt.Print(LevelWarn + msg + "\n")
	}
}
