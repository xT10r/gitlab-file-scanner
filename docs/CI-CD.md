# CI/CD

Проект собирает и тестирует GitHub Actions. Все workflow - в `.github/workflows/`.

## Файлы

| Файл | За что отвечает | Когда запускается |
|------|----------------|-------------------|
| `ci.yml` | Линтер, тесты, security, сборка | push, PR |
| `release.yml` | GoReleaser → GitHub Releases | тег `v*` |
| `docker.yml` | Docker build + push | push, тег |

## Порядок выполнения в CI

```
lint → test → security → build
```

Если линтер упал - дальше не пойдёт. Если тесты не прошли - сборка не запустится.

### Линтер

Запускает `go mod tidy`, `go vet` и `golangci-lint`. Если `go mod tidy` что-то меняет - CI упадёт с подсказкой запустить его локально.

### Тесты

```bash
go test -race -coverprofile=coverage.out -v ./test/...
```

Флаг `-race` ищет состояния гонки. Coverage сохраняется как артефакт на 30 дней.

### Security

`govulncheck` сканирует зависимости на известные уязвимости.

### Сборка

Кросс-компиляция для всех платформ:

| OS | amd64 | arm64 |
|----|:-----:|:-----:|
| Linux | ✅ | ✅ |
| macOS | ✅ | ✅ |
| Windows | ✅ | ✅ |

Каждый бинарник - статический (`CGO_ENABLED=0`), без отладочных символов (`-ldflags="-s -w"`).

## Релиз

При создании тега `v1.2.3` запускается GoReleaser:

```bash
git tag v1.2.3 && git push origin v1.2.3
```

Он:
1. Собирает бинарники для всех платформ
2. Создаёт архивы (tar.gz для Linux/macOS, zip для Windows)
3. Генерирует `checksums.txt` (SHA256)
4. Создаёт GitHub Release с changelog из коммитов

### Changelog в релизе

Автоматически из Conventional Commits:
- `feat(...)` → Features
- `fix(...)` → Bug Fixes
- `refactor(...)` → Refactoring

Коммиты `docs:`, `chore:`, `ci:` не попадают в changelog.

## Docker

### Теги образов

| Событие | Теги |
|---------|------|
| PR | `pr-<номер>`, `<sha>` |
| Push в ветку | `<ветка>`, `<sha>` |
| Тег `v1.2.3` | `1.2.3`, `1.2`, `latest` |

### Платформы

- `linux/amd64`
- `linux/arm64`

### Секреты

Нужно настроить в GitHub Settings → Secrets and variables → Actions:

| Секрет | Что туда |
|--------|----------|
| `DOCKER_REGISTRY_USERNAME` | Логин Docker Hub |
| `DOCKER_REGISTRY_PASSWORD` | Токен/пароль |
| `DOCKER_REGISTRY_PATH` | Namespace, например `xt10r` |

## Ручной запуск

Всё можно запустить вручную из вкладки **Actions** на GitHub.
