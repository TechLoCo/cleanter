# Const
#===============================================================
OS                   := $(shell uname | tr A-Z a-z )
SHELL                := /bin/bash
BUILD_OPTIONS        := -tags netgo -installsuffix netgo
BIN_DIR              := bin

# GO:必要なツール類をセットアップします
setup:
ifeq ($(shell command -v make2help 2> /dev/null),)
	go get -u github.com/Songmu/make2help/cmd/make2help
endif
ifeq ($(shell command -v golint 2> /dev/null),)
	go get -u golang.org/x/lint/golint
endif
ifeq ($(shell command -v goreturns 2> /dev/null),)
	go get -u github.com/sqs/goreturns
endif

## GO:全てのソースの整形を行います
.PHONY: fmt
fmt:
	for pkg in $$(go list -f {{.Dir}} ./... | grep -v ^$$(pwd)$$ ); do \
		goreturns -w $$pkg; \
	done

## GO:全てのソースのLINTを実行します
.PHONY: lint
lint:
	for pkg in $$(go list ./...); do \
		golint -set_exit_status $$pkg; \
	done

## GO:ユニットテストを実行します
.PHONY: test
test:
	go test $$(go list ./...)

## GO:ビルドを実行します
.PHONY: build
build: .go-set-revision
	$(eval ldflags  := -X 'main.revision=$(revision)' -extldflags '-static')
	GOOS=$(OS) GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(ldflags)" -o $(BIN_DIR)/$(ENTRY_POINT) $(BUILD_OPTIONS) main.go

.PHONY: release
release: setup fmt lint build test

.go-set-revision:
	$(eval REVISION := $(shell if [[ $$REV = "" ]]; then git rev-parse --short HEAD; else echo $$REV;fi;))