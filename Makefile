GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

PKG_NAME := github.com/fatindeed/github-proxy

EXTENSION:=
ifeq ($(GOOS),windows)
  EXTENSION:=.exe
endif

STATIC_FLAGS=CGO_ENABLED=0

GIT_TAG?=$(shell git rev-parse --short HEAD)

LDFLAGS="-s -w -X $(PKG_NAME)/cmd.Version=${GIT_TAG}"
GO_BUILD=$(STATIC_FLAGS) go build -trimpath -ldflags=$(LDFLAGS)

GITHUBPROXY_BINARY?=bin/github-proxy
GITHUBPROXY_BINARY_WITH_EXTENSION=$(GITHUBPROXY_BINARY)$(EXTENSION)

TAGS:=
ifdef BUILD_TAGS
  TAGS=-tags $(BUILD_TAGS)
  LINT_TAGS=--build-tags $(BUILD_TAGS)
endif

all: test build

.PHONY: test
test:
	go test $(TAGS) -cover ./...

.PHONY: lint
lint:
	golangci-lint run $(LINT_TAGS) --timeout 10m0s ./...

.PHONY: build
build:
	$(GO_BUILD) $(TAGS) -o $(GITHUBPROXY_BINARY_WITH_EXTENSION) main.go

.PHONY: cross
cross:
	GOOS=linux   GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(GITHUBPROXY_BINARY)-linux-x86_64 main.go
	GOOS=darwin  GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(GITHUBPROXY_BINARY)-darwin-x86_64 main.go
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(TAGS) -o $(GITHUBPROXY_BINARY)-windows-x86_64.exe main.go