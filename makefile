.PHONY: deps build run test

deps:
	go mod tidy
	go install github.com/spf13/cobra-cli@latest

build:
	go build -o bin/thresh .

run:
	go run main.go

test:
	go test -v ./...

clean:
	rm -rf bin/