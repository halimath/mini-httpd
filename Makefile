VERSION := 0.2.0
BUILD_TIMESTAMP := $(shell date -u --rfc-3339=seconds)
COMMIT := $(shell git rev-parse --short HEAD )

GO ?= go
GO_FLAGS ?=

TAR ?= tar
TAR_FLAGS ?= cfz

mini-httpd: main.go
	$(GO) $(GO_FLAGS) build -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X 'main.buildTimestamp=${BUILD_TIMESTAMP}'" -o $@ $<

.PHONY: all
all: linux windows darwin

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 $(MAKE) mini-httpd
	$(TAR) $(TAR_FLAGS) mini-httpd.linux-amd64.tar.gz mini-httpd

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 $(MAKE) mini-httpd
	$(TAR) $(TAR_FLAGS) mini-httpd.windows-amd64.tar.gz mini-httpd

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 $(MAKE) mini-httpd
	$(TAR) $(TAR_FLAGS) mini-httpd.darwin-amd64.tar.gz mini-httpd

.PHONY: clean
clean:
	rm -f mini-httpd*