
build: goreleaser
	@goreleaser build --snapshot --clean

generate:
	@go generate ./...

clean:
	rm -f main ossie dist/

go-mod-update:
	@go get -u ./...
	@go mod tidy

run:
	@go run cmd/ossie/ossie.go

goreleaser:
	@command -v gorelease || go install github.com/goreleaser/goreleaser/v2@latest

release: goreleaser
	@goreleaser release

test:
	@go test ./...
