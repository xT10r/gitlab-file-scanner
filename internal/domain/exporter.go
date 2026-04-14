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

// Exporter defines the interface for saving file lists to storage.
type Exporter interface {
	// Save writes the file list to the given path, returning the full file path.
	Save(path string, data *FileList) (string, error)
}

// Filter defines the interface for filtering file paths.
type Filter interface {
	// Apply filters file paths by the given mask, returning matched paths.
	Apply(filePaths []string, mask string) []string
}
