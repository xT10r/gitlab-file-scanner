# Параметры по умолчанию
TAG ?= 1.0.0
VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
GIT_BRANCH ?= $(shell git symbolic-ref --short HEAD 2>/dev/null || echo "detached")
BUILD_DATE ?= $(shell date -u +%Y-%m-%d)

# Go linker flags for version info
LDFLAGS = -s -w \
    -X gitlabFileScanner/internal/cli/commands.version=$(VERSION) \
    -X gitlabFileScanner/internal/cli/commands.commit=$(GIT_COMMIT) \
    -X gitlabFileScanner/internal/cli/commands.branch=$(GIT_BRANCH) \
    -X gitlabFileScanner/internal/cli/commands.buildDate=$(BUILD_DATE)

# Пути к учетным данным DockerHub
DOCKER_USERNAME ?= $(shell echo $$DOCKER_REGISTRY_USERNAME)
DOCKER_PASSWORD ?= $(shell echo $$DOCKER_REGISTRY_PASSWORD)

# Docker Registry URL
DOCKER_REGISTRY_URL ?= docker.io

# Путь к образу в Docker Registry
DOCKER_REGISTRY_PATH ?= $(shell echo $$DOCKER_REGISTRY_PATH)

# Имя образа
IMAGE_NAME := gitlab-file-scanner
RACE_TEST_IMAGE := $(IMAGE_NAME):race-test

# Пути к Docker и к Docker Hub
DOCKER := docker
DOCKER_LOGIN := $(DOCKER) login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY_URL)

# Путь до каталога временных файлов
FILES_LISTS="./files_lists"

GOOS ?= $(shell go env GOOS 2>/dev/null || echo "unknown")
BINARY_NAME := gitlab-file-scanner$(if $(filter windows,$(GOOS)),.exe)

# Команда для сборки Docker образа
build:
	@$(DOCKER) build -t $(IMAGE_NAME):$(TAG) .

# Сборка Go бинарника (локальная)
cli:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd

# Команда для запуска тестов
tests:
	go test ./... -v

test-race-docker:
	@$(DOCKER) build --target race-test -t $(RACE_TEST_IMAGE) .

run-test:
	@mkdir -p $(FILES_LISTS)
	@$(DOCKER) run --rm -v $(FILES_LISTS):/app/files_lists -e "GFS_URL=https://gitlab.com" -e "GFS_BRANCH=main" -e "GFS_PROJECT_ID=56903310" -e "GFS_EXPORT_PATH=/app/files_lists" $(IMAGE_NAME):$(TAG) scan
	@$(DOCKER) run --rm -v $(FILES_LISTS):/app/files_lists -e "GFS_URL=https://gitlab.com" -e "GFS_BRANCH=main" -e "GFS_PROJECT_ID=56903311" -e "GFS_EXPORT_PATH=/app/files_lists" -e "GFS_FILEMASK=*.md" $(IMAGE_NAME):$(TAG) scan

# Команда для публикации Docker образа в Docker Hub
push:
	@$(DOCKER_LOGIN) 2>/dev/null
	$(DOCKER) tag $(IMAGE_NAME):$(TAG) $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):$(TAG)
	$(DOCKER) push $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):$(TAG)

# Пометить образ как latest и отправить его в Docker Hub
push-latest:
	@$(DOCKER_LOGIN) 2>/dev/null
	$(DOCKER) tag $(IMAGE_NAME):$(TAG) $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):latest
	$(DOCKER) push $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):latest

# Удалить Docker образ
clean:
	@rm -rf $(FILES_LISTS) 2>/dev/null
	@$(DOCKER) rmi -f $(IMAGE_NAME):$(TAG) 2>/dev/null

# Проверить заголовки лицензий (для CI)
license-check:
	@echo "Checking licenses..."
	@go run github.com/google/addlicense@v1.2.0 -c "Alex Dobshikov" -check -f HEADER $$(find . -name "*.go" -not -path "./vendor/*")

# Применить заголовки лицензий (локально)
license-fix:
	@echo "Updating licenses from HEADER template..."
	@go run github.com/google/addlicense@v1.2.0 -c "Alex Dobshikov" -f HEADER $$(find . -name "*.go" -not -path "./vendor/*")
