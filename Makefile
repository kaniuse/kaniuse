.PHONY: image-data-scraper
image-data-scraper:
	docker build -t ghcr.io/kaniuse/data-scraper:latest ./data-scraper

.PHONY: data
data:
	cd data-scraper && go run . > ../server/data/gvk_api_lifecycle.json
