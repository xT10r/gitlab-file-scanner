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

package api

import (
	"context"
	"fmt"
	"sync"

	"gitlabFileScanner/internal/text"

	"gitlab.com/gitlab-org/api/client-go"
)

type gitlabAPI struct {
	client *gitlab.Client
	ctx    context.Context
}

func NewGitlabAPI(ctx context.Context, url, token string) (*gitlabAPI, error) {

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании клиента GitLab: %v", err)
	}

	return &gitlabAPI{
		client: client,
		ctx:    ctx,
	}, nil

}

func (api *gitlabAPI) GetProjects(projectsLimitCount int, projectIds ...int) ([]*gitlab.Project, error) {

	useProjectIdsFilter := len(projectIds) > 0 && projectIds[0] > 0
	message := "Получение проектов"
	if useProjectIdsFilter {
		message = fmt.Sprintf("%s (с использованием отбора по идентификаторам)", message)
	}

	fmt.Printf("%s...\n", message)

	var perPage int64
	if projectsLimitCount >= 100 {
		perPage = 100
	} else {
		perPage = int64(projectsLimitCount)
	}

	listOpts := gitlab.ListOptions{
		PerPage: perPage, // ограничение по количеству проектов на странице
		Page:    1,
		Sort:    "asc",
	}

	listProjectOpts := &gitlab.ListProjectsOptions{
		ListOptions:          listOpts,
		IncludePendingDelete: text.BoolPtr(false),
		IncludeHidden:        text.BoolPtr(false),
		Archived:             text.BoolPtr(false),
	}

	var projects []*gitlab.Project

	if useProjectIdsFilter {
		for _, projectId := range projectIds {
			project, _, err := api.client.Projects.GetProject(projectId, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("ошибка при получении проекта с ID %d: %v", projectId, err)
			}
			projects = append(projects, project)
		}

	} else {

		// Переменная для отслеживания количества полученных проектов
		var receivedProjectsCount int

		// Цикл для получения проектов постранично
		for {
			// Получаем список проектов на текущей странице
			pageProjects, resp, err := api.client.Projects.ListProjects(listProjectOpts)
			if err != nil {
				return nil, fmt.Errorf("ошибка при получении проектов: %v", err)
			}

			// Добавляем проекты из текущей страницы к общему списку
			projects = append(projects, pageProjects...)
			receivedProjectsCount += len(pageProjects)

			// Если получили нужное количество проектов или достигли последней страницы
			if receivedProjectsCount >= projectsLimitCount || int64(len(pageProjects)) < listOpts.PerPage {
				break
			}

			// Устанавливаем следующую страницу для запроса
			listOpts.Page = resp.NextPage
			listProjectOpts.ListOptions = listOpts
		}

	}

	fmt.Printf("Проекты получены в количестве %d шт.\n\n", len(projects))

	return projects, nil
}

// GetRepositoryFilePaths получает все пути файлов в репозитории для заданной ветки.
func (api *gitlabAPI) GetRepositoryFilePaths(projectId int, ref string) ([]string, error) {

	// Проверяем контекст перед началом работы
	select {
	case <-api.ctx.Done():
		return nil, api.ctx.Err()
	default:
	}

	var wg sync.WaitGroup
	var files []string
	var mu sync.Mutex // Общий мьютекс для защиты доступа к общему ресурсу (files).

	// Получаем список веток репозитория
	branches, _, err := api.client.Branches.ListBranches(projectId, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка веток: %v", err)
	}

	// Проверяем наличие указанной ветки в списке веток репозитория
	var refExists bool
	for _, branch := range branches {
		if branch.Name == ref {
			refExists = true
			break
		}
	}
	if !refExists {
		return nil, fmt.Errorf("ветка с именем '%s' отсутствует в репозитории", ref)
	}

	// Семафор для ограничения количества горутин
	semaphore := make(chan struct{}, 50)

	// Функция для сканирования каталога
	var scanDir func(string)

	scanDir = func(path string) {

		defer wg.Done()
		defer func() { <-semaphore }() // Освобождаем слот

		semaphore <- struct{}{} // Захватываем слот

		// Получаем содержимое каталога постранично
		var page int64 = 1
		for {

			listOpts := gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
				Sort:    "asc",
			}

			listTreeOpts := &gitlab.ListTreeOptions{
				Ref:         text.StringPtr(ref),
				Path:        text.StringPtr(path),
				ListOptions: listOpts,
			}

			tree, resp, err := api.client.Repositories.ListTree(projectId, listTreeOpts)
			if err != nil {
				fmt.Printf("Ошибка при получении содержимого каталога '%s': %v\n", path, err)
				return
			}

			// Проверяем, что resp не nil
			if resp == nil {
				fmt.Printf("Пустой ответ для каталога '%s'\n", path)
				return
			}

			// Перебираем элементы в каталоге
			for _, item := range tree {
				switch item.Type {
				case "blob": // Файл
					mu.Lock()
					files = append(files, item.Path)
					mu.Unlock()
				case "tree": // Каталог
					wg.Add(1)
					go func(itemPath string) {
						scanDir(itemPath)
					}(item.Path)
				}
			}

			// Проверяем, есть ли еще страницы
			if resp.CurrentPage >= resp.TotalPages {
				break
			}
			page++
		}
	}

	// Сканируем корневой каталог
	wg.Add(1)
	go scanDir("/")

	wg.Wait() // Ожидаем завершения всех горутин
	return files, nil
}
