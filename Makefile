export GO111MODULE=on

build: \
	build/example-simple \
	build/example-cache \
	build/example-httpsource \
	build/example-groupcache \
	build/example-advanced

build/example-simple:
	go build -o build/example-simple ./examples/simple

build/example-cache:
	go build -o build/example-cache ./examples/cache

build/example-httpsource:
	go build -o build/example-httpsource ./examples/httpsource

build/example-groupcache:
	go build -o build/example-groupcache ./examples/groupcache

build/example-advanced:
	go build -o build/example-advanced ./examples/advanced

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: \
	golangci-lint

GOLANGCI_LINT_VERSION=v1.17.0
GOLANGCI_LINT_DIR=$(shell go env GOPATH)/pkg/golangci-lint/$(GOLANGCI_LINT_VERSION)
$(GOLANGCI_LINT_DIR):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOLANGCI_LINT_DIR) $(GOLANGCI_LINT_VERSION)

.PHONY: install-golangci-lint
install-golangci-lint: $(GOLANGCI_LINT_DIR)

.PHONY: golangci-lint
golangci-lint: install-golangci-lint
	$(GOLANGCI_LINT_DIR)/golangci-lint run
.PHONY: clean
clean:
	rm -rf build
