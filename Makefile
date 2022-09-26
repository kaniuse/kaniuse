.PHONY: image-data-scraper
image-data-scraper:
	docker build -t ghcr.io/kaniuse/data-scraper:latest ./data-scraper
