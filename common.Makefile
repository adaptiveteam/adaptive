#######################################
# Common library for building lambdas #
#######################################

# go-sources(goSourcesPath) - Function that returns all *.go files and go.mod and go.sum 
#                             from the given directory.
go-sources = $(wildcard $(1)/*.go) $(wildcard $(1)/**/*.go)
# $(1)/go.mod $(1)/go.sum

tf-sources = $(1)/*.tf $(1)/*.tfvars $(1)/backends/*.tfbackend
tf-sources-only = $(1)/*.tf $(1)/backends/*.tfbackend

# build-go(srcDir,lambdaName,lambdaDir) - function that builds go source folder.
#                      As a result it produces $(lambdaDir)/$lambdaName.zip 
#                      with the executable and resources if available.
define build-go =
	set -e; \
	pushd $(1); \
	pwd;\
	mkdir -p zip-contents;\
	GOOS=linux GOARCH=amd64 go build -o zip-contents/$(2); \
	echo Built $(2) in $(1); \
	[ -d main/resources ] && cp -a main/resources/* zip-contents/
	mkdir -p $(3); \
	pushd zip-contents; \
	zip -r $(3)/$(2).zip *; \
	popd; \
	rm -r zip-contents; \
	popd
endef


# test-go(srcDir,logFilename) - function that tests go source folder.
define test-go =
	set -e; \
	pushd $(1); \
	pwd;\
	GOOS=linux GOARCH=amd64 go test -v > $(2); \
	popd
endef
