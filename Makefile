# Image URL to use all building/pushing image targets
REGISTRY ?= quay.io
REPOSITORY ?= $(REGISTRY)/opstree

get-depends:
	go get -v ./...

build-code:	get-depends
	go build -o druid-exporter

build-image:
	docker build -t quay.io/opstree/druid-exporter:v0.9 -f Dockerfile .

check-fmt:
	test -z "$(shell gofmt -l .)"

lint:
	OUTPUT="$(shell go list ./...)"; golint -set_exit_status $$OUTPUT

vet:
	VET_OUTPUT="$(shell go list ./...)"; GO111MODULE=on go vet $$VET_OUTPUT

test:
	go test -v -coverprofile=coverage.txt ./...

golangci-lint:
	golangci-lint run ./...
