# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /src

# go modules first —> layer caching
COPY go.mod go.sum ./
RUN go mod download

# rest of the source
COPY . .

# Compile only the cmd/api package
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/marketplace ./cmd/api

FROM scratch

COPY --from=builder /go/bin/marketplace /marketplace
EXPOSE 8080
ENTRYPOINT ["/marketplace"]