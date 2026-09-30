PLUGIN_ID ?= oracle-grafana-datasource
VERSION ?= $(shell cat package.json | grep '"version"' | head -n1 | cut -d'"' -f4)
DIST_DIR := dist

.PHONY: all clean test build build-backend build-all-backends dist sign package

all: test build-backend

test:
	@echo "==> Running Go test suite..."
	go test -v ./pkg/...

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(DIST_DIR) *.tar.gz gpx_*

build-backend:
	@echo "==> Building native Linux backend (CGO_ENABLED=0)..."
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_linux_amd64 ./pkg

build-all-backends:
	@echo "==> Cross-compiling backends for all supported Grafana targets..."
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_linux_amd64 ./pkg
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_linux_arm64 ./pkg
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_darwin_amd64 ./pkg
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_darwin_arm64 ./pkg
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(DIST_DIR)/gpx_oracle_grafana_windows_amd64.exe ./pkg

dist: build-all-backends
	@echo "==> Packaging distribution bundle..."
	cp -r src/img $(DIST_DIR)/ 2>/dev/null || true
	cp src/plugin.json $(DIST_DIR)/
	cp README.md LICENSE CHANGELOG.md $(DIST_DIR)/ 2>/dev/null || true
	@if [ -f src/module.js ]; then cp src/module.js* $(DIST_DIR)/; fi
	@echo "==> Dist directory ready."

sign: dist
	@echo "==> Signing plugin with @grafana/sign-plugin..."
	@if [ -z "$$GRAFANA_ACCESS_POLICY_TOKEN" ]; then \
		echo "ERROR: GRAFANA_ACCESS_POLICY_TOKEN is not set."; \
		echo "To sign: export GRAFANA_ACCESS_POLICY_TOKEN='<your-token>' && make sign"; \
		exit 1; \
	fi
	npx --yes @grafana/sign-plugin@latest --rootUrls $${GRAFANA_ROOT_URLS:-http://localhost:3000}

package: dist
	@echo "==> Creating release archive..."
	tar -czf $(PLUGIN_ID)-$(VERSION).tar.gz -C $(DIST_DIR) .
	@echo "Created $(PLUGIN_ID)-$(VERSION).tar.gz"
