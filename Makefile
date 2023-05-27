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
	re
