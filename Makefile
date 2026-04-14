# Параметры по умолчанию
TAG ?= 1.0.0

# Пути к учетным данным DockerHub
DOCKER_USERNAME ?= $(shell echo $$DOCKER_REGISTRY_USERNAME)
DOCKER_PASSWORD ?= $(shell echo $$DOCKER_REGISTRY_PASSWORD)

# Docker Registry URL
DOCKER_REGISTRY_URL ?= docker.io

# Путь к образу в Docker Registry
DOCKER_REGISTRY_PATH ?= $(shell echo $$DOCKER_REGISTRY_PATH)

# Имя образа
IMAGE_NAME := gitlab-file-scanner

# Пути к Docker и к Docker Hub
DOCKER := docker
DOCKER_LOGIN := $(DOCKER) login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY_URL)

# Путь до каталога временных файлов
FILES_LISTS="./files_lists"

# Команда для сборки Docker образа
build:
	@$(DOCKER) build -t $(IMAGE_NAME):$(TAG) .

# Команда для запуска тестов
tests:
	$(DOCKER) run --rm $(IMAGE_NAME):$(TAG) -test

run-test:
	@mkdir -p $(FILES_LISTS)
	@$(DOCKER) run --rm -v $(FILES_LISTS):/app/files_lists -e "GITLAB_FILE_SCANNER_SERVER_URL=https://gitlab.com" -e "GITLAB_FILE_SCANNER_BRANCH=main" -e "GITLAB_FILE_SCANNER_PROJECT_ID=56903310" -e "GITLAB_FILE_SCANNER_EXPORT_PATH=$(FILES_LISTS)" $(IMAGE_NAME):$(TAG)
	@$(DOCKER) run --rm -v $(FILES_LISTS):/app/files_lists -e "GITLAB_FILE_SCANNER_SERVER_URL=https://gitlab.com" -e "GITLAB_FILE_SCANNER_BRANCH=main" -e "GITLAB_FILE_SCANNER_PROJECT_ID=56903311" -e "GITLAB_FILE_SCANNER_EXPORT_PATH=$(FILES_LISTS)" -e "GITLAB_FILE_SCANNER_FILEMASK=*.md" $(IMAGE_NAME):$(TAG)

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
	@go run github.com/google/addlicense@latest -check -f HEADER $$(find . -name "*.go" -not -path "./vendor/*")

# Применить заголовки лицензий (локально)
license-fix:
	@go run github.com/google/addlicense@latest -f HEADER $$(find . -name "*.go" -not -path "./vendor/*")