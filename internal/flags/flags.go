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

package flags

import (
	"flag"
	"fmt"
	"gitlabFileScanner/internal/file"
	"net/url"
	"os"
	"strconv"
)

// Enum для идентификации каждого флага
type Flag string

type FlagSet struct {
	flags *flag.FlagSet
}

// Enum для флагов флагов
const (
	GitlabApiTokenFlag  Flag = "token"
	GitlabURLFlag       Flag = "url"
	GitlabBranchFlag    Flag = "branch"
	GitlabProjectIdFlag Flag = "project-id"
	ExportFilesPathFlag Flag = "export-files-path"
	ExportFilesMaskFlag Flag = "export-files-mask"

	// В случае появления нового флага, добавляем для него запись тут
)

// Создает новый набор флагов
func NewFlagSet() (*FlagSet, error) {
	fs := flag.NewFlagSet("myflags", flag.ErrorHandling(flag.ExitOnError))

	fs.String(string(GitlabApiTokenFlag), os.Getenv("GITLAB_FILE_SCANNER_API_TOKEN"), "Токен GitLab")
	fs.String(string(GitlabURLFlag), os.Getenv("GITLAB_FILE_SCANNER_SERVER_URL"), "URL-адрес сервера GitLab")
	fs.String(string(GitlabBranchFlag), os.Getenv("GITLAB_FILE_SCANNER_BRANCH"), "Ветка для получения списка файлов")

	projectId, _ := stringToNumber(os.Getenv("GITLAB_FILE_SCANNER_PROJECT_ID"))
	fs.Int(string(GitlabProjectIdFlag), projectId, "Идентификатор проекта")

	fs.String(string(ExportFilesPathFlag), os.Getenv("GITLAB_FILE_SCANNER_EXPORT_PATH"), "Путь для выгрузки списка файлов")
	fs.String(string(ExportFilesMaskFlag), os.Getenv("GITLAB_FILE_SCANNER_FILEMASK"), "Маска файлов")

	// В случае появления нового флага, добавляем для него запись тут

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге флагов: %v", err)
	}

	cfs := &FlagSet{
		flags: fs,
	}

	// Валидация флагов
	err := cfs.checkFlags()
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке значений флагов: %v", err)
	}

	return cfs, nil
}

// Возвращает целочисленное значение указанного флага
func (fs *FlagSet) GetValueInt(flag Flag) int {
	value := fs.GetValueStr(flag)
	num, _ := stringToNumber(value)
	return num
}

// Возвращает строковое значение указанного флага
func (fs *FlagSet) GetValueStr(flag Flag) string {
	return fs.flags.Lookup(string(flag)).Value.String()
}

// Записывает строковое значение указанного флага
func (fs *FlagSet) SetValueStr(flag Flag, value string) error {
	err := fs.flags.Lookup(string(flag)).Value.Set(value)
	if err != nil {
		return fmt.Errorf("не удалось установить значение флага: %s", string(flag))
	}
	return nil
}

// Валидация значений флагов
func (cfs *FlagSet) checkFlags() error {

	if err := cfs.checkGitLabToken(GitlabApiTokenFlag); err != nil {
		return err
	}
	if err := cfs.checkGitLabURL(GitlabURLFlag); err != nil {
		return err
	}
	if err := cfs.checkBranch(GitlabBranchFlag); err != nil {
		return err
	}
	if err := cfs.checkExportPath(ExportFilesPathFlag); err != nil {
		return err
	}
	if err := cfs.checkFileMask(ExportFilesMaskFlag); err != nil {
		return err
	}

	// В случае появления нового флага, добавляем для него проверку тут

	return nil
}

func (cfs *FlagSet) checkGitLabToken(flag Flag) error {
	token := cfs.GetValueStr(flag)
	if token == "" {
		return fmt.Errorf("не указан токен GitLab")
	}
	return nil
}

func (cfs *FlagSet) checkGitLabURL(flag Flag) error {
	urlStr := cfs.GetValueStr(flag)
	if urlStr == "" {
		return fmt.Errorf("не указан URL-адрес сервера GitLab")
	}

	_, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("некорректный URL-адрес сервера GitLab: %v", err)
	}

	return nil
}

func (cfs *FlagSet) checkBranch(flag Flag) error {
	branch := cfs.GetValueStr(flag)
	if branch == "" {
		return fmt.Errorf("не указана ветка по-умолчанию")
	}
	return nil
}

func (cfs *FlagSet) checkExportPath(flag Flag) error {
	path := cfs.GetValueStr(flag)
	if path == "" {
		return fmt.Errorf("не указан путь для экспорта списка файлов")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("ошибка при создании директории для экспорта: %v", err)
		}
	}
	return nil
}

func (cfs *FlagSet) checkFileMask(flag Flag) error {

	if cfs.GetValueStr(flag) == "" {
		cfs.SetValueStr(flag, "*")
	}

	mask := cfs.GetValueStr(flag)
	if mask == "" {
		return fmt.Errorf("маска файла не может быть пустой")
	}

	// Преобразование маски файла в регулярное выражение
	_, err := file.MaskToFileRegex(mask)
	if err != nil {
		return fmt.Errorf("неверный формат маски файлов: %v", err)
	}
	return nil
}

func stringToNumber(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("невозможно преобразовать строку в число: %v", err)
	}
	return num, nil
}
