.PHONY: build clean test test-race

VERSION=0.0.1
BIN=ddns-go
DIR_SRC=.
DOCKER_CMD=docker

GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X main.version=$(VERSION) -X 'main.buildTime=`date`' -extldflags -static"
GO=$(GO_ENV) $(shell which go)
GOROOT=$(shell `which go` env GOROOT)
GOPATH=$(shell `which go` env GOPATH)

build: init bindata $(DIR_SRC)/main.go
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

build_docker_image:
	@$(DOCKER_CMD) build -f ./Dockerfile -t ddns-go:$(VERSION) .

init:
	@go get -u github.com/go-bindata/go-bindata/...

test:
	@$(GO) test ./...

test-race:
	@$(GO) test -race ./...

bindata:
	@go-bindata -pkg util -o util/staticPages.go static/pages/...
	@go-bindata -pkg static -o asserts/html.go -fs -prefix "static/" static/

dev:
	@go-bindata -debug -pkg util -o util/staticPages.go static/pages/...

# clean all build result
clean:
	@rm -f util/staticPages.go asserts/html.go
	@$(GO) clean ./...
	@rm -f $(BIN)
	@rm -rf ./dist/*
