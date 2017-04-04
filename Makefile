GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_DEPS_UPDATE=$(GO_CMD) get -d -v -u
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GO_LINT=golint

COVERAGE_DIR=.cover
COVERAGE_PROFILE=$(COVERAGE_DIR)/cover.out
COVERAGE_MODE=atomic

.PHONY: all test coverage coverall lint

all: test

test:
	$(GO_TEST)

coverage:
	@rm -rf $(COVERAGE_DIR)
	@mkdir $(COVERAGE_DIR)
	$(GO_TEST) -covermode=$(COVERAGE_MODE) -coverprofile=$(COVERAGE_PROFILE)

coverall: coverage
	goveralls -coverprofile=$(COVERAGE_PROFILE)

lint:
	@$(GO_LINT)

get-deps-tests:
	$(GO_DEPS) github.com/stretchr/testify
	$(GO_DEPS) github.com/golang/lint
	$(GO_DEPS) github.com/mattn/goveralls
