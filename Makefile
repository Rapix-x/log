.PHONY: _prep setup build lint test check-build check

check: lint test check-build
	go mod tidy

setup:
	brew install go
	brew install golangci-lint

_prep:
	go generate ./...

run:
	echo "can't run this one, sorry"

build: export CGO_ENABLED=0
build: _prep
	@go build log

lint: _prep
	golangci-lint run ./...

test: _prep
	go test ./...

check-build: export CGO_ENABLED=1
check-build: _prep
	go build -race log
