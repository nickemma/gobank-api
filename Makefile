build: 
	@go build -o bin/gobank-api

run: build
	@./bin/gobank-api

test:
	@go test -v ./...