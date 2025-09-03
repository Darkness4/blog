GO_SRCS := $(shell find . -type f -name '*.go' -a ! \( -name 'zz_generated*' -o -name '*_test.go' \))
GO_TESTS := $(shell find . -type f -name '*_test.go')
TAG_NAME = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
TAG_NAME_DEV = $(shell git describe --tags --abbrev=0 2>/dev/null)
VERSION_CORE = $(shell echo $(TAG_NAME))
VERSION_CORE_DEV = $(shell echo $(TAG_NAME_DEV))
GIT_COMMIT = $(shell git rev-parse --short=7 HEAD)
VERSION = $(or $(and $(TAG_NAME),$(VERSION_CORE)),$(and $(TAG_NAME_DEV),$(VERSION_CORE_DEV)-dev),$(GIT_COMMIT))
DB_DSN ?= $(shell cat .env | grep DB_DSN | cut -d '=' -f 2)
DB_DSN := $(or $(DB_DSN),$(shell cat .env.local | grep DB_DSN | cut -d '=' -f 2))
MIGRATION_NAME ?= migration_name

golint :=  $(shell which golangci-lint)
ifeq ($(golint),)
golint := $(shell go env GOPATH)/bin/golangci-lint
endif

wgo :=  $(shell which wgo)
ifeq ($(wgo),)
wgo := $(shell go env GOPATH)/bin/wgo
endif

goose := $(shell which goose)
ifeq ($(goose),)
goose := $(shell go env GOPATH)/bin/goose
endif

sqlc := $(shell which sqlc)
ifeq ($(sqlc),)
sqlc := $(shell go env GOPATH)/bin/sqlc
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

.PHONY: migration
migration:
	$(goose) -dir db/migrations -s create $(MIGRATION_NAME) sql

.PHONY: up
up: $(MIGRATIONS)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	$(goose) -dir db/migrations postgres $(DB_DSN) up

.PHONY: drop
drop:
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	$(goose) -dir db/migrations postgres $(DB_DSN) reset

.PHONY: sql
sql: $(sqlc)
ifndef DB_DSN
	$(error DB_DSN is not defined)
endif
	@DB_DSN=$(subst cockroachdb://,postgres://,$(DB_DSN)) $(sqlc) generate
	@echo "sqlc: done"

.PHONY: watch
watch: $(wgo)
	$(wgo) -xdir "bin/" -xdir "web/gen/" sh -c 'while nc -vz 127.0.0.1 3000 > /dev/null 2>&1; do sleep 1; done; make run || exit 1' --signal SIGTERM

.PHONY: unit
unit:
	go test -race -covermode=atomic -timeout=30s ./...

.PHONY: lint
lint: $(golint)
	$(golint) run ./...

.PHONY: clean
clean:
	rm -rf bin/

$(goose):
	go install -tags 'no_clickhouse,no_mssql,no_mysql,no_turso,no_vertica,no_ydb' github.com/pressly/goose/v3/cmd/goose

$(sqlc):
	go install github.com/sqlc-dev/sqlc/cmd/sqlc

$(wgo):
	go install github.com/bokwoon95/wgo@latest

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
