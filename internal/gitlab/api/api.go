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
	"fmt"
	"sync"

	"gitlabFileScanner/internal/text"

	"github.com/xanzy/go-gitlab"
)

type gitlabAPI struct {
	client *gitlab.Client
}

func NewGitlabAPI(url, token string) (*gitlabAPI, error) {

	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании клиента GitLab: %v", err)
	}

	return &gitlabAPI{
		client: client,
	}, nil

}

func (api *gitlabAPI) GetProjects(projectsLimitCount int, projectIds ...int) ([]*gitlab.Project, error) {

	fmt.Printf("Получение проектов...\n")

	listOpts := gitlab.ListOptions{
		PerPage: 100, // ограничение по количеству проектов на странице
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

	if len(projectIds) > 0 && projectIds[0] > 0 {
		fmt.Printf("Используется фильтр по ID проектов")
		for _, projectId := range projectIds {
			project, _, err := api.client.Projects.GetProject(projectId, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("ошибка при получении проекта с ID %d: %v", projectId, err)
			}
			projects = append(projects, project)
		}
		return projects, nil
	}

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
		if receivedProjectsCount >= projectsLimitCount || len(pageProjects) < listOpts.PerPage {
			break
		}

		// Устанавливаем следующую страницу для запроса
		listOpts.Page = resp.NextPage
		listProjectOpts.ListOptions = listOpts
	}

	fmt.Printf("Проекты получены в количестве %d шт.\n\n", len(projects))
	return projects, nil
}

func (api *gitlabAPI) GetRepositoryFilePaths(projectId int, ref string) []string {

	var wg sync.WaitGroup
	var files []string

	// Буферизованный канал для передачи результатов сканирования
	resultCh := make(chan string, 100)

	// Переменная для подсчета активных горутин
	var activeGoroutines sync.WaitGroup

	// Получаем список веток репозитория
	branches, _, err := api.client.Branches.ListBranches(projectId, nil)
	if err != nil {
		fmt.Printf("Ошибка получения списка веток: %v\n", err)
		return files
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
		fmt.Printf("Указанная ветка '%s' отсутствует в репозитории\n", ref)
		return files
	}

	listTreeOpts := &gitlab.ListTreeOptions{
		Ref:  text.StringPtr(ref),
		Path: text.StringPtr("/"),
	}

	// Получаем список корневых каталогов репозитория
	tree, _, err := api.client.Repositories.ListTree(projectId, listTreeOpts)
	if err != nil {
		fmt.Printf("Ошибка получения списка файлов: %v\n", err)
		return files
	}

	// Сканируем каждый корневой каталог в отдельной горутине
	for _, dir := range tree {
		wg.Add(1)
		activeGoroutines.Add(1)
		go func(dir *gitlab.TreeNode) {
			defer activeGoroutines.Done()
			scanDir(api.client, projectId, ref, dir, resultCh, &wg)
		}(dir)
	}

	// Запускаем горутину для ожидания завершения всех задач сканирования
	go func() {
		wg.Wait()
		close(resultCh)
		activeGoroutines.Wait()
	}()

	// Собираем пути файлов из канала
	for path := range resultCh {
		files = append(files, path)
	}

	return files
}

func scanDir(client *gitlab.Client, projectId int, ref string, dir *gitlab.TreeNode, resultCh chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var mu sync.Mutex

	// Для случая, когда текущий элемент дерева является файлом
	if dir.Type == "blob" {
		mu.Lock()
		resultCh <- dir.Path
		mu.Unlock()
		return
	}

	// Получаем содержимое каталога
	tree, _, err := client.Repositories.ListTree(projectId, &gitlab.ListTreeOptions{
		Ref:  text.StringPtr(ref),
		Path: &dir.Path,
	})
	if err != nil {
		fmt.Println("Ошибка при получении содержимого каталога:", err)
		return
	}
	// Перебираем файлы в каталоге и добавляем их путь в канал resultCh
	for _, file := range tree {
		if file.Type == "blob" {
			mu.Lock()
			resultCh <- file.Path // Отправляем результат сканирования
			mu.Unlock()
		} else if file.Type == "tree" {
			wg.Add(1)
			go scanDir(client, projectId, ref, file, resultCh, wg) // Отправляем задачу сканирования дочернего каталога
		}
	}
}
