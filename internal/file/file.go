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

package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func SaveFilesListToJSON(exportPath string, filePaths []string, projectName string, projectId int, ref string) error {
	// Создаем структуру для сохранения данных
	data := struct {
		Name      string   `json:"Name"`
		ID        int      `json:"Id"`
		Branch    string   `json:"Branch"`
		FilePaths []string `json:"filePaths"`
	}{
		Name:      projectName,
		ID:        projectId,
		Branch:    ref,
		FilePaths: filePaths,
	}

	// Преобразуем структуру в JSON
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	// Генерируем имя файла
	filename := fmt.Sprintf("%d.json", projectId)

	// Создаем директорию для сохранения файла, если она не существует
	err = os.MkdirAll(exportPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Полный путь к файлу
	filePath := filepath.Join(exportPath, filename)

	// Создаем файл для записи данных
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем данные в файл
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Printf("Список файлов проекта: %s\n", filePath)
	return nil
}

// Функция для фильтрации списка путей файлов по маске
func FilterFilesByMask(filePaths []string, mask string) []string {
	var filteredPaths []string

	// Разделяем маску на отдельные выражения
	maskParts := strings.Split(mask, "|")

	// Создаем регулярное выражение для каждой части маски
	for _, maskPart := range maskParts {

		// Преобразуем маску в регулярное выражение и компилируем его
		r, err := MaskToFileRegex(maskPart)
		if err != nil {
			fmt.Println("Ошибка компиляции регулярного выражения:", err)
			continue
		}

		// Проверяем каждый путь файла на соответствие регулярному выражению
		for _, filePath := range filePaths {
			fileName := filepath.Base(filePath)
			if r.MatchString(fileName) {
				filteredPaths = append(filteredPaths, filePath)
			}
		}
	}

	return filteredPaths
}

func MaskToFileRegex(mask string) (*regexp.Regexp, error) {
	// Преобразуем маску файла в регулярное выражение
	regex := strings.ReplaceAll(mask, ".", `\.`) // Экранируем точки
	regex = strings.ReplaceAll(regex, "*", ".*") // Заменяем "*" на ".*"
	return regexp.Compile("^" + regex + "$")
}
