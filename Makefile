BINARY := googlemapscli
PKG    := ./cmd/googlemapscli

.PHONY: build test lint clean

build:
	go build -o $(BINARY) $(PKG)

test:
	go test ./... -count=1

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

coverage:
	go test ./... -coverprofile=coverage.txt -count=1
	go tool cover -func=coverage.txt
