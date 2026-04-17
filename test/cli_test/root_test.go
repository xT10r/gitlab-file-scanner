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

package cli_test

import (
	"testing"

	"gitlabFileScanner/internal/cli"
)

func TestRootCmd_HasVersionCommand(t *testing.T) {
	t.Parallel()

	root := cli.NewRootCmd()
	sub := root.Commands()

	found := false
	for _, cmd := range sub {
		if cmd.Name() == "version" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'version' subcommand to be registered")
	}
}

func TestRootCmd_HasScanCommand(t *testing.T) {
	t.Parallel()

	root := cli.NewRootCmd()
	sub := root.Commands()

	found := false
	for _, cmd := range sub {
		if cmd.Name() == "scan" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected 'scan' subcommand to be registered")
	}
}

func TestRootCmd_HasJSONFlag(t *testing.T) {
	t.Parallel()

	root := cli.NewRootCmd()

	flag := root.PersistentFlags().Lookup("json")
	if flag == nil {
		t.Fatal("expected --json flag to be defined")
	}

	if flag.Shorthand != "j" {
		t.Errorf("expected shorthand 'j', got %q", flag.Shorthand)
	}
}

func TestRootCmd_VersionCommand_Output(t *testing.T) {
	t.Parallel()

	// Version command writes directly to os.Stdout/os.Stderr,
	// not through cobra's out/err. Just verify it doesn't error.
	root := cli.NewRootCmd()
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCmd_VersionCommand_JSON(t *testing.T) {
	t.Parallel()

	root := cli.NewRootCmd()
	root.SetArgs([]string{"version", "--json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCmd_VersionCommand_InvalidFlag(t *testing.T) {
	t.Parallel()

	root := cli.NewRootCmd()
	root.SetArgs([]string{"version", "--unknown-flag"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for unknown flag")
	}
}
