
SHELL := /bin/bash


common/test.log: common/*.go
	pushd common; go test -v; date > test.log; popd


user/test.log: user/*.go
	pushd user; go test -v; date > test.log; popd


coachingRelationship/test.log: coachingRelationship/*.go
	pushd coachingRelationship; go test -v; date > test.log; popd


userFeedback/test.log: userFeedback/*.go
	pushd userFeedback; go test -v; date > test.log; popd


userEngagement/test.log: userEngagement/*.go
	pushd userEngagement; go test -v; date > test.log; popd


adaptiveCommunityUser/test.log: adaptiveCommunityUser/*.go
	pushd adaptiveCommunityUser; go test -v; date > test.log; popd


adaptiveCommunity/test.log: adaptiveCommunity/*.go
	pushd adaptiveCommunity; go test -v; date > test.log; popd


userAttribute/test.log: userAttribute/*.go
	pushd userAttribute; go test -v; date > test.log; popd


adHocHoliday/test.log: adHocHoliday/*.go
	pushd adHocHoliday; go test -v; date > test.log; popd


userObjective/test.log: userObjective/*.go
	pushd userObjective; go test -v; date > test.log; popd


userObjectiveProgress/test.log: userObjectiveProgress/*.go
	pushd userObjectiveProgress; go test -v; date > test.log; popd


adaptiveValue/test.log: adaptiveValue/*.go
	pushd adaptiveValue; go test -v; date > test.log; popd


clientPlatformToken/test.log: clientPlatformToken/*.go
	pushd clientPlatformToken; go test -v; date > test.log; popd


strategyObjective/test.log: strategyObjective/*.go
	pushd strategyObjective; go test -v; date > test.log; popd


strategyInitiative/test.log: strategyInitiative/*.go
	pushd strategyInitiative; go test -v; date > test.log; popd


visionMission/test.log: visionMission/*.go
	pushd visionMission; go test -v; date > test.log; popd


strategyCommunity/test.log: strategyCommunity/*.go
	pushd strategyCommunity; go test -v; date > test.log; popd


capabilityCommunity/test.log: capabilityCommunity/*.go
	pushd capabilityCommunity; go test -v; date > test.log; popd


strategyInitiativeCommunity/test.log: strategyInitiativeCommunity/*.go
	pushd strategyInitiativeCommunity; go test -v; date > test.log; popd


dialogEntry/test.log: dialogEntry/*.go
	pushd dialogEntry; go test -v; date > test.log; popd


contextAliasEntry/test.log: contextAliasEntry/*.go
	pushd contextAliasEntry; go test -v; date > test.log; popd


objectiveTypeDictionary/test.log: objectiveTypeDictionary/*.go
	pushd objectiveTypeDictionary; go test -v; date > test.log; popd


postponedEvent/test.log: postponedEvent/*.go
	pushd postponedEvent; go test -v; date > test.log; popd


slackTeam/test.log: slackTeam/*.go
	pushd slackTeam; go test -v; date > test.log; popd


community/test.log: community/*.go
	pushd community; go test -v; date > test.log; popd


channelMember/test.log: channelMember/*.go
	pushd channelMember; go test -v; date > test.log; popd


migration/test.log: migration/*.go
	pushd migration; go test -v; date > test.log; popd

test-all: common/test.log user/test.log coachingRelationship/test.log userFeedback/test.log userEngagement/test.log adaptiveCommunityUser/test.log adaptiveCommunity/test.log userAttribute/test.log adHocHoliday/test.log userObjective/test.log userObjectiveProgress/test.log adaptiveValue/test.log clientPlatformToken/test.log strategyObjective/test.log strategyInitiative/test.log visionMission/test.log strategyCommunity/test.log capabilityCommunity/test.log strategyInitiativeCommunity/test.log dialogEntry/test.log contextAliasEntry/test.log objectiveTypeDictionary/test.log postponedEvent/test.log slackTeam/test.log community/test.log channelMember/test.log migration/test.log
