.PHONY: all
all: lint test cover

.PHONY: test
test: setup
	go test -v ./...

.PHONY: cover
cover: setup
	mkdir -p coverage
	gocov test ./... | gocov-html > coverage/coverage.html

sources = $(shell find . -name '*.go' -not -path './vendor/*')
.PHONY: goimports
goimports: setup
	goimports -w $(sources)

.PHONY: lint
lint: setup
	gometalinter --disable=golint --enable=goimports --vendor ./...

.PHONY: install
install: setup
	go install

BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep
GOCOV := $(BIN_DIR)/gocov
GOCOV_HTML := $(BIN_DIR)/gocov-html

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

$(GOCOV):
	go get -u github.com/axw/gocov/gocov

$(GOCOV_HTML):
	go get -u gopkg.in/matm/v1/gocov-html

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

tools: $(GOIMPORTS) $(GOMETALINTER) $(GOCOV) $(GOCOV_HTML) $(DEP)

vendor: $(DEP)
	dep ensure

setup: tools vendor

updatedeps:
	dep ensure -update
