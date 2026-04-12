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

// Logger writes to stdout.
type Logger struct{}

// New creates a console logger.
func New() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Printf(msg+"\n", args...)
	} else {
		fmt.Println(msg)
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Printf("ERROR: "+msg+"\n", args...)
	} else {
		fmt.Println("ERROR: " + msg)
	}
}

func (l *Logger) Warn(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Printf("WARN: "+msg+"\n", args...)
	} else {
		fmt.Println("WARN: " + msg)
	}
}
