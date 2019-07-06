ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SHELL := /bin/bash

SOURCES = $(shell find $(ROOT_DIR) -name "*.go" -print)

GOOS ?= linux
GOARCH ?= amd64
GOPATH ?= $(shell pwd)

export GO111MODULE = on

.PHONY: test check .coverprofile

default: all

all: build check

check: checkfmt test

.coverprofile:
	go test -coverprofile .coverprofile

cover: .coverprofile
	go tool cover -func .coverprofile

showcover: .coverprofile
	go tool cover -html .coverprofile

build:
	go build ./...

test:
	go test -race ./...

checkfmt:
	@[ -z $$(gofmt -l $(SOURCES)) ] || (echo "Sources not formatted correctly. Fix by running: make fmt" && false)

fmt: $(SOURCES)
	gofmt -s -w $(SOURCES)
