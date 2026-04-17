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
	"runtime"

	"github.com/spf13/cobra"
)

// Version info set via ldflags.
var (
	version   = "dev"
	commit    = "dev"
	branch    = "unknown"
	buildDate = "unknown"
	goVersion = runtime.Version()
)

// GetVersionMinimal returns the minimal version string.
func GetVersionMinimal() string {
	if version != "dev" {
		return version
	}
	return commit + "-" + branch
}

// GetVersionFull returns the full version string.
func GetVersionFull() string {
	return "gitlab-file-scanner " + GetVersionMinimal() + " (built: " + buildDate + ")"
}

// Version returns the application version.
func Version() string {
	return GetVersionMinimal()
}

func NewVersionCmd() *cobra.Command {
	var full bool

	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print version information",
		RunE: func(c *cobra.Command, _ []string) error {
			f, err := getFormatter(c)
			if err != nil {
				return err
			}

			if f != nil {
				// For JSON, always print full details
				if c.Flag("json").Changed {
					data := map[string]string{
						"version":    version,
						"commit":     commit,
						"branch":     branch,
						"build_date": buildDate,
						"go_version": goVersion,
					}
					f.Print(data)
					return nil
				}

				if full {
					f.Println(GetVersionFull())
				} else {
					f.Println(GetVersionMinimal())
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&full, "full", "V", false, "Print full version info")

	return cmd
}
