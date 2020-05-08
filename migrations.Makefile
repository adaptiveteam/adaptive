.PHONY: backdate-feedback

${PWD}/bin/backdate: ${PWD}/migrations/backdate-feedback/*.go
	pushd ${PWD}/migrations/backdate-feedback;\
	go build -o ${PWD}/bin/backdate;\
	popd;

backdate-build: ${PWD}/bin/backdate

backdate-feedback: ${PWD}/bin/backdate
	export LOG_NAMESPACE="backdate" ;\
#	export PLATFORM_ID="NO-PLATFORM" ;\
	${PWD}/bin/backdate

rename-user-engagement: backup-all 
	mv ./dump/${ADAPTIVE_CLIENT_ID}_adaptive_users_engagements ./dump/${ADAPTIVE_CLIENT_ID}_user_engagement ;\
	pushd infra/core/terraform ;\
	terraform apply -target=aws_dynamodb_table.adaptive_user_engagements_dynamo_table ;\
	popd ;\
	$(call restore-table,${ADAPTIVE_CLIENT_ID}_user_engagement)

update-issues.bin:
	pushd migrations/update-issues ;\
	go build ;\
	popd

update-issues: backup-all update-issues.bin
	export LOG_NAMESPACE="MAKE update-issues" ;\
	export CLIENT_ID=${ADAPTIVE_CLIENT_ID} ;\
	./migrations/update-issues/update-issues -all ;\
