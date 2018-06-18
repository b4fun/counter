.PHONY: docker-image
docker-image:
	docker build -t b4fun/counter:latest .
