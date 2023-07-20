.PHONY: build clean test test-race

# 如果找不到 tag 则使用 HEAD commit
VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null || git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BIN=ddns-go
DIR_SRC=.
DOCKER_ENV=DOCKER_BUILDKIT=1
DOCKER=$(DOCKER_ENV) docker

GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)' -extldflags -static -s -w" -trimpath
GO=$(GO_ENV) $(shell which go)
GOROOT=$(shell `which go` env GOROOT)
GOPATH=$(shell `which go` env GOPATH)

build: $(DIR_SRC)/main.go
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

build_docker_image:
	@$(DOCKER) build -f ./Dockerfile -t ddns-go:$(VERSION) .

test:
	@$(GO) test ./...

test-race:
	@$(GO) test -race ./...

# clean all build result
clean:
	@$(GO) clean ./...
	@rm -f $(BIN)
	@rm -rf ./dist/*