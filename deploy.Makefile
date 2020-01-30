CORE_TERRAFORM_SRC := ${PWD}/infra/core/terraform

.PHONY: core-deploy terraform-validate terraform-apply-auto-approve

# deploy.log is updated during terraform run. If all zips are older, then 
# terraform is not invoked. Otherwise it is triggered and it has own means to 
# check for changes.
${CORE_TERRAFORM_SRC}/deploy.log: adaptive-build \
	$(call tf-sources-only,${CORE_TERRAFORM_SRC})
	pushd $(@D);\
	pwd;\
	time terraform apply;\
	date > deploy.log;\
	date ;\
	popd

# deploy-auto.log is updated during terraform run. If all zips are older, then 
# terraform is not invoked. Otherwise it is triggered and it has own means to 
# check for changes.
${CORE_TERRAFORM_SRC}/deploy-auto.log: $(CORE_LAMBDA_BINS) \
	$(call tf-sources-only,${CORE_TERRAFORM_SRC})
	pushd $(@D);\
	pwd;\
	terraform apply -auto-approve;\
	date > deploy-auto.log;\
	popd

# core-deploy is the top level goal for installing core lambdas
core-deploy: ${CORE_TERRAFORM_SRC}/deploy.log

terraform-validate: $(CORE_LAMBDA_BINS)
	cd terraform;\
	terraform validate

terraform-apply-auto-approve: ${CORE_TERRAFORM_SRC}/deploy-auto.log

core-init:
	pushd ${CORE_TERRAFORM_SRC};\
	./backends/init.sh;\
	popd
