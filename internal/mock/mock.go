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

// Package mock provides test doubles for domain interfaces.
package mock

import (
	"sync"

	"gitlabFileScanner/internal/domain"
)

// Scanner is a mock implementation of domain.Scanner.
type Scanner struct {
	mu             sync.RWMutex
	Projects       []domain.Project
	Files          []string
	ProjectsErr    error
	FilesErr       error
	GetProjectsCalled bool
	GetFilesCalled   bool
}

func (m *Scanner) GetProjects(limit int, ids ...int) ([]domain.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetProjectsCalled = true
	return m.Projects, m.ProjectsErr
}

func (m *Scanner) GetFilePaths(projectID int64, branch string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.GetFilesCalled = true
	return m.Files, m.FilesErr
}

// Filter is a mock implementation of domain.Filter.
type Filter struct {
	mu          sync.RWMutex
	Result      []string
	ApplyCalled bool
	LastMask    string
	LastPaths   []string
}

func (m *Filter) Apply(filePaths []string, mask string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ApplyCalled = true
	m.LastMask = mask
	m.LastPaths = filePaths
	if m.Result != nil {
		return m.Result
	}
	return filePaths
}

// Exporter is a mock implementation of domain.Exporter.
type Exporter struct {
	mu         sync.RWMutex
	SavedPath  string
	SaveErr    error
	SaveData   *domain.FileList
	SaveCalled bool
}

func (m *Exporter) Save(path string, data *domain.FileList) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SaveCalled = true
	m.SaveData = data
	return m.SavedPath, m.SaveErr
}

// Logger is a mock implementation of domain.Logger.
type Logger struct {
	mu           sync.RWMutex
	DebugCalls   []string
	InfoCalls    []string
	WarnCalls    []string
	ErrorCalls   []string
}

func (l *Logger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.DebugCalls = append(l.DebugCalls, msg)
}

func (l *Logger) Info(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.InfoCalls = append(l.InfoCalls, msg)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.WarnCalls = append(l.WarnCalls, msg)
}

func (l *Logger) Error(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ErrorCalls = append(l.ErrorCalls, msg)
}
