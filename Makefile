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

.PHONY: schema
schema:
ifeq ($(OS),Windows_NT)
	powershell -NoProfile -ExecutionPolicy Bypass -File scripts/gen-schema.ps1
else
	bash scripts/gen-schema.sh
endif