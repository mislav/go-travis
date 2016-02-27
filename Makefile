GO = go
NAMESPACE = github.com/mislav/go-travis

SOURCES = $(shell find . -name '*.go' -maxdepth 2)
MODULES = $(shell find . -name '*.go' -depth 2 | cut -d/ -f1-2 | sort -u)
GITMODULES = $(shell cat .gitmodules | awk '/path =/ {print $$(NF)}')

.PHONY: fmt fmt-all

$(GITMODULES:=/.git):
	git submodule update --init --recursive

bin/travis: $(filter-out %_test.go,$(SOURCES)) $(GITMODULES:=/.git)
	@mkdir -p bin
	GO15VENDOREXPERIMENT=1 $(GO) build -o $@ $(NAMESPACE)

fmt:
	$(GO) fmt main.go
	$(GO) fmt $(MODULES)

fmt-all:
	$(GO) fmt ./...
