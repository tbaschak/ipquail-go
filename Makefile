APP := ipquail
PACKAGE := ./ipquail
BIN_DIR := bin
GO ?= go
GOFLAGS ?=

.PHONY: all build run test fmt vet clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP) $(PACKAGE)

run:
	$(GO) run $(GOFLAGS) $(PACKAGE)

test:
	$(GO) test $(GOFLAGS) ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet $(GOFLAGS) ./...

clean:
	$(GO) clean
	$(RM) $(BIN_DIR)/$(APP)
	@rmdir $(BIN_DIR) 2>/dev/null || true
