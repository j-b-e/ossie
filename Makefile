generate:
	@go generate ./...

build:
	@go build main.go

clean:
	rm -f main ossie

go-mod-update:
	@go get -u ./...
	@go mod tidy
