BINARY=projscan

build:
	go build -o $(BINARY)

run:
	go run . .

test:
	go test ./...

release-snapshot:
	goreleaser release --snapshot --clean