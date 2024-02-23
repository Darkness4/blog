GO_SRCS := $(shell find . -type f -name '*.go' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
GO_TESTS := $(shell find . -type f -name '*_test.go')
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
TAG_NAME_DEV = $(shell git describe --tags --abbrev=0 2>/dev/null)
VERSION_CORE = $(shell echo $(TAG_NAME))
VERSION_CORE_DEV = $(shell echo $(TAG_NAME_DEV))
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
VERSION = $(or $(and $(TAG_NAME),$(VERSION_CORE)),$(and $(TAG_NAME_DEV),$(VERSION_CORE_DEV)-dev),$(GIT_COMMIT))

golint :=  $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

gow := $(shell which gow)
ifeq ($(gow),)
gow := $(shell go env GOPATH)/bin/gow
endif

.PHONY: bin/blog
bin/blog: $(GO_SRCS) generate
	go build -ldflags "-s -w -X main.version=${VERSION}" -o "$@" ./main.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: run
run: generate
	go run ./main.go

.PHONY: watch
watch: $(gow)
	$(gow) -e=go,mod,html,tmpl,env,local,htmx run ./main.go

.PHONY: unit
unit:
	go test -race -covermode=atomic -timeout=30s ./...

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: clean
clean:
	rm -rf bin/

$(gow):
	go install github.com/mitranim/gow@latest

$(golint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: version
version:
	@echo VERSION_CORE=${VERSION_CORE}
	@echo VERSION_CORE_DEV=${VERSION_CORE_DEV}
	@echo VERSION=${VERSION}

.PHONY: lighthouse
lighthouse:
	pnpm dlx unlighthouse --site http://localhost:3000
