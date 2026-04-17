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

package commands

import (
	"gitlabFileScanner/internal/cli/formatter"

	"github.com/spf13/cobra"
)

// outputter defines the interface for command output.
type outputter interface {
	Print(v any)
	Println(msg string)
	Printf(format string, args ...any)
	Error(err error)
	Log(msg string)
}

// getFormatter returns a Formatter based on the --json flag.
func getFormatter(cmd *cobra.Command) (outputter, error) {
	jsonMode, err := cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}

	return formatter.New(jsonMode), nil
}
