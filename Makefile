GO ?= go
TOOLCHAIN := $(GO) -C toolchain

.PHONY: all build quick linux darwin windows webapp test test-integration install clean

all: build

build:
	$(TOOLCHAIN) run build.go

quick:
	$(TOOLCHAIN) run build.go --quick

linux:
	$(TOOLCHAIN) run build.go linux

darwin:
	$(TOOLCHAIN) run build.go darwin

windows:
	$(TOOLCHAIN) run build.go windows

webapp:
	$(TOOLCHAIN) run build-webapp.go

test:
	$(GO) -C source test ./...

test-integration:
	$(TOOLCHAIN) run test.go

install:
	$(TOOLCHAIN) run install.go

clean:
	$(TOOLCHAIN) run clean.go
