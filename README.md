# gitlab-file-scanner

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI](https://github.com/xT10r/gitlab-file-scanner/actions/workflows/ci.yml/badge.svg)](https://github.com/xT10r/gitlab-file-scanner/actions/workflows/ci.yml)
[![Docker](https://img.shields.io/badge/Docker-alpine-2496ED?logo=docker)](Dockerfile)

Утилита для сканирования репозиториев GitLab и выгрузки списков файлов в JSON.

## Что умеет

- Обходит все директории проекта и собирает пути файлов
- Работает параллельно - до 50 горутин с ограничением через семафор
- Фильтрует файлы по маске (`*.go`, `*.go|*.py`, `*.test.*`)
- Сохраняет результат в JSON - по одному файлу на проект
- Корректно останавливается по `Ctrl+C`
- Собирается в один статический бинарник

## Установка

### Скачать готовый бинарник

Страница [Releases](https://github.com/xT10r/gitlab-file-scanner/releases) - там бинарники для Linux, macOS и Windows.

### Собрать из исходников

```bash
git clone https://github.com/xT10r/gitlab-file-scanner.git
cd gitlab-file-scanner
go build -o gitlab-file-scanner ./cmd
```

### Docker

```bash
docker pull docker.io/xt10r/gitlab-file-scanner:latest
```

## Быстрый старт

```bash
gitlab-file-scanner scan \
  --url https://gitlab.com \
  --token glpat-xxx \
  --branch main \
  --files-mask "*.go" \
  --export-files-path ./output
```

Или через короткие переменные окружения:

```bash
export GFS_URL=https://gitlab.com
export GFS_API_TOKEN=glpat-xxx
export GFS_BRANCH=main
export GFS_FILEMASK="*.go"
export GFS_EXPORT_PATH=./output

gitlab-file-scanner scan
```

## Как пользоваться

Команда сканирования — `gitlab-file-scanner scan`.

### Через аргументы (параметры) командной строки

```bash
# Все проекты, только Markdown
gitlab-file-scanner scan \
  --url https://gitlab.com \
  --branch main \
  --files-mask "*.md" \
  --export-files-path ./output

# Один конкретный проект, все файлы
gitlab-file-scanner scan \
  --url https://gitlab.com \
  --branch main \
  --project-id 56903310 \
  --export-files-path ./output
```

### Через переменные окружения

```bash
export GFS_URL=https://gitlab.com
export GFS_BRANCH=main
export GFS_EXPORT_PATH=./output

gitlab-file-scanner scan
```

### Docker

```bash
docker run --rm \
  -v ./output:/app/files_lists \
  -e GFS_URL=https://gitlab.com \
  -e GFS_BRANCH=main \
  -e GFS_EXPORT_PATH=/app/files_lists \
  docker.io/xt10r/gitlab-file-scanner:latest scan
```

## Параметры

| Флаг | Env | По умолчанию | Описание |
|------|-----|:------------:|----------|
| `--url` | `GFS_URL` | - | URL GitLab |
| `--token` | `GFS_API_TOKEN` | - | Токен API |
| `--token-file` | `GFS_API_TOKEN_FILE` | - | Файл с токеном API |
| `--branch` | `GFS_BRANCH` | `main` | Ветка |
| `--project-id` | `GFS_PROJECT_ID` | - | ID проекта |
| `--projects-limit` | `GFS_PROJECTS_LIMIT` | 100 | Лимит проектов |
| `--export-files-path` | `GFS_EXPORT_PATH` | - | Куда сохранять JSON |
| `--files-mask` | `GFS_FILEMASK` | `*` | Маска файлов |

## Что на выходе

Один JSON-файл на проект:

```json
{
    "name": "my-project",
    "web_url": "https://gitlab.com/user/my-project",
    "id": 56903310,
    "branch": "main",
    "files": [
        "README.md",
        "src/main.go",
        "src/utils/helper.go"
    ]
}
```

## Маски файлов

| Маска | Что найдёт |
|-------|-----------|
| `*` | Всё |
| `*.go` | Только Go |
| `*.go\|*.mod` | Go и go.mod |
| `*.(go\|py)` | Файлы .go или .py |

## Подробнее

- [CI/CD](docs/CI-CD.md) - пайплайн, релизы, Docker
- [Development](docs/DEVELOPMENT.md) - коммиты, ветвление, тесты
- [Contributing](CONTRIBUTING.md) - правила контрибьюции

## Лицензия

[Apache License 2.0](LICENSE)
