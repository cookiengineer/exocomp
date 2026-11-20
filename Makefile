GO ?= go
ARCH ?= $(shell $(GO) env GOARCH)
TOOLCHAIN := $(GO) -C toolchain

.PHONY: all build quick linux darwin windows webapp lint lint-toolchain test test-toolchain test-integration install clean

all: build

build:
	$(TOOLCHAIN) run ./build --arch=$(ARCH)

quick:
	$(TOOLCHAIN) run ./build --quick --arch=$(ARCH)

linux:
	$(TOOLCHAIN) run ./build --arch=$(ARCH) linux

darwin:
	$(TOOLCHAIN) run ./build --arch=$(ARCH) darwin

windows:
	$(TOOLCHAIN) run ./build --arch=$(ARCH) windows

webapp:
	$(TOOLCHAIN) run ./build-webapp

lint:
	$(TOOLCHAIN) run ./lint

lint-toolchain:
	$(TOOLCHAIN) run ./lint -root .

test:
	$(GO) -C source test ./...

test-toolchain:
	$(TOOLCHAIN) test ./...

test-integration:
	$(TOOLCHAIN) run ./test

install:
	$(TOOLCHAIN) run ./install

clean:
	$(TOOLCHAIN) run ./clean
