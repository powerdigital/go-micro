PACKAGE = $(shell go list -m)
VERSION ?= $(shell git describe --exact-match --tags 2> /dev/null || head -1 CHANGELOG.md 2> /dev/null | cut -d ' ' -f 2)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")
COMMIT ?= $(shell git rev-parse HEAD)
LDFLAGS = -ldflags "-w -X ${PACKAGE}/internal/version.Version=${VERSION} -X ${PACKAGE}/internal/version.BuildDate=${BUILD_DATE} -X ${PACKAGE}/internal/version.Commit=${COMMIT}"
TAGS =

.PHONY: *
build-binary: ## build a binary
	go build -tags '${TAGS}' ${LDFLAGS} -o bin/app

run-rest-server:
	make build-binary && ./bin/app rest
run-grpc-server:
	make build-binary && ./bin/app grpc

# Generate gRPC server using proto files
# Requires installed buf (https://buf.build/docs/installation)
gen-grpc:
	${UTILS_COMMAND} buf generate -v --template api/buf.gen.yaml api/grpc

# Proto specifications validation
# Requires installed buf (https://buf.build/docs/installation)
lint-grpc:
	${UTILS_COMMAND} buf lint api/grpc

test:
	go test -tags mock,integration -race -cover ./...
test-no-cache:
	go test -tags mock,integration -race -cover -count=1 ./...
