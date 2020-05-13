export GO111MODULE = on

build:
	go build -o output/autowire cmd/main.go
.PHONY: build

test:
	go test ./...
.PHONY: test

vendor:
	go mod vendor
.PHONY: vendor

clean:
	rm -rf output
.PHONY: clean
