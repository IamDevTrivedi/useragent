SHELL := /bin/sh

GO ?= go
GOLANGCI_LINT ?= golangci-lint
COVERAGE_FILE ?= coverage.out
FUZZTIME ?= 30s
BENCH_PATTERN ?= .
VERSION ?= v0.1.0

.PHONY: all check verify fmt fmt-check test race vet lint cover bench fuzz fuzz-all cli release-assets clean

all: check

check: fmt-check vet test race lint cover

verify: check

fmt:
	$(GO) fmt ./...

fmt-check:
	@out="$$(gofmt -l .)"; if [ -n "$$out" ]; then \
		printf 'gofmt required for:\n%s\n' "$$out"; exit 1; \
	fi

test:
	$(GO) test ./...

race:
	$(GO) test -race ./...

vet:
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run ./...

cover:
	$(GO) test -coverprofile=$(COVERAGE_FILE) ./...
	@total="$$( $(GO) tool cover -func=$(COVERAGE_FILE) | tail -n 1 | awk '{print $$3}' | tr -d '%' )"; \
	printf 'total coverage: %s%%\n' "$$total"; \
	awk -v total="$$total" 'BEGIN { if (total < 90) { printf "coverage below 90%%: %s%%\n", total; exit 1 } }'

bench:
	$(GO) test -run '^$$' -bench '$(BENCH_PATTERN)' -benchmem ./...

fuzz:
	$(GO) test -run '^$$' -fuzz '^FuzzParse$$' -fuzztime '$(FUZZTIME)' .

fuzz-all:
	@set -e; \
	for target in FuzzParse FuzzTokenizer FuzzParseVersion FuzzParseBrandList; do \
		$(GO) test -run '^$$' -fuzz "^$$target$$" -fuzztime '$(FUZZTIME)' .; \
	done; \
	$(GO) test -run '^$$' -fuzz '^FuzzParseRequest$$' -fuzztime '$(FUZZTIME)' ./uahttp

cli:
	@mkdir -p bin
	$(GO) build -trimpath -buildvcs=false -o bin/ua ./cmd/ua

release-assets:
	./scripts/release-assets.sh $(VERSION)

clean:
	rm -rf bin dist $(COVERAGE_FILE) *.prof *.test
