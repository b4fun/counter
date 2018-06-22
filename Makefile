.PHONY: docker-image
docker-image: docker-image-counter

.PHONY: docker-image-counter
docker-image-counter:
	docker build -t b4fun/counter:latest -f dockerfile/counter/Dockerfile .
