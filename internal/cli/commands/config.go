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
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gitlabFileScanner/internal/domain"

	"github.com/spf13/cobra"
)

// ScanFlags holds parsed CLI flags for the scan command.
type ScanFlags struct {
	GitLabURL       string
	GitLabToken     string
	GitLabTokenFile string
	GitLabBranch    string
	ProjectID       int
	ProjectsLimit   int
	ExportPath      string
	FilesMask       string
}

// BindScanFlags registers scan flags on the command.
func BindScanFlags(cmd *cobra.Command) {
	cmd.Flags().String("url", "", "GitLab server URL")
	cmd.Flags().String("token", "", "GitLab API token")
	cmd.Flags().String("token-file", "", "Path to a file containing the GitLab API token")
	cmd.Flags().String("branch", "", "GitLab branch to scan")
	cmd.Flags().Int("project-id", 0, "Single GitLab project ID")
	cmd.Flags().Int("projects-limit", 0, "Maximum number of projects to scan")
	cmd.Flags().String("export-files-path", "", "Output directory for exported JSON files")
	cmd.Flags().String("files-mask", "", "File mask for filtering exported files")
}

// LoadScanFlags reads flags from CLI and environment variables with defaults.
func LoadScanFlags(cmd *cobra.Command) ScanFlags {
	f := ScanFlags{
		GitLabURL:       getEnv("GFS_URL"),
		GitLabToken:     getEnv("GFS_API_TOKEN"),
		GitLabTokenFile: getEnv("GFS_API_TOKEN_FILE"),
		GitLabBranch:    getEnvOr("GFS_BRANCH", "main"),
		ExportPath:      getEnv("GFS_EXPORT_PATH"),
		FilesMask:       getEnvOr("GFS_FILEMASK", "*"),
	}

	if v := getEnv("GFS_PROJECT_ID"); v != "" {
		f.ProjectID, _ = strconv.Atoi(v)
	}

	if v := getEnv("GFS_PROJECTS_LIMIT"); v != "" {
		f.ProjectsLimit, _ = strconv.Atoi(v)
	}

	applyFlagOverrides(cmd, &f)

	if f.ProjectsLimit <= 0 {
		f.ProjectsLimit = 100
	}

	return f
}

// ToDomain converts ScanFlags to domain.Config, validating required fields.
func (f ScanFlags) ToDomain() (domain.Config, error) {
	cfg := domain.DefaultConfig()

	if f.GitLabURL == "" {
		return cfg, fmt.Errorf("GFS_URLis required")
	}

	u, err := url.Parse(f.GitLabURL)
	if err != nil {
		return cfg, fmt.Errorf("invalid GitLab URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return cfg, fmt.Errorf("invalid GitLab URL: expected scheme and host (e.g. https://gitlab.com)")
	}

	if f.GitLabBranch == "" {
		return cfg, fmt.Errorf("GFS_BRANCH is required")
	}

	if f.ExportPath == "" {
		return cfg, fmt.Errorf("GFS_EXPORT_PATH is required")
	}

	token, err := resolveToken(f.GitLabToken, f.GitLabTokenFile)
	if err != nil {
		return cfg, err
	}

	cfg.GitLabURL = f.GitLabURL
	cfg.GitLabToken = token
	cfg.GitLabBranch = f.GitLabBranch
	cfg.ExportPath = f.ExportPath
	cfg.FilesMask = f.FilesMask
	cfg.ProjectsLimit = f.ProjectsLimit

	if f.ProjectID > 0 {
		cfg.ProjectIDs = []int{f.ProjectID}
	}

	return cfg, nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvOr(key, def string) string {
	if v := getEnv(key); v != "" {
		return v
	}
	return def
}

func applyFlagOverrides(cmd *cobra.Command, f *ScanFlags) {
	if cmd == nil {
		return
	}

	if cmd.Flags().Changed("url") {
		f.GitLabURL, _ = cmd.Flags().GetString("url")
	}
	if cmd.Flags().Changed("token") {
		f.GitLabToken, _ = cmd.Flags().GetString("token")
	}
	if cmd.Flags().Changed("token-file") {
		f.GitLabTokenFile, _ = cmd.Flags().GetString("token-file")
	}
	if cmd.Flags().Changed("branch") {
		f.GitLabBranch, _ = cmd.Flags().GetString("branch")
	}
	if cmd.Flags().Changed("project-id") {
		f.ProjectID, _ = cmd.Flags().GetInt("project-id")
	}
	if cmd.Flags().Changed("projects-limit") {
		f.ProjectsLimit, _ = cmd.Flags().GetInt("projects-limit")
	}
	if cmd.Flags().Changed("export-files-path") {
		f.ExportPath, _ = cmd.Flags().GetString("export-files-path")
	}
	if cmd.Flags().Changed("files-mask") {
		f.FilesMask, _ = cmd.Flags().GetString("files-mask")
	}
}

func resolveToken(token, tokenFile string) (string, error) {
	if token != "" {
		return token, nil
	}
	if tokenFile == "" {
		return "", nil
	}

	content, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("reading token file: %w", err)
	}

	resolved := strings.TrimSpace(string(content))
	if resolved == "" {
		return "", fmt.Errorf("token file is empty: %s", tokenFile)
	}

	return resolved, nil
}
