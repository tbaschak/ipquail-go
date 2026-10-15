APP := ipquail
PACKAGE := ./ipquail
BIN_DIR := bin
GO ?= go
GOFLAGS ?=
NVM_DIR ?= $(HOME)/.nvm
NVM := . "$(NVM_DIR)/nvm.sh" &&
PUBLIC_DIR := ipquail/public

.PHONY: all build run test fmt vet verify web web-install web-deploy clean

all: build

build: web-deploy
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP) $(PACKAGE)

run: web-deploy
	$(GO) run $(GOFLAGS) $(PACKAGE)

web: web-deploy

web-install: node_modules/.package-lock.json

node_modules/.package-lock.json: package-lock.json
	$(NVM) nvm exec --silent npm ci

web-deploy: web-install
	$(NVM) nvm exec --silent npm run build

test:
	$(GO) test $(GOFLAGS) ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet $(GOFLAGS) ./...

verify:
	@unformatted="$$(gofmt -l $$(find . -type f -name '*.go' -not -path './vendor/*'))"; \
	if [ -n "$$unformatted" ]; then \
		echo "The following Go files need formatting:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	$(GO) mod tidy -diff

clean:
	$(GO) clean
	$(RM) $(BIN_DIR)/$(APP)
	$(RM) $(PUBLIC_DIR)/css/bootstrap*.css
	$(RM) $(PUBLIC_DIR)/css/bootstrap*.css.map
	$(RM) $(PUBLIC_DIR)/js/vendor/jquery*.js
	$(RM) $(PUBLIC_DIR)/js/vendor/jquery*.map
	@rmdir $(BIN_DIR) 2>/dev/null || true
