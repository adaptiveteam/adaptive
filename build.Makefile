include lambdas.Makefile

.PHONY: core-build core-clean adaptive-build adaptive-clean
# List of all lambda binaries.
CORE_LAMBDA_BINS := $(foreach lambda,$(CORE_LAMBDA_NAMES),${PWD}/bin/$(strip ${lambda}))
ADAPTIVE_BIN := target/adaptive

${PWD}/target/%: infra/core/src/%/*.go lambdas/%/*.go
	pushd ${PWD}/infra/core/src/$*;\
	GOOS=linux GOARCH=amd64 go build -o ${PWD}/target/$*;\
	cp -r ${PWD}/resources/* ${PWD}/target ;\
	popd

core-build: $(CORE_LAMBDA_BINS)

# core-clean is used when we want to rebuild all lambdas. In particular,
# when some of the root libraries have been changed.
core-clean:
	rm -r ${PWD}/target/

# We do not need zip files, because they are produced by terraform itself
# .PHONY: core-zips
# 
# # List of all lambda zip-files.
# CORE_LAMBDA_ZIPS := $(foreach lambda,$(CORE_LAMBDA_NAMES),${PWD}/bin/$(strip ${lambda}).zip)
# 
# ${PWD}/bin/%.zip: $(call go-sources,${PWD}/infra/core/src/%)
# 	$(call build-go,${PWD}/infra/core/src/$*,$*,${PWD}/bin)
# 
# core-zips: $(CORE_LAMBDA_ZIPS)

${ADAPTIVE_BIN}: go.mod go.sum $(shell find . -name "*.go" -print)
	GOOS=linux GOARCH=amd64 go build -o ${ADAPTIVE_BIN} ;\
	cp -r ${PWD}/resources/* ${PWD}/target

adaptive-build: ${ADAPTIVE_BIN}

adaptive-clean:
	rm ${ADAPTIVE_BIN}
