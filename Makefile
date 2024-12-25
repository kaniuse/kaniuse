.DEFAULT_GOAL:=help

.PHONY: image-data-scraper
image-data-scraper: ## Build image ghcr.io/kaniuse/data-scraper:latest
	docker build -t ghcr.io/kaniuse/data-scraper:latest ./data-scraper

.PHONY: data-scraper
data-scraper: ## Build binary executable for data-scraper
	cd data-scraper && go build -o ./data-scraper ./cmd/data-scraper

.PHONY: data
data: data-api-lifecycle data-kinds data-fields ## Update data JSON files

.PHONY: data-api-lifecycle
data-api-lifecycle:
		cd data-scraper && \
		go run ./cmd/data-scraper api-lifecycle -w ../public/data/gvk_api_lifecycle.json

.PHONY: data-kinds
data-kinds:
	cd data-scraper && \
		go run ./cmd/data-scraper kinds -w ../public/data/kinds.json

.PHONY: data-fields
data-fields:
	cd data-scraper && \
		go run ./cmd/data-scraper fields -w ../public/data/fields.json

# The help will print out all targets with their descriptions organized bellow their categories. The categories are represented by `##@` and the target descriptions by `##`.
# The awk commands is responsible to read the entire set of makefiles included in this invocation, looking for lines of the file as xyz: ## something, and then pretty-format the target and help. Then, if there's a line with ##@ something, that gets pretty-printed as a category.
# More info over the usage of ANSI control characters for terminal formatting: https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info over awk command: http://linuxcommand.org/lc3_adv_awk.php
.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
