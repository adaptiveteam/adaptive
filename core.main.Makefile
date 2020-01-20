
SHELL := /bin/bash

PWD ?= $(shell pwd)

CORE_LAMBDA_SOURCES_DIR := ${PWD}/src

.PHONY: core-apply compile-core install-core install clone-adaptive-core-lambdas

include common.Makefile
include core.tf.Makefile

./dynamodump/dynamodump.py:
	git clone git@github.com:bchew/dynamodump.git;\
	pushd dynamodump;\
	pip install -r requirements.txt;\
	pip install flake8;\
	popd

backup-all: ./dynamodump/dynamodump.py
	python dynamodump/dynamodump.py -m backup  -r ${AWS_REGION} -s ${ADAPTIVE_CLIENT_ID}*

restore-all: ./dynamodump/dynamodump.py
	python dynamodump/dynamodump.py -m restore -r ${AWS_REGION} --dataOnly -s ${ADAPTIVE_CLIENT_ID}*

install: install-core

dev: ${CORE_TERRAFORM_SRC}/deploy-auto.log

test-core:
	cd src; \
	go test -timeout 1h -v

init:
	cd terraform;\
	./backends/init.sh

restore-table-user-objective: ./dynamodump/dynamodump.py
	./rename-backup.sh
	python dynamodump/dynamodump.py --skipThroughputUpdate -m restore -r ${AWS_REGION} --dataOnly -s ${ADAPTIVE_CLIENT_ID}_user_objective

rename-resource-user-objective:
	cd terraform;\
	terraform state mv 'aws_dynamodb_table.user_objectives' 'aws_dynamodb_table.user_objective_dynamodb_table'
