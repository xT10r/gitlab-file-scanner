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

// Package domain defines core business entities and interfaces.
package domain

// Project represents a GitLab project.
type Project struct {
	ID     int64
	Name   string
	WebURL string
}

// FileList represents the result of scanning a project.
type FileList struct {
	Name      string
	WebURL    string
	ID        int64
	Branch    string
	FilePaths []string
}

// Config holds application configuration.
type Config struct {
	GitLabURL     string
	GitLabToken   string
	GitLabBranch  string
	ProjectIDs    []int
	ProjectsLimit int
	ExportPath    string
	FilesMask     string
}

// DefaultConfig returns config with default values applied.
func DefaultConfig() Config {
	return Config{
		ProjectsLimit: 100,
		FilesMask:     "*",
		GitLabBranch:  "main",
	}
}
