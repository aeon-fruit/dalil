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
TEST_COVERAGE_PATH=$(APP_BUILD_PATH)/coverage
MOCKGEN=mockgen

# Dependencies variables
DEP_VERSION_MOCKGEN=v1.6.0
DEP_VERSION_GINKGO=v2.9.1
DEP_VERSION_GOMEGA=v1.27.4

all: build test
.PHONY: all

deps:
.PHONY: deps

deps-test:
	$(GO) install github.com/golang/mock/mockgen@$(DEP_VERSION_MOCKGEN)
	$(GO) install github.com/onsi/ginkgo/v2/ginkgo@$(DEP_VERSION_GINKGO)
	$(GO) get github.com/onsi/ginkgo/v2@$(DEP_VERSION_GINKGO)
	$(GO) get github.com/onsi/gomega@$(DEP_VERSION_GOMEGA)
.PHONY: deps-test

build: download
	@mkdir -p $(APP_BUILD_PATH)
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
	$(MOCKGEN) -source=$(APP_ROOT)/internal/pkg/tasks/dao/repository.go -destination=$(TEST_MOCKS_PATH)/tasks/dao/repository_mock.go
	$(MOCKGEN) -source=$(APP_ROOT)/internal/pkg/tasks/service/service.go -destination=$(TEST_MOCKS_PATH)/tasks/service/service_mock.go
	$(MOCKGEN) -source=$(APP_ROOT)/internal/pkg/tasks/controller/controller.go -destination=$(TEST_MOCKS_PATH)/tasks/controller/controller_mock.go
	$(MOCKGEN) -source=$(APP_ROOT)/internal/pkg/log/log.go -destination=$(TEST_MOCKS_PATH)/log/log_mock.go
	$(MOCKGEN) -destination=$(TEST_MOCKS_PATH)/logr/logr_mock.go github.com/go-logr/logr LogSink
.PHONY: gen-test

test-coverage: clean-test-coverage gen-test build
	@mkdir -p $(TEST_COVERAGE_PATH)
	$(GO) test ./internal/... -covermode=count -coverprofile=$(TEST_COVERAGE_PATH)/coverage.out
	$(GO) tool cover -html $(TEST_COVERAGE_PATH)/coverage.out -o $(TEST_COVERAGE_PATH)/coverage.html
.PHONY: test-coverage

docker-build:
	echo "Build $(APP_NAME) Docker image"
	docker build -t $(APP_NAME) .
.PHONY: docker-build

clean-gen-test:
	rm -rf $(TEST_MOCKS_PATH)
.PHONY: clean-gen-test

clean-test-coverage:
	rm -rf $(TEST_COVERAGE_PATH)
.PHONY: clean-test-coverage

clean: clean-gen-test clean-test-coverage
	rm -rf $(APP_BUILD_PATH)
.PHONY: clean
