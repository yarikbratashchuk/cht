GO ?= go

# Install from source.
install:
	@$(GO) install ./...

.PHONY: install
