BINARY=demo

build:
	go build -o $(BINARY) ./cmd/demo

run:
	go run ./cmd/demo

device:
	go run ./mockdevice

test-go:
	go test ./...

test-py:
	pytest -q uat/python
