
SHELL := /bin/bash

PWD ?= $(shell pwd)

CORE_LAMBDA_SOURCES_DIR := ${PWD}/src

.PHONY: core-apply compile-core install-core install clone-adaptive-core-lambdas

include common.Makefile
include core.tf.Makefile

install: install-core

dev: ${CORE_TERRAFORM_SRC}/deploy-auto.log

test-core:
	cd src; \
	go test -timeout 1h -v

init:
	cd terraform;\
	./backends/init.sh
