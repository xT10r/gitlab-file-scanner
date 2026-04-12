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

package scanner_test

import (
	"context"
	"testing"

	"gitlabFileScanner/internal/domain"
	"gitlabFileScanner/internal/mock"
	scanner "gitlabFileScanner/internal/usecase/scanner"
)

func TestService_Run_Success(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{
			{ID: 1, Name: "project-a", WebURL: "https://gitlab.com/a"},
			{ID: 2, Name: "project-b", WebURL: "https://gitlab.com/b"},
		},
		Files: []string{"main.go", "README.md", "utils.go"},
	}

	mf := &mock.Filter{
		Result: []string{"main.go", "utils.go"},
	}

	me := &mock.Exporter{
		SavedPath: "/output/1.json",
	}

	ml := &mock.Logger{}

	cfg := domain.Config{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "/output",
		FilesMask:    "*.go",
		ProjectsLimit: 10,
	}

	svc := scanner.New(ms, mf, me, ml, cfg)
	results, err := svc.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if !ms.GetProjectsCalled {
		t.Error("expected GetProjects to be called")
	}

	if ms.GetProjectsCalled && !ms.GetFilesCalled {
		t.Error("expected GetFilePaths to be called")
	}

	if !mf.ApplyCalled {
		t.Error("expected Filter.Apply to be called")
	}

	if !me.SaveCalled {
		t.Error("expected Exporter.Save to be called")
	}
}

func TestService_Run_NoProjects(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{},
	}

	svc := scanner.New(ms, &mock.Filter{}, &mock.Exporter{}, &mock.Logger{}, domain.Config{
		GitLabURL:    "https://gitlab.com",
		ProjectsLimit: 10,
	})

	results, err := svc.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestService_Run_GetProjectsError(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		ProjectsErr: context.DeadlineExceeded,
	}

	svc := scanner.New(ms, &mock.Filter{}, &mock.Exporter{}, &mock.Logger{}, domain.Config{
		GitLabURL:    "https://gitlab.com",
		ProjectsLimit: 10,
	})

	_, err := svc.Run(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_Run_GetFilePathsError(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{
			{ID: 1, Name: "project-a"},
		},
		FilesErr: context.DeadlineExceeded,
	}

	svc := scanner.New(ms, &mock.Filter{}, &mock.Exporter{}, &mock.Logger{}, domain.Config{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ProjectsLimit: 10,
	})

	results, err := svc.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if results[0].Err == nil {
		t.Error("expected error in result, got nil")
	}
}

func TestService_Run_EmptyFilterResult(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{
			{ID: 1, Name: "project-a"},
		},
		Files: []string{"main.py", "test.js"},
	}

	mf := &mock.Filter{
		Result: []string{}, // no matches
	}

	me := &mock.Exporter{}
	ml := &mock.Logger{}

	svc := scanner.New(ms, mf, me, ml, domain.Config{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "/output",
		FilesMask:    "*.go",
		ProjectsLimit: 10,
	})

	results, err := svc.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if me.SaveCalled {
		t.Error("expected Save not to be called for empty filter result")
	}
}

func TestService_Run_ContextCancellation(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{
			{ID: 1, Name: "project-a"},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	svc := scanner.New(ms, &mock.Filter{}, &mock.Exporter{}, &mock.Logger{}, domain.Config{
		GitLabURL:    "https://gitlab.com",
		ProjectsLimit: 10,
	})

	results, err := svc.Run(ctx)

	if err == nil {
		t.Fatal("expected context cancelled error, got nil")
	}

	// May or may not have processed any projects depending on timing
	_ = results
}

func TestService_Run_ExportError(t *testing.T) {
	t.Parallel()

	ms := &mock.Scanner{
		Projects: []domain.Project{
			{ID: 1, Name: "project-a"},
		},
		Files: []string{"main.go"},
	}

	mf := &mock.Filter{
		Result: []string{"main.go"},
	}

	me := &mock.Exporter{
		SaveErr: context.DeadlineExceeded,
	}

	ml := &mock.Logger{}

	svc := scanner.New(ms, mf, me, ml, domain.Config{
		GitLabURL:    "https://gitlab.com",
		GitLabBranch: "main",
		ExportPath:   "/output",
		FilesMask:    "*.go",
		ProjectsLimit: 10,
	})

	results, err := svc.Run(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}

	if results[0].Err == nil {
		t.Error("expected error in result, got nil")
	}
}
