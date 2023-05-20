# Provider metadata and versionning
PACKAGE_NAME = neuvector
FULL_PACKAGE_NAME = github.com/theobori/terraform-provider-neuvector
VERSION = $(shell git describe --tags --always)
PROVIDER_NAME = terraform-provider-$(PACKAGE_NAME)_$(VERSION)
PROVIDER_PATH = bin/$(PROVIDER_NAME)

# Formatted Go files
GOFMT_FILES ?= $(shell find . -name "*.go")

# Terraform plugins directory
# TODO

default: build

build:
	go build -o $(PROVIDER_PATH)

install:
	echo $(HOME)

fmt:
	gofmt -w $(GOFMT_FILES)

clean:
	$(RM) -r $(PROVIDER_PATH)

uninstall:
	echo uninstall

fclean: clean uninstall

re: fclean build install

.PHONY: \
	build \
	install \
	uninstall \
	fmt \
	clean \
	fclean \
	re
