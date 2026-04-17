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

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"gitlabFileScanner/internal/cli/commands"

	"github.com/spf13/cobra"
)

func TestLoadScanFlags_Defaults(t *testing.T) {
	// Cannot use t.Parallel with t.Setenv

	// Clear env vars that might affect the test
	t.Setenv("GFS_URL", "https://gitlab.com")
	t.Setenv("GFS_BRANCH", "")
	t.Setenv("GFS_EXPORT_PATH", "/tmp/output")
	t.Setenv("GFS_PROJECTS_LIMIT", "")
	t.Setenv("GFS_PROJECT_ID", "")
	t.Setenv("GFS_FILEMASK", "")
	t.Setenv("GFS_API_TOKEN", "")

	flags := commands.LoadScanFlags(nil)

	if flags.ProjectsLimit != 100 {
		t.Errorf("expected default limit 100, got %d", flags.ProjectsLimit)
	}

	if flags.FilesMask != "*" {
		t.Errorf("expected default mask '*', got %q", flags.FilesMask)
	}

	if flags.GitLabBranch != "main" {
		t.Errorf("expected default branch 'main', got %q", flags.GitLabBranch)
	}
}

func TestLoadScanFlags_ReadsGFSVariables(t *testing.T) {
	t.Setenv("GFS_URL", "https://short.example")
	t.Setenv("GFS_BRANCH", "release")
	t.Setenv("GFS_EXPORT_PATH", "/tmp/short")
	t.Setenv("GFS_API_TOKEN", "short-token")

	flags := commands.LoadScanFlags(nil)

	if flags.GitLabURL != "https://short.example" {
		t.Fatalf("expected GFS URL, got %q", flags.GitLabURL)
	}
	if flags.GitLabBranch != "release" {
		t.Fatalf("expected GFS branch, got %q", flags.GitLabBranch)
	}
	if flags.ExportPath != "/tmp/short" {
		t.Fatalf("expected GFS export path, got %q", flags.ExportPath)
	}
	if flags.GitLabToken != "short-token" {
		t.Fatalf("expected GFS token, got %q", flags.GitLabToken)
	}
}

func TestLoadScanFlags_CommandFlagsOverrideEnv(t *testing.T) {
	t.Setenv("GFS_URL", "https://env.example")
	t.Setenv("GFS_BRANCH", "env-branch")
	t.Setenv("GFS_EXPORT_PATH", "/tmp/env")
	t.Setenv("GFS_FILEMASK", "*.env")
	t.Setenv("GFS_PROJECT_ID", "10")
	t.Setenv("GFS_PROJECTS_LIMIT", "20")
	t.Setenv("GFS_API_TOKEN", "env-token")
	t.Setenv("GFS_API_TOKEN_FILE", "/tmp/env-token-file")

	cmd := &cobra.Command{Use: "scan"}
	commands.BindScanFlags(cmd)
	if err := cmd.ParseFlags([]string{
		"--url", "https://flag.example",
		"--branch", "flag-branch",
		"--export-files-path", "/tmp/flag",
		"--files-mask", "*.go",
		"--project-id", "42",
		"--projects-limit", "55",
		"--token", "flag-token",
		"--token-file", "/tmp/flag-token-file",
	}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}

	flags := commands.LoadScanFlags(cmd)

	if flags.GitLabURL != "https://flag.example" ||
		flags.GitLabBranch != "flag-branch" ||
		flags.ExportPath != "/tmp/flag" ||
		flags.FilesMask != "*.go" ||
		flags.ProjectID != 42 ||
		flags.ProjectsLimit != 55 ||
		flags.GitLabToken != "flag-token" ||
		flags.GitLabTokenFile != "/tmp/flag-token-file" {
		t.Fatalf("expected CLI flags to override env, got %+v", flags)
	}
}

func TestScanFlags_ToDomain_Valid(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:     "https://gitlab.com",
		GitLabBranch:  "main",
		ExportPath:    "/tmp/output",
		ProjectsLimit: 50,
		FilesMask:     "*.go",
	}

	cfg, err := flags.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.GitLabURL != flags.GitLabURL {
		t.Errorf("URL mismatch: expected %q, got %q", flags.GitLabURL, cfg.GitLabURL)
	}

	if cfg.ProjectsLimit != flags.ProjectsLimit {
		t.Errorf("limit mismatch: expected %d, got %d", flags.ProjectsLimit, cfg.ProjectsLimit)
	}
}

func TestScanFlags_ToDomain_WithTokenFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token.txt")
	if err := os.WriteFile(tokenFile, []byte(" secret-token \n"), 0o600); err != nil {
		t.Fatalf("write token file: %v", err)
	}

	flags := commands.ScanFlags{
		GitLabURL:       "https://gitlab.com",
		GitLabTokenFile: tokenFile,
		GitLabBranch:    "main",
		ExportPath:      "/tmp/output",
	}

	cfg, err := flags.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.GitLabToken != "secret-token" {
		t.Fatalf("expected token from file, got %q", cfg.GitLabToken)
	}
}

func TestScanFlags_ToDomain_TokenWinsOverTokenFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token.txt")
	if err := os.WriteFile(tokenFile, []byte("file-token"), 0o600); err != nil {
		t.Fatalf("write token file: %v", err)
	}

	flags := commands.ScanFlags{
		GitLabURL:       "https://gitlab.com",
		GitLabToken:     "inline-token",
		GitLabTokenFile: tokenFile,
		GitLabBranch:    "main",
		ExportPath:      "/tmp/output",
	}

	cfg, err := flags.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.GitLabToken != "inline-token" {
		t.Fatalf("expected inline token to win, got %q", cfg.GitLabToken)
	}
}

func TestScanFlags_ToDomain_EmptyTokenFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "empty-token.txt")
	if err := os.WriteFile(tokenFile, []byte(" \n\t "), 0o600); err != nil {
		t.Fatalf("write token file: %v", err)
	}

	flags := commands.ScanFlags{
		GitLabURL:       "https://gitlab.com",
		GitLabTokenFile: tokenFile,
		GitLabBranch:    "main",
		ExportPath:      "/tmp/output",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for empty token file")
	}
}

func TestScanFlags_ToDomain_MissingURL(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabBranch: "main",
		ExportPath:   "/tmp/output",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestScanFlags_ToDomain_InvalidURL(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:    "://invalid",
		GitLabBranch: "main",
		ExportPath:   "/tmp/output",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestScanFlags_ToDomain_URLWithoutScheme(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:    "gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "/tmp/output",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for URL without scheme/host")
	}
}

func TestScanFlags_ToDomain_MissingBranch(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "",
		ExportPath:   "/tmp/output",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for missing branch")
	}
}

func TestScanFlags_ToDomain_MissingExportPath(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "",
	}

	_, err := flags.ToDomain()
	if err == nil {
		t.Fatal("expected error for missing export path")
	}
}

func TestScanFlags_ToDomain_WithProjectID(t *testing.T) {
	t.Parallel()

	flags := commands.ScanFlags{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "/tmp/output",
		ProjectID:    12345,
	}

	cfg, err := flags.ToDomain()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.ProjectIDs) != 1 || cfg.ProjectIDs[0] != 12345 {
		t.Errorf("expected ProjectIDs [12345], got %v", cfg.ProjectIDs)
	}
}
