.POSIX:
.SUFFIXES:
.SHELL ?= /bin/sh

PFTEST ?= bin/pftest
GO_FILES := cmd/pftest/main.go

all: $(PFTEST)

$(PFTEST): $(GO_FILES)
	CGO_ENABLED=0 GOOS=openbsd go build -o $(PFTEST) ./cmd/pftest/main.go

clean:
	rm -rf $(PFTEST)

.PHONY: all clean
