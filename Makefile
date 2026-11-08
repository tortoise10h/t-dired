.PHONY: format build run

format:
	go fmt ./...

build:
	go build

run:
	go run main.go
