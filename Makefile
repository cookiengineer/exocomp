GO ?= go
ARCH ?= $(shell $(GO) env GOARCH)
TOOLCHAIN := $(GO) -C toolchain

.PHONY: all build quick linux darwin windows webapp test test-integration install clean

all: build

build:
	$(TOOLCHAIN) run build.go --arch=$(ARCH)

quick:
	$(TOOLCHAIN) run build.go --quick --arch=$(ARCH)

linux:
	$(TOOLCHAIN) run build.go linux --arch=$(ARCH)

darwin:
	$(TOOLCHAIN) run build.go darwin --arch=$(ARCH)

windows:
	$(TOOLCHAIN) run build.go windows --arch=$(ARCH)

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
