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

// Package filefilter implements domain.Filter using glob-style masks.
package filefilter

import (
	"path/filepath"
	"sort"

	"gitlabFileScanner/internal/strutil"
)

// Filter implements domain.Filter.
type Filter struct{}

// New creates a new file filter.
func New() *Filter {
	return &Filter{}
}

// Apply filters file paths by the given mask.
func (f *Filter) Apply(filePaths []string, mask string) []string {
	if len(filePaths) == 0 {
		return nil
	}

	maskParts := strutil.SplitMask(mask)
	var filtered []string

	for _, part := range maskParts {
		re, err := strutil.MaskToFileRegex(part)
		if err != nil {
			continue
		}
		for _, p := range filePaths {
			if re.MatchString(filepath.Base(p)) {
				filtered = append(filtered, p)
			}
		}
	}

	sort.Strings(filtered)
	return filtered
}
