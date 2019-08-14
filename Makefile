.PHONY: install format test cover

install:
	go mod download

format:
	gofmt -s -w .

test:
	go test -v ./...

bench:
	go test -run=__ -bench=. -cpuprofile test/profile_cpu.out ./...
	go tool pprof -svg test/profile_cpu.out > test/profile_cpu.svg

cover:
	go test -v -coverpkg=./... -coverprofile ./test/coverage.out ./...
	go tool cover -html ./test/coverage.out
