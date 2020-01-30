.PHONY: backdate-feedback

${PWD}/bin/backdate: ${PWD}/migrations/backdate-feedback/*.go
	pushd ${PWD}/migrations/backdate-feedback;\
	go build -o ${PWD}/bin/backdate;\
	popd;

backdate-build: ${PWD}/bin/backdate

backdate-feedback: ${PWD}/bin/backdate
	export LOG_NAMESPACE="backdate" ;\
	export PLATFORM_ID="NO-PLATFORM" ;\
	${PWD}/bin/backdate
