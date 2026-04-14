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

package domain

// Scanner defines the interface for interacting with GitLab repositories.
type Scanner interface {
	// GetProjects returns a list of projects, optionally filtered by IDs.
	GetProjects(limit int, ids ...int) ([]Project, error)

	// GetFilePaths returns all file paths in a project for the given branch.
	GetFilePaths(projectID int64, branch string) ([]string, error)
}
