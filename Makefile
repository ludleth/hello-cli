BINARY_NAME=hello-cli
# Try to get the latest Git tag, otherwise use 'dev'
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
# Get current commit hash
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
# Get build time
BUILD_TIME=$(shell date +%FT%T%z)

LDFLAGS=-ldflags "-s -w -X github.com/ludleth/hello-cli/cmd.Version=${VERSION} \
                  -X github.com/ludleth/hello-cli/cmd.Commit=${COMMIT} \
                  -X github.com/ludleth/hello-cli/cmd.BuildTime=${BUILD_TIME}"

.PHONY: build
build:
	mkdir -p ./build/
	go build ${LDFLAGS} -o ./build/${BINARY_NAME} main.go

.PHONY: install
install:
	go install ${LDFLAGS}


.PHONY: test
test:
	go test ./... -v
