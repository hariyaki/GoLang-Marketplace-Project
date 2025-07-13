SWAG := $(shell go env GOPATH)/bin/swag

.PHONY: docs
docs:
	@echo "› Generating Swagger docs"
	"$(SWAG)" init --dir ./cmd/api,./internal --output ./docs

test:
	go vet ./...
	go test ./...

build:
	docker build -t golang-marketplace-api:dev .