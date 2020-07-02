all: build test format presubmit

test: 
	@echo "--running tests--"
	@go test -race ./...

tidy:
	@echo "--tidying up--"
	@go mod tidy

format:
	@echo "--formatting code--"
	@go fmt ./...

vet:
	@echo "--vetting code--"
	@go vet ./...

build:
	@echo "--building code--"
	@go build -o $(PWD)/systemd_exporter github.com/dashpole/systemd_exporter

docker:
	@echo "--building docker container--"
	@docker build -t systemd_exporter:$(shell git rev-parse --short HEAD) -f deploy/Dockerfile .

presubmit: vet lint

lint:
	@echo "--linting code--"
	@$(shell go env GOPATH)/bin/golangci-lint run

.PHONY: all build docker format presubmit test tidy vet