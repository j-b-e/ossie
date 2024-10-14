
build:
	@go build .

generate:
	@go generate ./...

clean:
	rm -f main ossie

go-mod-update:
	@go get -u ./...
	@go mod tidy

run:
	@go run .