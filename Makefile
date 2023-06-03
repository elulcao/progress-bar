APP_NAME := progress_bar
GOOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH := $(subst x86_64,amd64,$(shell uname -m))
GO_FILES := $(shell find . -type f -not -path './vendor/*' -name '*.go')

.PHONY all: clean test build
.PHONY test: go-fmt go-vet go-test
.PHONY build: clean
.PHONY install: build

go-test:
	@echo "Running go test"
	go test ./...

go-vet:
	@echo "Running go vet"
	go vet ./...

go-fmt:
	@echo "Running go fmt"
	go fmt ./...

build: clean
	env GO111MODULE=auto GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o "$(APP_NAME)"

install: build
	install -m 0755 "$(APP_NAME)" "$(GOPATH)/bin/"

clean:
	find . -type f -name "$(APP_NAME).*" -exec rm -rf {} \;
