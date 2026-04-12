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

// Package scanner provides the usecase for scanning GitLab projects and exporting file lists.
package scanner

import (
	"context"
	"fmt"
	"time"

	"gitlabFileScanner/internal/domain"
)

// Service orchestrates the scanning workflow.
type Service struct {
	scanner  domain.Scanner
	filter   domain.Filter
	exporter domain.Exporter
	logger   domain.Logger
	cfg      domain.Config
}

// New creates a new scanning service.
func New(s domain.Scanner, f domain.Filter, e domain.Exporter, l domain.Logger, cfg domain.Config) *Service {
	return &Service{
		scanner:  s,
		filter:   f,
		exporter: e,
		logger:   l,
		cfg:      cfg,
	}
}

// Result holds the outcome of scanning a single project.
type Result struct {
	ProjectID   int64
	ProjectName string
	TotalFiles  int
	MatchedFiles int
	SavedPath   string
	Err         error
}

// Run executes the full scanning workflow.
func (s *Service) Run(ctx context.Context) ([]Result, error) {
	startTime := time.Now()

	s.logger.Info("\n[Project Scan]\n")
	projects, err := s.scanner.GetProjects(s.cfg.ProjectsLimit, s.cfg.ProjectIDs...)
	if err != nil {
		return nil, fmt.Errorf("getting projects: %w", err)
	}

	s.logger.Info("Found %d projects\n", len(projects))

	results := make([]Result, 0, len(projects))
	total := len(projects)

	for i, project := range projects {
		select {
		case <-ctx.Done():
			s.logger.Info("Cancelled after %d/%d projects", i, total)
			return results, ctx.Err()
		default:
		}

		res := s.scanProject(i+1, total, project)
		results = append(results, res)
	}

	s.logger.Info("Completed in %s\n", s.durationString(time.Since(startTime)))
	return results, nil
}

func (s *Service) scanProject(num, total int, project domain.Project) Result {
	res := Result{
		ProjectID:   project.ID,
		ProjectName: project.Name,
	}

	s.logger.Info("[%d/%d] %d | %s", num, total, project.ID, project.Name)

	files, err := s.scanner.GetFilePaths(project.ID, s.cfg.GitLabBranch)
	if err != nil {
		res.Err = err
		s.logger.Info(" | %v\n", err)
		return res
	}

	filtered := s.filter.Apply(files, s.cfg.FilesMask)
	res.TotalFiles = len(files)
	res.MatchedFiles = len(filtered)

	if len(filtered) == 0 {
		s.logger.Info(" | %d/%d | no files matching mask\n", len(files), len(filtered))
		return res
	}

	fileList := &domain.FileList{
		Name:      project.Name,
		WebURL:    project.WebURL,
		ID:        project.ID,
		Branch:    s.cfg.GitLabBranch,
		FilePaths: filtered,
	}

	path, err := s.exporter.Save(s.cfg.ExportPath, fileList)
	if err != nil {
		res.Err = err
		s.logger.Info(" | save error: %v\n", err)
		return res
	}

	res.SavedPath = path
	s.logger.Info(" | %d/%d | %s\n", len(files), len(filtered), path)
	return res
}

func (s *Service) durationString(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	if len(parts) == 0 {
		return "0s"
	}
	return fmt.Sprintf("%s", parts)
}
