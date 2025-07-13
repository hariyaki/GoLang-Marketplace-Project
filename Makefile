.PHONY: test build

test:
	go vet ./...
	go test ./...

build:
	docker build -t golang-marketplace-api:dev .