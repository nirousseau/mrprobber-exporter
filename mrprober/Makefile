#!/usr/bin/make -sf
BINARY_NAME := mrprober

BUILDID = $(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')
VERSION = $(shell cat mrprober.spec | grep Version | cut -d ':' -f2 | xargs)
CWD = $(shell pwd)

#DOCKER_USER := $(shell whoami)
DOCKER_REGISTRY := docker.io
DOCKER_REPO := $(DOCKER_REGISTRY)/miscellaneous
DOCKER_IMAGE := $(BINARY_NAME):$(VERSION)

.DEFAULT := $(BINARY_NAME)
.PHONY: install docker docker-run docker-push dep clean

# Linux binary
$(BINARY_NAME):
	# https://fedoraproject.org/wiki/PackagingDrafts/Go#Build_ID
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags "-X main.VERSION=$(VERSION) -s -w -B 0x$(BUILDID)" -v

# Build docker image from scratch
docker: $(BINARY_NAME)
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker
	docker run --rm -v ${CWD}/conf/samples/mrprober.defaults.yaml:/etc/mrprober/mrprober.yaml -p 9100:9100 -it $(DOCKER_IMAGE)

docker-push: docker
	@docker login $(DOCKER_REGISTRY) -u "$(DOCKER_USER)" -p "$(DOCKER_PASSWD)"
	docker tag $(DOCKER_IMAGE) $(DOCKER_REPO)/$(DOCKER_IMAGE)
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE)

test:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download
	go mod tidy

clean:
	go clean
	docker system prune -f
	rm -f mrprober
	rm -f mrprober.exe
