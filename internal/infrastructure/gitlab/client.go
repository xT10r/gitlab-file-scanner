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

// Package gitlab implements the domain.Scanner interface using the GitLab API.
package gitlab

import (
	"context"
	"fmt"
	"sync"

	"gitlabFileScanner/internal/domain"
	"gitlabFileScanner/internal/strutil"

	"gitlab.com/gitlab-org/api/client-go"
)

// Client implements domain.Scanner.
type Client struct {
	gl  *gitlab.Client
	ctx context.Context
}

// NewClient creates a new GitLab client.
func NewClient(ctx context.Context, url, token string) (*Client, error) {
	gl, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, fmt.Errorf("creating GitLab client: %w", err)
	}
	return &Client{gl: gl, ctx: ctx}, nil
}

// GetProjects returns a list of projects, optionally filtered by IDs.
func (c *Client) GetProjects(limit int, ids ...int) ([]domain.Project, error) {
	if len(ids) > 0 && ids[0] > 0 {
		return c.getProjectsByIDs(ids)
	}
	return c.listProjects(limit)
}

// GetFilePaths returns all file paths in a project for the given branch.
func (c *Client) GetFilePaths(projectID int64, branch string) ([]string, error) {
	return c.scanRepository(int(projectID), branch)
}

func (c *Client) getProjectsByIDs(ids []int) ([]domain.Project, error) {
	var projects []domain.Project
	for _, id := range ids {
		p, _, err := c.gl.Projects.GetProject(id, nil)
		if err != nil {
			return nil, fmt.Errorf("getting project %d: %w", id, err)
		}
		projects = append(projects, toDomainProject(p))
	}
	return projects, nil
}

func (c *Client) listProjects(limit int) ([]domain.Project, error) {
	var perPage int64 = 100
	if limit < 100 {
		perPage = int64(limit)
	}

	opts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: perPage,
			Page:    1,
		},
		IncludePendingDelete: strutil.BoolPtr(false),
		IncludeHidden:        strutil.BoolPtr(false),
		Archived:             strutil.BoolPtr(false),
	}

	var projects []domain.Project
	var received int

	for {
		page, resp, err := c.gl.Projects.ListProjects(opts)
		if err != nil {
			return nil, fmt.Errorf("listing projects: %w", err)
		}

		for _, p := range page {
			projects = append(projects, toDomainProject(p))
		}
		received += len(page)

		if received >= limit || len(page) < int(perPage) {
			break
		}

		opts.Page = resp.NextPage
	}

	return projects, nil
}

func (c *Client) scanRepository(projectID int, ref string) ([]string, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	default:
	}

	branches, _, err := c.gl.Branches.ListBranches(projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("listing branches: %w", err)
	}

	var refExists bool
	for _, b := range branches {
		if b.Name == ref {
			refExists = true
			break
		}
	}
	if !refExists {
		return nil, fmt.Errorf("branch '%s' not found in repository", ref)
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		files     []string
		semaphore = make(chan struct{}, 50)
	)

	var scanDir func(string)
	scanDir = func(path string) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()
		defer wg.Done()

		page := int64(1)
		for {
			tree, resp, err := c.gl.Repositories.ListTree(projectID, &gitlab.ListTreeOptions{
				Ref:  strutil.StringPtr(ref),
				Path: strutil.StringPtr(path),
				ListOptions: gitlab.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			})
			if err != nil {
				return
			}
			if resp == nil {
				return
			}

			for _, item := range tree {
				switch item.Type {
				case "blob":
					mu.Lock()
					files = append(files, item.Path)
					mu.Unlock()
				case "tree":
					wg.Add(1)
					go func(p string) {
						scanDir(p)
					}(item.Path)
				}
			}

			if resp.CurrentPage >= resp.TotalPages {
				break
			}
			page++
		}
	}

	wg.Add(1)
	go scanDir("/")
	wg.Wait()

	return files, nil
}

func toDomainProject(p *gitlab.Project) domain.Project {
	return domain.Project{
		ID:     p.ID,
		Name:   p.Name,
		WebURL: p.WebURL,
	}
}
