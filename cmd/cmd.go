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

package cmd

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

	api, err := api.NewGitlabAPI(fs.GetValueStr(flags.GitlabURLFlag), fs.GetValueStr(flags.GitlabApiTokenFlag))
	if err != nil {
		return fmt.Errorf("ошибка при создании клиента GitLab: %v", err)

	}

	// Получаем все доступные проекты
	projects, _ := api.GetProjects(fs.GetValueInt(flags.GitlabProjectIdFlag))
	fmt.Printf("Получено проектов: %d\n\n", len(projects))

	// Получаем файлы из проектов
	for _, project := range projects {
		fmt.Printf("Проект: %s (%d)\n", project.Name, project.ID)
		files := api.GetRepositoryFilePaths(project.ID, fs.GetValueStr(flags.GitlabBranchFlag))
		if len(files) == 0 {
			fmt.Println("---")
			continue
		}

		filteredFiles := file.FilterFilesByMask(files, fs.GetValueStr(flags.ExportFilesMaskFlag))

		if len(filteredFiles) > 0 {
			err := file.SaveFilesListToJSON(fs.GetValueStr(flags.ExportFilesPathFlag), filteredFiles, project.Name, project.ID, fs.GetValueStr(flags.GitlabBranchFlag))
			if err != nil {
				fmt.Printf("Ошибка при сохранении списка файлов: %v\n", err)
			}
		}
		fmt.Println("---")

	}

	fmt.Printf("Затрачено времени: %s\n", text.GetDurationString(time.Since(startTime)))

	return nil

}
