# Provider metadata and versionning
PROVIDER = neuvector
VERSION = $(shell git describe --tags --always)

# Terraform metadata for installation
BIN = terraform-provider-$(PROVIDER)
HOSTNAME = github.com
NAMESPACE = theobori
VERSION = 1.0.0
OS_ARCH = linux_amd64

# Terraform
TF_PLUGINS_DIR = ~/.terraform.d/plugins
TF_CACHE = $(shell find examples/ -name ".terraform*")

# Output binary paths
BIN_DIR = $(TF_PLUGINS_DIR)/$(HOSTNAME)/$(NAMESPACE)/$(PROVIDER)/$(VERSION)/$(OS_ARCH)
BIN_PATH = $(BIN_DIR)/$(BIN)

# Formatted Go files
GOFMT_FILES ?= $(shell find . -name "*.go")


default: install

build:
	go build -o $(BIN)

fmt:
	gofmt -w $(GOFMT_FILES)

clean_test:
	go clean -testcache

clean: clean_test
	$(RM) -r $(BIN)

uninstall: clean
	$(RM) -r $(BIN_DIR)

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/$(BINARY)_$(VERSION)_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/$(BINARY)_$(VERSION)_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/$(BINARY)_$(VERSION)_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/$(BINARY)_$(VERSION)_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/$(BINARY)_$(VERSION)_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/$(BINARY)_$(VERSION)_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(BINARY)_$(VERSION)_windows_amd64

install: build
	mkdir -p $(BIN_DIR)
	mv $(BIN) $(BIN_PATH)

fclean: clean uninstall
	$(RM) -r $(TF_CACHE)

re: fclean build install

.PHONY: \
	build \
	install \
	uninstall \
	fmt \
	clean_test \
	clean \
	release \
	re
