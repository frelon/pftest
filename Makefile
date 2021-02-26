.POSIX:
.SUFFIXES:
.SHELL ?= /bin/sh

PFTEST ?= bin/pftest
GO_FILES := $(wildcard **/*.go)
GO_MAIN := ./cmd/pftest/main.go

all: $(PFTEST)

$(PFTEST): $(GO_FILES)
	go build -o $(PFTEST) $(GO_MAIN)

test:
	go test ./...

clean:
	rm $(PFTEST)

lint:
	golangci-lint run

.PHONY: all clean test lint echo
