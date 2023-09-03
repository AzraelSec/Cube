.PHONY: build

build:
	go build -o repl ./cmd/repl/main.go
	go build -o cube ./cmd/interpreter/main.go
