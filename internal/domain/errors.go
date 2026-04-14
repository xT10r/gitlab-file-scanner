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

import "errors"

// Sentinel errors for the application.
var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrBranchNotFound    = errors.New("branch not found in repository")
	ErrNoFilesFound      = errors.New("no files found in repository")
	ErrNoMatchingFiles   = errors.New("no files matching the specified mask")
	ErrInvalidFileMask   = errors.New("invalid file mask")
	ErrExportFailed      = errors.New("failed to export file list")
)
