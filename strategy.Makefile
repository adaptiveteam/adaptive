STRATEGY_TERRAFORM := ${ADAPTIVE_REPOS}/adaptive-strategy-infra-terraform/terraform

STRATEGY_LAMBDAS := ${STRATEGY_TERRAFORM}/lambdas

.PHONY: strategy-apply compile-strategy install-strategy

# Compiles go project. This target depends on all go sources - *.go and go.*
#                      As a result it produces `.binary` file.
#                      Note that makefile doesn't support two % in target name. Hence we have to introduce a new filename.
${ADAPTIVE_REPOS}/adaptive-strategy-lambdas/%/.binary: $(call go-sources,${ADAPTIVE_REPOS}/adaptive-strategy-lambdas/%/src)
	pushd $(@D)/src
	go build -o $@
	popd

# Packages binary of lambda.
#           $(notdir $(<D)) gets directory (`D`) of the first (`$<`) dependency. 
#           In our case it's path to `.binary`. And from that directory takes just filename.
#           This gives us exactly lambda name.
#           Then we create lambda executable by copying `.binary` and 
#           zip it.
${STRATEGY_LAMBDAS}/%.zip: ${ADAPTIVE_REPOS}/adaptive-strategy-lambdas/%/.binary
	pushd $(<D)
	cp .binary $(notdir $(<D))
	zip $@ $(notdir $(<D))
	rm $(notdir $(<D))
	popd

# List of all zips with lambdas.
STRATEGY_LAMBDA_ZIPS := ${STRATEGY_LAMBDAS}/strategy-slack-message-processor-lambda-go.zip

compile-strategy: $(STRATEGY_LAMBDA_ZIPS)
	echo Coaching-lambdas zip files are up to date
	ls -l ${STRATEGY_LAMBDAS}
	
# deploy.log is updated during terraform run. If all zips are older, then 
# terraform is not invoked. Otherwise it is triggered and it has own means to 
# check for changes.
${STRATEGY_TERRAFORM}/deploy.log: $(STRATEGY_LAMBDA_ZIPS)
	pushd $(@D)
	terraform apply -auto-approve > deploy.log
	popd

# install-strategy is the top level goal for installing strategy lambdas
install-strategy: ${STRATEGY_TERRAFORM}/deploy.log
