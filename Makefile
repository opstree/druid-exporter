# # Image URL to use all building/pushing image targets
# REGISTRY ?= quay.io
# REPOSITORY ?= $(REGISTRY)/opstree
# IMAGE ?= druid-exporter
# VERSION ?= v0.11

# Image URL to use all building/pushing image targets
REPOSITORY ?= iunera
IMAGE ?= druid-exporter
VERSION ?= v0.12.6

get-depends:
	go get -v ./...

build-code:	get-depends
	go build -o druid-exporter

build-image:
	docker build -t ${REPOSITORY}/${IMAGE}:${VERSION} -f Dockerfile .

push-image: build-image
	docker push ${REPOSITORY}/${IMAGE}:${VERSION}

fmt:
	go fmt ./...

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
