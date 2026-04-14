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
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gitlabFileScanner/internal/domain"
	"gitlabFileScanner/internal/infrastructure/consolelog"
	"gitlabFileScanner/internal/infrastructure/filefilter"
	fs "gitlabFileScanner/internal/infrastructure/filesystem"
	"gitlabFileScanner/internal/infrastructure/gitlab"
	"gitlabFileScanner/internal/usecase/scanner"
)

// Start is the application entry point.
func Start() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("\nShutting down...")
	}()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger := consolelog.New()

	logger.Info("\n[Configuration]\n")
	logger.Info("Server URL: %s", cfg.GitLabURL)
	if cfg.GitLabToken != "" {
		logger.Info("API Token: set")
	} else {
		logger.Warn("API Token not set — project visibility may be limited")
	}
	logger.Info("Branch: %s", cfg.GitLabBranch)
	logger.Info("Export Path: %s", cfg.ExportPath)
	logger.Info("Files Mask: %s", cfg.FilesMask)
	logger.Info("Projects Limit: %d", cfg.ProjectsLimit)
	if len(cfg.ProjectIDs) > 0 {
		logger.Info("Project IDs: %v", cfg.ProjectIDs)
	}

	client, err := gitlab.NewClient(ctx, cfg.GitLabURL, cfg.GitLabToken)
	if err != nil {
		return fmt.Errorf("creating GitLab client: %w", err)
	}

	svc := scanner.New(
		client,
		filefilter.New(),
		fs.NewWriter(),
		logger,
		cfg,
	)

	results, err := svc.Run(ctx)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Summary
	errors := 0
	for _, r := range results {
		if r.Err != nil {
			errors++
		}
	}
	if errors > 0 {
		logger.Info("\n%d/%d projects had errors", errors, len(results))
	}

	return nil
}

func loadConfig() (domain.Config, error) {
	cfg := domain.DefaultConfig()

	cfg.GitLabURL = getEnv("GITLAB_FILE_SCANNER_SERVER_URL")
	cfg.GitLabToken = getEnv("GITLAB_FILE_SCANNER_API_TOKEN")
	cfg.GitLabBranch = getEnvOr("GITLAB_FILE_SCANNER_BRANCH", "main")
	cfg.ExportPath = getEnv("GITLAB_FILE_SCANNER_EXPORT_PATH")

	if v := getEnv("GITLAB_FILE_SCANNER_PROJECTS_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("invalid projects limit: %w", err)
		}
		cfg.ProjectsLimit = n
	}

	if v := getEnv("GITLAB_FILE_SCANNER_PROJECT_ID"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return cfg, fmt.Errorf("invalid project ID: %w", err)
		}
		cfg.ProjectIDs = []int{n}
	}

	if v := getEnv("GITLAB_FILE_SCANNER_FILEMASK"); v != "" {
		cfg.FilesMask = v
	}

	if cfg.GitLabURL == "" {
		return cfg, fmt.Errorf("GITLAB_FILE_SCANNER_SERVER_URL is required")
	}
	if cfg.ExportPath == "" {
		return cfg, fmt.Errorf("GITLAB_FILE_SCANNER_EXPORT_PATH is required")
	}

	return cfg, nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getEnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
