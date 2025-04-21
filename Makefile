PACKAGE = $(shell go list -m)
VERSION ?= $(shell git describe --exact-match --tags 2> /dev/null || head -1 CHANGELOG.md 2> /dev/null | cut -d ' ' -f 2)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")
COMMIT ?= $(shell git rev-parse HEAD)
LDFLAGS = -ldflags "-w -X ${PACKAGE}/internal/version.Version=${VERSION} -X ${PACKAGE}/internal/version.BuildDate=${BUILD_DATE} -X ${PACKAGE}/internal/version.Commit=${COMMIT}"
TAGS =

.PHONY: *
build-binary: ## build a binary
	go build -tags '${TAGS}' ${LDFLAGS} -o bin/app

run-rest:
	make build-binary && ./bin/app rest
run-grpc:
	make build-binary && ./bin/app grpc
run-graphql:
	make build-binary && ./bin/app graphql
run-kafka:
	make build-binary && ./bin/app kafka all

up:
	docker compose -f .docker/docker-compose.yml up -d
down:
	docker compose -f .docker/docker-compose.yml down --remove-orphans

migrate-up:
	go run main.go migrate mysql up; \
	go run main.go migrate postgres up; \
	true

migrate-down:
	go run main.go migrate mysql up; \
	go run main.go migrate postgres down; \
	true

lint:
	golangci-lint run -v

gen-grpc:
	${UTILS_COMMAND} buf generate -v --template api/grpc/buf.gen.yaml api/grpc
lint-grpc:
	${UTILS_COMMAND} buf lint api/grpc

gen-graphql:
	${UTILS_COMMAND} go get github.com/99designs/gqlgen@latest && go run github.com/99designs/gqlgen generate --config api/graphql/gqlgen.yml

test:
	CGO_ENABLED=1 go test -tags mock,integration -race -cover ./...
test-no-cache:
	CGO_ENABLED=1 go test -tags mock,integration -race -cover -count=1 ./...
