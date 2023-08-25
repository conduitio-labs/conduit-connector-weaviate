.PHONY: build test test-integration generate install-paramgen install-mocken

VERSION				=  $(shell git describe --tags --dirty --always)
MOCKGEN_VERSION		?= v0.2.0
PARAMGEN_VERSION	?= v0.7.2
GOLANG_CI_LINT_VER	:= v1.54.2

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-weaviate.version=${VERSION}'" -o conduit-connector-weaviate cmd/connector/main.go

install-mocken:
	go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)

test: generate
	go test $(GOTEST_FLAGS) -race ./...

test-integration:
	# run required docker containers, execute integration tests, stop containers after tests
	docker compose -f test/docker-compose.yml up -d
	go test $(GOTEST_FLAGS) -v -race ./...; ret=$$?; \
		docker compose -f test/docker-compose.yml down; \
		exit $$ret

generate: install-mocken install-paramgen
	go generate ./...

install-paramgen:
	go install github.com/conduitio/conduit-connector-sdk/cmd/paramgen@$(PARAMGEN_VERSION)


.PHONY: install-golangci-lint
install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_CI_LINT_VER)

.PHONY: lint
lint: install-golangci-lint
	golangci-lint run -v
