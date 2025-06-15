main_package_path = ./cmd
binary_name = ghosty_link

DOCKER_USERNAME = adanil
APPLICATION_NAME = ghosty_link
GIT_HASH ?= $(shell git log --format="%h" -n 1)

PROJECT_PATH = $(CURDIR)
BINDIR = $(PROJECT_PATH)/bin

GOLANGCI = $(BINDIR)/golangci-lint

.PHONY: dev-up
dev-up:
	docker compose --file ./docker-compose.dev.yml up --detach

.PHONY: dev-down
dev-down:
	docker compose --file ./docker-compose.dev.yml down

.PHONY: install-lint
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.6

.PHONY: lint
lint: install-lint
	$(GOLANGCI) run --config=$(PROJECT_PATH)/.golangci.yml -v $(PROJECT_PATH)/...

.PHONY: fmt
fmt: install-lint
	$(GOLANGCI) fmt -v $(PROJECT_PATH)/...
	go fmt ./...

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the application
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o=./bin/${binary_name} ${main_package_path}

## run: run the  application
.PHONY: run
run: build
	./bin/${binary_name}

## build-docker: build docker image
.PHONY: build-docker
build-docker: build
	docker build --tag ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH} .

## release-docker: push docker images with latest version
.PHONY: release-docker
release-docker: build-docker
	docker tag  ${DOCKER_USERNAME}/${APPLICATION_NAME}:${GIT_HASH} ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}:latest

.PHONY: generate-oapi
generate-oapi:
	oapi-codegen -config api/public/cfg.yaml api/public/api.yaml