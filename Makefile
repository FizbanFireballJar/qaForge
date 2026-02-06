.PHONY: build build-all test run clean

build:
	go build -o qaforge cmd/qaforge/*.go

build-all:
	GOOS=darwin GOARCH=arm64 go build -o qaforge-darwin-arm64 cmd/qaforge/*.go
	GOOS=darwin GOARCH=amd64 go build -o qaforge-darwin-amd64 cmd/qaforge/*.go
	GOOS=linux GOARCH=amd64 go build -o qaforge-linux-amd64 cmd/qaforge/*.go
	GOOS=windows GOARCH=amd64 go build -o qaforge-windows-amd64.exe cmd/qaforge/*.go

test:
	go test ./...

run:
	go run cmd/qaforge/*.go

clean:
	rm -f qaforge qaforge-*