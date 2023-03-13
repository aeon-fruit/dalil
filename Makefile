# App variables
APP_NAME=dalil
APP_ROOT=.
APP_ROOT_FILES=cmd/dalil/dalil.go
APP_BUILD_PATH=$(APP_ROOT)/.build
APP_OUTPUT=$(APP_BUILD_PATH)/$(APP_NAME)

# Tools variables
GO=go
GO_TEST=$(GO) test ./cmd/... ./internal/...
GO_TEST_VERBOSE=$(GO_TEST) -v
GO_DEBUG_OPTIONS=-gcflags='all=-N -l'

# Testing variables
TEST_MOCKS_PATH=$(APP_ROOT)/test/mocks

# Dependencies variables
DEP_VERSION_MOCKGEN=1.6.0

all: build test
.PHONY: all

deps:
.PHONY: deps

deps-test:
# TODO - Enable when tests are implemented
# $(GO) install github.com/golang/mock/mockgen@$(DEP_VERSION_MOCKGEN)
.PHONY: deps-test

build: download
	@mkdir -p .build
	$(GO) build -o $(APP_OUTPUT) $(APP_ROOT_FILES)
.PHONY: build

build-debug: download
	$(GO) build -o $(APP_OUTPUT) $(GO_DEBUG_OPTIONS) $(APP_ROOT_FILES)
.PHONY: build-debug

download:
	$(GO) mod tidy
	$(GO) mod download all
.PHONY: download

run:
	$(GO) run $(APP_ROOT_FILES)
.PHONY: run

test: build gen-test quick-test
.PHONY: test

quick-test:
	$(GO_TEST_VERBOSE)
.PHONY: quick-test

test-race: build gen-test quick-test-race
.PHONY: test-race

quick-test-race:
	$(GO_TEST_VERBOSE) -race
.PHONY: quick-test-race

gen-test: clean-gen-test
	@mkdir -p $(TEST_MOCKS_PATH)
# TODO - Enable when tests are implemented
# mockgen -source=$(APP_ROOT)/internal/pkg/<package>/<file>.go -destination=$(TEST_MOCKS_PATH)/<package>/<file>_mock.go
.PHONY: gen-test

docker-build:
	echo "Build $(APP_NAME) Docker image"
# TODO - Enable when Dockerfile is specified
# docker build -t $(APP_NAME) .
.PHONY: docker-build

clean-gen-test:
	rm -rf $(TEST_MOCKS_PATH)
.PHONY: clean-gen-test

clean: clean-gen-test
	rm -rf $(APP_BUILD_PATH)
.PHONY: clean
