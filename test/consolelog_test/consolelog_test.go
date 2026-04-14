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

package consolelog_test

import (
	"testing"

	"gitlabFileScanner/internal/infrastructure/consolelog"
)

func TestLogger_Info_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	// These should not panic
	l.Info("simple message")
	l.Info("formatted: %s", "value")
	l.Info("")
	l.Info("many args: %s %d %v", "str", 42, true)
}

func TestLogger_Error_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Error("error message")
	l.Error("error with detail: %v", 42)
	l.Error("")
}

func TestLogger_Warn_NoPanic(t *testing.T) {
	t.Parallel()

	l := consolelog.New()

	l.Warn("warning message")
	l.Warn("warning with detail: %s", "detail")
	l.Warn("")
}
