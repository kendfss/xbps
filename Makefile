ROOT := $(abspath .)
NAME := $(notdir $(ROOT))

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" || echo "unknown")

LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE)

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o ./bin/$(NAME) .

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: install
install:
	CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" .
	sudo ln -fs "$(shell go env GOPATH)/bin/$(NAME)" /usr/bin/$(NAME)
