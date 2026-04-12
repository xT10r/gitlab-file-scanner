# Разработка

## Быстрый старт

```bash
git clone https://github.com/xT10r/gitlab-file-scanner.git
cd gitlab-file-scanner
make install-hooks    # git hook для проверки коммитов
go test ./test/... -v
go build -o gitlab-file-scanner ./cmd
```

## Коммиты

Формат - [Conventional Commits](https://conventionalcommits.org):

```
type(scope): описание
```

**Примеры:**

```
feat(api): add semaphore for concurrent file scanning
fix(text): handle pluralization for numbers 11-14
refactor(file): replace sort.StringSlice with sort.Strings
test(file): add FilterFilesByMask test cases
```

**Правила:**

- Английский, императив - "add" а не "added"
- Без точки в конце
- Scope обязательный

### Доступные типы

| Тип | Когда использовать |
|-----|-------------------|
| `feat` | Новая фича |
| `fix` | Баг-фикс |
| `refactor` | Рефакторинг, поведение не меняется |
| `docs` | Документация |
| `test` | Тесты |
| `chore` | Обслуживание, конфиги |
| `ci` | CI/CD, пайплайны |
| `build` | Зависимости, go.mod |

### Scope

| Scope | Что меняем |
|-------|-----------|
| `app` | `internal/app` |
| `flags` | `internal/flags` |
| `api` | `internal/gitlab/api` |
| `file` | `internal/file` |
| `text` | `internal/text` |
| `cmd` | `cmd/` |
| `docker` | Dockerfile |
| `ci` | Makefile, GitHub Actions |
| `deps` | go.mod, go.sum |

## Ветки

| Ветка | Зачем |
|-------|-------|
| `main` | Релизы |
| `develop` | Разработка |
| `feature/*` | Новые фичи |
| `fix/*` | Баг-фиксы |

## Тесты

```bash
# Запустить всё
go test ./test/... -v

# С покрытием
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

## Docker

```bash
make build       # собрать образ
make run-test    # запустить тестовый сценарий
make push        # опубликовать
```

## Полезные команды

```bash
make install-hooks     # установить git hook
make uninstall-hooks   # удалить git hook
make tests             # запустить тесты в Docker
```
