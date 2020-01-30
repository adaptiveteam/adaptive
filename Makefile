
SHELL := /bin/bash

#COACHING_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-coaching-lambdas
#CORE_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-core-lambdas
#STRATEGY_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-strategy-lambdas
#USER_COMMUNITY_LAMBDA = $(ADAPTIVE_REPOS)/adaptive-user-community-lambdas

.DEFAULT_GOAL := all

.SHELLFLAGS = -ec

.ONESHELL:

.PHONY: help test core-apply core-lambdas install
LAMBDAS_SRC_DIR := $(shell cd ..; pwd)
PWD := $(shell pwd)
AMM_BIN := ${PWD}/bin/amm
AMM_VERSION := 2.0.4
AMM := ${AMM_BIN}-${AMM_VERSION}

#CORE_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-core-lambdas
#COACHING_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-coaching-lambdas
#STRATEGY_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-strategy-lambdas
#USER_COMMUNITY_LAMBDA_SOURCES_DIR := ${ADAPTIVE_REPOS}/adaptive-user-community-lambdas

include common.Makefile
include build.Makefile
include deploy.Makefile
#include core.main.Makefile

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

# all-git-pull:
# 	pushd $(ADAPTIVE_REPOS) ;\
# 	for d in ./*/ ; do (pushd "$$d" && pwd && git pull && popd); done ;\
# 	popd

# all-git-branch:
# 	pushd $(ADAPTIVE_REPOS) ;\
# 	for d in ./*/ ; do (pushd "$$d">>/dev/null && echo "branch:$$(git rev-parse --abbrev-ref HEAD) in $$d" && popd >> null); done ;\
# 	popd

# all-git-status:
# 	pushd $(ADAPTIVE_REPOS) ;\
# 	for d in ./*/ ; do (pushd "$$d">>/dev/null && echo "$$d status:$$(git s)" && popd >> /dev/null); done ;\
# 	popd

# all-echo-all:
# 	pushd $(ADAPTIVE_REPOS) ;\
# 	for d in ./*/ ; do (echo "dir:$$d"); done ;\
# 	popd

# backup-all:
# 	pushd $(ADAPTIVE_REPOS)/core-infra-terraform ;\
# 	make backup-all ;\
# 	popd

restore-all:
	pushd $(ADAPTIVE_REPOS)/core-infra-terraform ;\
	make restore-all ;\
	popd

all:
	echo "all"

test-with-localstack:
	docker-compose up -d ;\
	go test ${TEST_OPS} -v ./...  -coverprofile=cover.out ;\
	docker-compose down

test:
	go test -v ${TEST_OPS} ./...  -coverprofile=cover.out
test-short:
	go test -short -v ${TEST_OPS} ./...  -coverprofile=cover.out
clean:
	go clean
deps:
	go build -v ./...
upgrade:
	go get -u
coverage: test
	go tool cover -html=cover.out
