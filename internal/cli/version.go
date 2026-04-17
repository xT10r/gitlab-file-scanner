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

package cli

import "fmt"

var (
	Version   = "unknown"
	Commit    = "unknown"
	Branch    = "unknown"
	BuildTime = "unknown"
)

type VersionResult struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	BuildTime string `json:"buildTime"`
	Resolved  string `json:"resolved"`
}

func versionResult() VersionResult {
	return VersionResult{
		Version:   Version,
		Commit:    Commit,
		Branch:    Branch,
		BuildTime: BuildTime,
		Resolved:  resolveVersion(),
	}
}

func resolveVersion() string {
	if Version != "" && Version != "unknown" {
		return Version
	}

	if Branch != "unknown" && Commit != "unknown" {
		return fmt.Sprintf("%s-%s", Branch, Commit[:7])
	}

	if BuildTime != "unknown" {
		return "build-" + BuildTime
	}

	return "unknown"
}

/*
var Version = "dev"

type VersionCommand struct{}

func NewVersionCommand() *VersionCommand {
    return &VersionCommand{}
}

func (c *VersionCommand) Name() string {
    return "version"
}

func (c *VersionCommand) Description() string {
    return "Show application version"
}

func (c *VersionCommand) Run(args []string) error {
    fmt.Println("gitlab-file-scanner version:", Version)
    return nil
}
*/
