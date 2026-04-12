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

// Package filesystem implements domain.Exporter for writing to the local filesystem.
package filesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gitlabFileScanner/internal/domain"
)

// Writer implements domain.Exporter by writing JSON files to disk.
type Writer struct{}

// NewWriter creates a new filesystem writer.
func NewWriter() *Writer {
	return &Writer{}
}

// Save writes the file list as JSON to the given directory.
func (w *Writer) Save(dir string, data *domain.FileList) (string, error) {
	jsonData, err := json.MarshalIndent(toExportStruct(data), "", "    ")
	if err != nil {
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}

	filename := fmt.Sprintf("%d.json", data.ID)
	fullPath := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("creating file: %w", err)
	}
	defer func() { _ = file.Close() }()

	if _, err := file.Write(jsonData); err != nil {
		return "", fmt.Errorf("writing file: %w", err)
	}

	return fullPath, nil
}

// exportStruct is the JSON representation of a file list.
type exportStruct struct {
	Name      string   `json:"name"`
	WebURL    string   `json:"web_url"`
	ID        int64    `json:"id"`
	Branch    string   `json:"branch"`
	FilePaths []string `json:"files"`
}

func toExportStruct(data *domain.FileList) exportStruct {
	return exportStruct{
		Name:      data.Name,
		WebURL:    data.WebURL,
		ID:        data.ID,
		Branch:    data.Branch,
		FilePaths: data.FilePaths,
	}
}
