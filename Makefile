GO ?= go
GOFMT ?= gofmt
STATICCHECK_VERSION := v0.8.1

.PHONY: check fmt-check fmt vet lint test

check: fmt-check vet lint test

fmt-check:
	@set -e; files="$$($(GOFMT) -l cmd internal)"; \
	if [ -n "$$files" ]; then \
		printf '%s\n' "Run make fmt to format:" "$$files"; exit 1; \
	fi

fmt:
	$(GOFMT) -w cmd internal

vet:
	$(GO) vet -mod=readonly -tags=nomsgpack ./...

lint:
	$(GO) run honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION) -tags=nomsgpack ./...

test:
	$(GO) test -mod=readonly -tags=nomsgpack -coverpkg=./... -coverprofile=coverage.out ./...
