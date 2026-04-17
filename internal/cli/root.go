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

// Package cli provides CLI commands for the application.
package cli

import (
	"os"

	"gitlabFileScanner/internal/cli/commands"

	"github.com/spf13/cobra"
)

// Execute runs the root CLI command.
func Execute() {
	c := NewRootCmd()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewRootCmd creates the root CLI command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gitlab-file-scanner",
		Short: "CLI tool for scanning GitLab repositories",
		Long: "Scans GitLab projects and exports file lists to JSON. " +
			"Supports filtering by file mask and concurrent scanning.",
		SilenceUsage: true,
		Version:      commands.Version(),
	}

	root.PersistentFlags().BoolP("json", "j", false, "Output in JSON format")
	root.SetVersionTemplate("{{.Version}}\n")

	root.AddCommand(commands.NewVersionCmd())
	root.AddCommand(commands.NewScanCmd())

	return root
}
