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

.PHONY: all test lint dep

all: dep test lint

test:
	./tools/test.sh

lint:
	./tools/lint.sh

dep:
	curl -fsSL -o /tmp/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64
	chmod +x /tmp/dep
	/tmp/dep ensure -vendor-only
