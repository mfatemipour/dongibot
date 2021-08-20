PROJECT_NAME := dongibot
DONGIBOT_DOCKER_IMAGE := mfatemipour/dongibot
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags $(COMMIT) 2> /dev/null || echo "$(COMMIT)")
TAG = $(DONGIBOT_DOCKER_IMAGE):$(VERSION)
TAG_LATEST = $(DONGIBOT_DOCKER_IMAGE):latest
LINTER = golangci-lint
LINTER_VERSION = v1.36.0

build:
	@go build ./cmd/dongibot

run:
	@go run cmd/dongibot/main.go

test:
	@go test ./...

test-alpine:
	@go test -tags musl ./...

fmt:
	@go fmt ./...

dep:
	@go mod tidy

lintdeps:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(LINTER_VERSION)

lint:
	$(LINTER) run --config=.golangci-lint.yml ./...

lint-alpine:
	$(LINTER) run --build-tags=musl --config=.golangci-lint.yml ./...

clean:
	@go clean ./...
	@rm -f ./$(PROJECT_NAME)

build-alpine:
	@go build -tags musl ./cmd/dongibot/

docker: docker-build docker-push

docker-build:
	docker build --pull -t $(TAG) .

docker-push:
	docker push $(TAG)

docker-push-latest: docker-pull
	docker tag $(TAG) $(TAG_LATEST)
	docker push $(TAG_LATEST)

docker-pull:
	docker pull $(TAG)
