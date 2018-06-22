.PHONY: docker-image
docker-image: docker-image-counter docker-image-redis

.PHONY: docker-image-counter
docker-image-counter:
	docker build -t b4fun/counter:latest -f dockerfile/counter/Dockerfile .

.PHONY: docker-image-redis
docker-image-redis:
	docker build -t b4fun/counter-redis-slave:latest dockerfile/redis-slave
