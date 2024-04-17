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

package text

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Возвращает указатель на строку
func StringPtr(s string) *string {
	return &s
}

// Возвращает указатель на булево
func BoolPtr(b bool) *bool {
	return &b
}

func MaskToFileRegex(mask string) (*regexp.Regexp, error) {

	// Удаляем пробелы из маски
	mask = strings.ReplaceAll(mask, " ", "")

	// Преобразуем маску в регулярное выражение
	regex := strings.ReplaceAll(mask, ".", `\.`)  // Экранируем точки
	regex = strings.ReplaceAll(regex, "*", ".*")  // Заменяем "*" на ".*"
	regex = strings.ReplaceAll(regex, "(", "(?:") // Преобразуем открывающиеся скобки в незахватывающие группы

	return regexp.Compile("^" + regex + "$")
}

func SplitMask(mask string) []string {
	var parts []string
	var part strings.Builder
	var level int

	for _, char := range mask {
		switch char {
		case '|':
			if level == 0 {
				parts = append(parts, part.String())
				part.Reset()
			} else {
				part.WriteRune(char)
			}
		case '(':
			level++
			part.WriteRune(char)
		case ')':
			level--
			part.WriteRune(char)
		default:
			part.WriteRune(char)
		}
	}

	parts = append(parts, part.String())

	return parts
}

func GetDurationString(duration time.Duration) string {

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, pluralize(days, "день", "дня", "дней"))
	}
	if hours > 0 {
		parts = append(parts, pluralize(hours, "час", "часа", "часов"))
	}
	if minutes > 0 {
		parts = append(parts, pluralize(minutes, "минута", "минуты", "минут"))
	}
	if seconds > 0 {
		parts = append(parts, pluralize(seconds, "секунда", "секунды", "секунд"))
	}

	return fmt.Sprintf("%s", strings.Join(parts, " "))
}

func pluralize(count int, form1, form2, form5 string) string {
	if count == 0 {
		return ""
	}
	if count == 1 {
		return fmt.Sprintf("%d %s", count, form1)
	}
	if count >= 2 && count <= 4 {
		return fmt.Sprintf("%d %s", count, form2)
	}
	return fmt.Sprintf("%d %s", count, form5)
}
