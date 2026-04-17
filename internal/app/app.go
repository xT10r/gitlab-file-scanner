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

// Package app wires dependencies and runs the application.
package app

import (
	"context"
	"fmt"

	"gitlabFileScanner/internal/domain"
	"gitlabFileScanner/internal/infrastructure/consolelog"
	"gitlabFileScanner/internal/infrastructure/filefilter"
	fs "gitlabFileScanner/internal/infrastructure/filesystem"
	"gitlabFileScanner/internal/infrastructure/gitlab"
	"gitlabFileScanner/internal/usecase/scanner"
)

// NewService creates a fully wired scanning service.
func NewService(ctx context.Context, cfg domain.Config) (*scanner.Service, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := gitlab.NewClient(ctx, cfg.GitLabURL, cfg.GitLabToken)
	if err != nil {
		return nil, fmt.Errorf("creating GitLab client: %w", err)
	}

	return scanner.New(
		client,
		filefilter.New(),
		fs.NewWriter(),
		consolelog.New(),
		cfg,
	), nil
}
