GO = go

SOURCES = $(shell find . -name '*.go' -maxdepth 2)
MODULES = $(shell find . -name '*.go' -depth 2 | cut -d/ -f1-2 | sort -u)

.PHONY: fmt fmt-all

bin/travis: $(filter-out %_test.go,$(SOURCES))
	mkdir -p bin
	GO15VENDOREXPERIMENT=1 $(GO) build -o bin/travis github.com/mislav/go-travis

fmt:
	$(GO) fmt main.go
	$(GO) fmt $(MODULES)

fmt-all:
	$(GO) fmt ./...
