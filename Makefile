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
