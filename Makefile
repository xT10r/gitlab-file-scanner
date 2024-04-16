# Параметры по умолчанию
TAG ?= 1.0.0

# Пути к учетным данным DockerHub
DOCKER_USERNAME ?= $(shell echo $$DOCKER_REGISTRY_USERNAME)
DOCKER_PASSWORD ?= $(shell echo $$DOCKER_REGISTRY_PASSWORD)

# Docker Registry URL
DOCKER_REGISTRY_URL ?= docker.io

# Путь к образу в Docker Registry
DOCKER_REGISTRY_PATH ?= gitlab

# Имя образа
IMAGE_NAME := gitlab-file-scanner

# Пути к Docker и к Docker Hub
DOCKER := docker
DOCKER_LOGIN := $(DOCKER) login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY_URL)

# Команда для сборки Docker образа
build:
	@$(DOCKER) build -t $(IMAGE_NAME):$(TAG) .

# Команда для запуска тестов
test:
	$(DOCKER) run --rm $(IMAGE_NAME):$(TAG) -test

# Команда для публикации Docker образа в Docker Hub
push:
	$(DOCKER_LOGIN)
	$(DOCKER) tag $(IMAGE_NAME):$(TAG) $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):$(TAG)
	$(DOCKER) push $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):$(TAG)

# Пометить образ как latest и отправить его в Docker Hub
push-latest:
	$(DOCKER_LOGIN)
	$(DOCKER) tag $(IMAGE_NAME):$(TAG) $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):latest
	$(DOCKER) push $(DOCKER_REGISTRY_URL)/$(DOCKER_REGISTRY_PATH)/$(IMAGE_NAME):latest

# Удалить Docker образ
clean:
	@$(DOCKER) rmi -f $(IMAGE_NAME):$(TAG) 2>/dev/null