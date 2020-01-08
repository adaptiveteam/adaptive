
SHELL := /bin/bash

ADAPTIVE_REPOS ?= $(shell cd ../ ; pwd)

COACHING_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-coaching-lambdas
CORE_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-core-lambdas
STRATEGY_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-strategy-lambdas
USER_COMMUNITY_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-user-community-lambdas

.DEFAULT_GOAL := all

.SHELLFLAGS = -ec

.ONESHELL:

.PHONY: help test core-apply core-lambdas install
ADAPTIVE_REPOS := $(shell cd ..; pwd)
PWD := $(shell pwd)
AMM_BIN := ${PWD}/bin/amm
AMM_VERSION := 1.8.2
AMM := ${AMM_BIN}-${AMM_VERSION}
CORE_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-core-lambdas
COACHING_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-coaching-lambdas
STRATEGY_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-strategy-lambdas
USER_COMMUNITY_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-user-community-lambdas

# include ../adaptive-core-lambdas/common.Makefile
# include ../adaptive-core-lambdas/core.Makefile
# include ../core-infra-terraform/core.tf.Makefile
# include ../adaptive-coaching-lambdas/coaching.Makefile
# include ../adaptive-coaching-infra-terraform/coaching.tf.Makefile
# include ../adaptive-strategy-lambdas/strategy.Makefile
# include ../adaptive-strategy-infra-terraform/strategy.tf.Makefile
# include ../adaptive-user-community-lambdas/user-community.Makefile
# include ../adaptive-user-community-infra-terraform/user-community.tf.Makefile

# Check preconditions

ifndef ADAPTIVE_REPOS
echo "ADAPTIVE_REPOS=${ADAPTIVE_REPOS}"
$(error ADAPTIVE_REPOS is not defined)
endif

# go-sources(goSourcesPath) - Function that returns all *.go files and go.mod and go.sum 
#                             from the given directory.
go-sources = $(wildcard $(1)/*.go) $(1)/go.mod $(1)/go.sum

# install is the top level goal for installing all lambdas
# install: install-core install-coaching install-user-community install-strategy

# compile: compile-core compile-coaching compile-user-community compile-strategy

help: ## Prints description of all goals
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-ammonite: ${AMM}

${AMM}:
	mkdir -p ${PWD}/bin
	rm -f ${AMM_BIN}
	sh -c '(echo "#!/usr/bin/env sh" && curl -L https://github.com/lihaoyi/Ammonite/releases/download/${AMM_VERSION}/2.13-${AMM_VERSION}) > ${AMM} && chmod +x ${AMM} && ln -s ${AMM} ${AMM_BIN}' 

generate: ${AMM}
	${AMM_BIN} scripts/Gen.sc

generate-dry-run: ${AMM}
	${AMM_BIN} scripts/Gen.sc --dry-run 

all-git-pull:
	pushd $(ADAPTIVE_REPOS) ;\
	for d in ./*/ ; do (pushd "$$d" && pwd && git pull && popd); done ;\
	popd

all-git-branch:
	pushd $(ADAPTIVE_REPOS) ;\
	for d in ./*/ ; do (pushd "$$d">>/dev/null && echo "branch:$$(git rev-parse --abbrev-ref HEAD) in $$d" && popd >> null); done ;\
	popd

all-git-status:
	pushd $(ADAPTIVE_REPOS) ;\
	for d in ./*/ ; do (pushd "$$d">>/dev/null && echo "$$d status:$$(git s)" && popd >> /dev/null); done ;\
	popd

all-echo-all:
	pushd $(ADAPTIVE_REPOS) ;\
	for d in ./*/ ; do (echo "dir:$$d"); done ;\
	popd

backup-all:
	pushd $(ADAPTIVE_REPOS)/core-infra-terraform ;\
	make backup-all ;\
	popd

restore-all:
	pushd $(ADAPTIVE_REPOS)/core-infra-terraform ;\
	make restore-all ;\
	popd

all: test

test-with-localstack:
	docker-compose up -d ;\
	go test -v ./...  -coverprofile=cover.out ;\
	docker-compose down

test:
	go test -v ./...  -coverprofile=cover.out
clean:
	go clean
deps:
	go build -v ./...
upgrade:
	go get -u
coverage: test
	go tool cover -html=cover.out
