build:
	@go build -o bin/gobank

run: 
	@air 5000

test:
	@go test -v ./...

tailwind:
	@npx tailwindcss -i ./static/input.css -o ./dist/output.css --watch
