USER_COMMUNITY_TERRAFORM := ${ADAPTIVE_REPOS}/adaptive-user-community-infra-terraform/terraform

USER_COMMUNITY_LAMBDAS := ${USER_COMMUNITY_TERRAFORM}/lambdas

.PHONY: user-community-apply compile-user-community install-user-community

# Compiles go project. This target depends on all go sources - *.go and go.*
#                      As a result it produces `.binary` file.
#                      Note that makefile doesn't support two % in target name. Hence we have to introduce a new filename.
${ADAPTIVE_REPOS}/adaptive-user-community-lambdas/%/.binary: $(call go-sources,${ADAPTIVE_REPOS}/adaptive-user-community-lambdas/%/src)
	pushd $(@D)/src
	go build -o $@
	popd

# Packages binary of lambda.
#           $(notdir $(<D)) gets directory (`D`) of the first (`$<`) dependency. 
#           In our case it's path to `.binary`. And from that directory takes just filename.
#           This gives us exactly lambda name.
#           Then we create lambda executable by copying `.binary` and 
#           zip it.
${USER_COMMUNITY_LAMBDAS}/%.zip: ${ADAPTIVE_REPOS}/adaptive-user-community-lambdas/%/.binary
	pushd $(<D)
	cp .binary $(notdir $(<D))
	zip $@ $(notdir $(<D))
	rm $(notdir $(<D))
	popd

# List of all zips with lambdas.
USER_COMMUNITY_LAMBDA_ZIPS := \
			${USER_COMMUNITY_LAMBDAS}/adaptive-community-slack-message-processor-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/adaptive-community-watcher-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/adaptive-hello-world-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/adaptive-holidays-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/adaptive-user-objectives-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/adaptive-values-lambda-go.zip \
			${USER_COMMUNITY_LAMBDAS}/user-engagement-scheduler-lambda-go.zip


compile-user-community: $(USER_COMMUNITY_LAMBDA_ZIPS)
	echo Coaching-lambdas zip files are up to date
	ls -l ${USER_COMMUNITY_LAMBDAS}
	
# deploy.log is updated during terraform run. If all zips are older, then 
# terraform is not invoked. Otherwise it is triggered and it has own means to 
# check for changes.
${USER_COMMUNITY_TERRAFORM}/deploy.log: $(USER_COMMUNITY_LAMBDA_ZIPS)
	pushd $(@D)
	terraform apply -auto-approve > deploy.log
	popd

# install-user-community is the top level goal for installing user-community lambdas
install-user-community: ${USER_COMMUNITY_TERRAFORM}/deploy.log
