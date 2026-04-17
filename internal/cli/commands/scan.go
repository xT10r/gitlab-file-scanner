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
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"gitlabFileScanner/internal/app"
	"gitlabFileScanner/internal/usecase/scanner"

	"github.com/spf13/cobra"
)

// scanOutput holds scan command output data.
type scanOutput struct {
	Version        string          `json:"version"`
	Duration       string          `json:"duration"`
	Total          int             `json:"total_projects"`
	Success        int             `json:"successful"`
	Failed         int             `json:"failed"`
	ProjectResults []projectResult `json:"projects"`
}

type projectResult struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Matched int    `json:"files_matched"`
	Total   int    `json:"total_files"`
	Path    string `json:"output_path,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan GitLab projects for files",
		Long: `Scans GitLab projects and exports file lists to JSON.

Environment variables:
  GFS_URL                 GitLab server URL (required)
  GFS_API_TOKEN           API token
  GFS_API_TOKEN_FILE      Path to a file containing the API token
  GFS_BRANCH              Branch to scan (default: main)
  GFS_PROJECT_ID          Single project ID
  GFS_PROJECTS_LIMIT      Max projects (default: 100)
  GFS_EXPORT_PATH         Output directory (required)
  GFS_FILEMASK            File mask (default: *)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			f, err := getFormatter(cmd)
			if err != nil {
				return err
			}

			startTime := time.Now()

			flags := LoadScanFlags(cmd)

			cfg, err := flags.ToDomain()
			if err != nil {
				f.Error(err)
				return err
			}

			ctx, cancel := signal.NotifyContext(
				context.Background(), syscall.SIGINT, syscall.SIGTERM,
			)
			defer cancel()

			svc, err := app.NewService(ctx, cfg)
			if err != nil {
				f.Error(err)
				return err
			}

			results, err := svc.Run(ctx)
			if err != nil {
				f.Error(fmt.Errorf("scan failed: %w", err))
				return err
			}

			duration := time.Since(startTime)
			out := buildScanOutput(version, duration, results)

			f.Print(out)

			return nil
		},
	}

	BindScanFlags(cmd)

	return cmd
}

func buildScanOutput(ver string, dur time.Duration, results []scanner.Result) scanOutput {
	out := scanOutput{
		Version:  ver,
		Duration: dur.Round(time.Second).String(),
	}

	for _, r := range results {
		pr := projectResult{
			ID:      r.ProjectID,
			Name:    r.ProjectName,
			Total:   r.TotalFiles,
			Matched: r.MatchedFiles,
		}

		if r.Err != nil {
			pr.Status = "error"
			pr.Error = r.Err.Error()
			out.Failed++
		} else {
			pr.Status = "success"
			pr.Path = r.SavedPath
			out.Success++
		}

		out.ProjectResults = append(out.ProjectResults, pr)
		out.Total++
	}

	return out
}
