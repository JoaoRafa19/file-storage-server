build:
	@go build -o ./bin/fs

run: build
	@./bin/fs

test:
	@go test ./...

clear:
	rm -rf ./:4000_network/ ./bin
