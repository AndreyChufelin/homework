BIN := "./bin/calendar"
SCHEDULER := "./bin/scheduler"
SENDER := "./bin/sender"
DOCKER_IMG="calendar:develop"
DB_DRIVER=postgres
DB_STRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run-calendar: build
	$(BIN) -config ./configs/config-calendar.toml

run-scheduler: build
	$(SCHEDULER) -config ./configs/config-scheduler.toml

run-sender: build
	$(SENDER) -config ./configs/config-sender.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

migrate:
	goose -dir migrations $(DB_DRIVER) $(DB_STRING) up

migrate-down:
	goose -dir migrations $(DB_DRIVER) $(DB_STRING) down

up:
	docker-compose -f ./deployments/docker-compose.yaml up

down:
	docker-compose -f ./deployments/docker-compose.yaml down


.PHONY: build run build-img run-img version test lint
