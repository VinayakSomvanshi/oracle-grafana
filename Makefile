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
	rm -rf $(DIST_DIR) *.tar.gz *.zip gpx_*

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
	cp -r dashboards $(DIST_DIR)/ 2>/dev/null || true
	@if [ -f src/module.js ]; then cp src/module.js* $(DIST_DIR)/; fi
	@echo "==> Dist directory ready."

sign:
	@echo "==> Signing plugin with @grafana/sign-plugin..."
	@if [ -z "$$GRAFANA_ACCESS_POLICY_TOKEN" ]; then \
		echo "ERROR: GRAFANA_ACCESS_POLICY_TOKEN is not set."; \
		echo "To sign: export GRAFANA_ACCESS_POLICY_TOKEN='<your-token>' && make sign"; \
		exit 1; \
	fi
	@URLS="$${GRAFANA_ROOT_URLS}"; \
	if [ -z "$$URLS" ]; then \
		URLS="http://localhost:3000/,http://*:3000/,https://*:3000/,http://localhost:8080/,http://*:8080/,https://*:8080/,http://localhost:8443/,http://*:8443/,https://*:8443/,http://localhost:9000/,http://*:9000/,https://*:9000/,http://localhost:3001/,http://*:3001/,https://*:3001/,http://*,https://*"; \
	elif ! echo "$$URLS" | grep -q "3000"; then \
		URLS="http://localhost:3000/,http://*:3000/,https://*:3000/,$$URLS"; \
	fi; \
	echo "==> Signing with root URLs: $$URLS"; \
	npx --yes @grafana/sign-plugin@latest --rootUrls "$$URLS"

package:
	@echo "==> Creating release archives..."
	tar -czf $(PLUGIN_ID)-$(VERSION).tar.gz -C $(DIST_DIR) .
	rm -rf /tmp/$(PLUGIN_ID) && mkdir -p /tmp/$(PLUGIN_ID) && cp -r $(DIST_DIR)/* /tmp/$(PLUGIN_ID)/
	cd /tmp && zip -q -r $(CURDIR)/$(PLUGIN_ID)-$(VERSION).zip $(PLUGIN_ID) && rm -rf /tmp/$(PLUGIN_ID)
	@echo "Created $(PLUGIN_ID)-$(VERSION).tar.gz and $(PLUGIN_ID)-$(VERSION).zip"
