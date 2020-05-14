# force to use bash
SHELL = /bin/bash

WRKDIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
COMMIT_VER := $(shell git rev-parse HEAD)
DOCKER_IMAGE_NAME = gitlab.connectwisedev.com/platform/platform-tasking-service
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
DOCKER_IMAGE_VERSION = master
SRCDIR = src
TRANSLATIONPATH=./src/translation

bindata:
	${GOPATH}/bin/go-bindata -pkg genlang -o $(TRANSLATIONPATH)/genlang/bindata_gen.go $(TRANSLATIONPATH)/languages/

docker-image-build: $(OUT_PATH)
	@ (cd $(SRCDIR) && glide install)
	@ make bindata
	@ cd $(WRKDIR)

	# Building image
	@ docker build --label version=$(DOCKER_IMAGE_VERSION) --label commit=$(COMMIT_VER) --label build_date=$(BUILD_DATE) -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VERSION) $(WRKDIR)
.PHONY: docker-image-build