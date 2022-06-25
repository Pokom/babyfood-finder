DOCKER_IMAGE_TAG ?= $(subst /,-,$(shell git rev-parse --short HEAD))
build-pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o babyfood-finder


build-docker:
	docker build -t babyfood-finder:latest .
	docker tag babyfood-finder:latest babyfood-finder:$(DOCKER_IMAGE_TAG)
	docker tag babyfood-finder:latest us-east4-docker.pkg.dev/sandbox-307502/docker/babyfood-finder:$(DOCKER_IMAGE_TAG)
	docker tag babyfood-finder:latest us-east4-docker.pkg.dev/sandbox-307502/docker/babyfood-finder:latest

release-gcr: build-docker
	docker push us-east4-docker.pkg.dev/sandbox-307502/docker/babyfood-finder:$(DOCKER_IMAGE_TAG)
	docker push us-east4-docker.pkg.dev/sandbox-307502/docker/babyfood-finder:latest

run-local: build-docker
	docker run --env-file=.env -it babyfood-finder:latest -to=+17324067063

