build:
	@go build -o bin/gobank

run: 
	@air 5000

test:
	@go test -v ./...
