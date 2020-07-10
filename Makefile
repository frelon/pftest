.POSIX:
.SUFFIXES:
.SHELL ?= /bin/sh

PFTEST ?= bin/pftest
GO_FILES := $(wildcard *.go)

all: $(PFTEST)

$(PFTEST): $(GO_FILES)
	go build -o $(PFTEST) ./cmd/pftest/main.go
	# CGO_ENABLED=0 GOOS=openbsd go build -o $(PFTEST) ./cmd/pftest/main.go

test:
	go test ./...

clean:
	rm -rf $(PFTEST)

lint:
	golangci-lint run

.PHONY: all clean test lint
