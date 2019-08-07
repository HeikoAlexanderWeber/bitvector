ifeq ($(OS),Windows_NT)
	HOST_OS := windows
	PROGRAM := program.exe
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        HOST_OS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        HOST_OS := darwin
    endif
	PROGRAM := program
endif

ifndef VERSION
	VERSION := $(shell sed -n 1p ./version)
endif

.PHONY: install format test cover

install:
	go mod download

format:
	gofmt -s -w .

test:
	go test -v ./...

cover:
	go test -v -coverpkg=./... -coverprofile ./test/coverage.out ./...
	go tool cover -html ./test/coverage.out
