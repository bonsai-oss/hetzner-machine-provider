build: export GOOS=linux
build: export CGO_ENABLED=0

NAME := "hmp"
VERSION := $(shell git describe --tags --always --dirty)

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o bin/$(NAME) cmd/$(NAME)/main.go