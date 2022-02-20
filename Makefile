APP_NAME  := progress_bar
GO_FILES  := $(shell find . -type f -not -path './vendor/*' -name '*.go')
UNAME_S   := $(shell uname -s)
goos_type := ""

ifeq ($(UNAME_S),Linux)
	goos_type = "linux"
endif

ifeq ($(UNAME_S),Darwin)
	goos_type = "darwin"
endif

.PHONY all: clean test build

.PHONY test: go-fmt go-vet go-lint go-test

go-test:
	@echo "Running go test"
	go test ./...

go-vet:
	@echo "Running go vet"
	go vet ./...

go-lint:
	@echo "Running go lint"
	go list ./... | grep -v upstream-go | xargs $(shell go env GOPATH)/bin/golint -set_exit_status=1

go-fmt:
	@echo "Running go fmt"
	go fmt ./...

.PHONY: build

build: clean
	env GO111MODULE=auto GOOS=darwin GOARCH=amd64       go build -v -o "$(APP_NAME).darwin" && \
	env GO111MODULE=auto GOOS=linux  GOARCH=amd64       go build -v -o "$(APP_NAME).linux"

.PHONY: install

install: build
	install -m 0755 "$(APP_NAME).$(goos_type)" "/usr/local/bin/$(APP_NAME)"

.PHONY: clean

clean:
	find . -type f -name "$(APP_NAME).*" -exec rm -rf {} \;
