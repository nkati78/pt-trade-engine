DOCKER_IMAGE=paper-thesis/trade-engine

# create migration timestamp
# Usage: make migration-date
migration-date:
	@echo $(shell date +%Y%m%d%H%M%S)

# run docker container
# Usage: make run
run:
	docker-compose up

# build docker container
# Usage: make build
build:
	docker build --platform linux/amd64 -t $(DOCKER_IMAGE) .

build-mac:
	docker build -t $(DOCKER_IMAGE) .