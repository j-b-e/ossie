
build: goreleaser
	@goreleaser build --snapshot --clean

generate:
	@go generate ./...

clean:
	rm -f main ossie

go-mod-update:
	@go get -u ./...
	@go mod tidy

run:
	@go run .

goreleaser:
	@command -v gorelease || go install github.com/goreleaser/goreleaser/v2@latest

release: gorelease
	@gorelease release