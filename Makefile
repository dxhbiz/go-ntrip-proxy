all: build


MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}bin
RELEASE_FILE := ${RELEASE_DIR}/ntrip-proxy
GO_PATH := $(shell go env | grep GOPATH | awk -F '"' '{print $$2}')

# Version
RELEASE?=v1.0.0

# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

# Build Flags
GO_LD_FLAGS= "-s -w -X github.com/dxhbiz/go-ntrip-proxy/pkg/version.RELEASE=${RELEASE} -X github.com/dxhbiz/go-ntrip-proxy/pkg/version.COMMIT=${GIT_COMMIT} -X github.com/dxhbiz/go-ntrip-proxy/pkg/version.REPO=${GIT_REPO_INFO}"

.PHONY: build
build:
	@echo "build ntrip-proxy"
	cd ${MKFILE_DIR} && \
	go build -v -trimpath -ldflags ${GO_LD_FLAGS} \
	-o ${RELEASE_FILE} ${MKFILE_DIR}cmd/ntrip-proxy

.PHONY: clean
clean:
	@echo "clean ntrip-proxy"
	cd ${MKFILE_DIR} && \
	rm -f ${RELEASE_FILE}