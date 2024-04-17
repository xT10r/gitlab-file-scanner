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

package app

import (
	"fmt"
	"time"

	"gitlabFileScanner/internal/file"
	"gitlabFileScanner/internal/flags"
	"gitlabFileScanner/internal/gitlab/api"
	"gitlabFileScanner/internal/text"
)

func Start() error {

	// Получение флагов и их парсинг
	fs, err := flags.NewFlagSet()
	if err != nil {
		return err
	}

	startTime := time.Now()

	api, err := api.NewGitlabAPI(fs.GetValue(flags.GitlabURLFlag), fs.GetValue(flags.GitlabApiTokenFlag))
	if err != nil {
		return fmt.Errorf("ошибка при создании клиента GitLab: %v", err)

	}

	// Получаем все доступные проекты
	fmt.Printf("\n[Работа с проектами Gitlab]\n")
	projects, _ := api.GetProjects(fs.GetValueInt(flags.GitlabProjectsLimitFlag), fs.GetValueInt(flags.GitlabProjectIdFlag))

	// Получаем пути файлов из проектов с учетом маски
	projectsTotal := len(projects)
	for projectIndex, project := range projects {
		projectNumberStr := fmt.Sprintf("%0*d", len(fmt.Sprintf("%d", projectsTotal)), projectIndex+1)

		// Выводим информацию о проекте
		fmt.Printf("%s/%d | %d | %s", projectNumberStr, projectsTotal, project.ID, project.Name)

		files, err := api.GetRepositoryFilePaths(project.ID, fs.GetValue(flags.GitlabBranchFlag))
		if err != nil {
			// Выводим информацию об ошибке
			fmt.Printf(" | %v\n", err)
			continue
		}

		filteredFilePaths := file.FilterFilesByMask(files, fs.GetValue(flags.FilesMaskFlag))
		if len(filteredFilePaths) == 0 {
			// Выводим информацию об отсутствии файлов после отбора по маске
			fmt.Printf(" | %d/%d | не найдено файлов соответствующих заданной маске\n", len(files), len(filteredFilePaths))
			continue
		}

		// Создаем структуру для сохранения данных
		fileData := &file.GitlabFilePathsStruct{
			Name:      project.Name,
			WebURL:    project.WebURL,
			ID:        project.ID,
			Branch:    fs.GetValue(flags.GitlabBranchFlag),
			FilePaths: filteredFilePaths,
		}

		filePath, err := file.SaveFilesListToJSON(fs.GetValue(flags.ExportFilesPathFlag), fileData)
		if err != nil {
			fmt.Printf("Ошибка при сохранении списка файлов: %v\n", err)
		}

		fmt.Printf(" | %d/%d | %s\n", len(files), len(filteredFilePaths), filePath)
	}

	fmt.Printf("Затрачено времени: %s\n", text.GetDurationString(time.Since(startTime)))

	return nil

}
