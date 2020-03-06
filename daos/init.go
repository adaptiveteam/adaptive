package daos

import (
	"log"

	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunity"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/capabilityCommunity"
	"github.com/adaptiveteam/adaptive/daos/clientPlatformToken"
	"github.com/adaptiveteam/adaptive/daos/coachingRelationship"
	"github.com/adaptiveteam/adaptive/daos/contextAliasEntry"
	"github.com/adaptiveteam/adaptive/daos/dialogEntry"
	"github.com/adaptiveteam/adaptive/daos/strategyCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiative"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiativeCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/user"
	"github.com/adaptiveteam/adaptive/daos/userFeedback"
	"github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	"github.com/adaptiveteam/adaptive/daos/visionMission"
)

func init() {
	log.Printf("Initializing table names")

	adHocHoliday.TableNameSuffixVar = "_ad_hoc_holidays"
	clientPlatformToken.TableNameSuffixVar = "_adaptive_client_config"
	// dialogEntry.TableNameSuffixVar = "_adaptive_dialog"
	userFeedback.TableNameSuffixVar = "_adaptive_user_feedback"
	user.TableNameSuffixVar = "_adaptive_users"
	// .TableNameSuffixVar = "_adaptive_value"
	capabilityCommunity.TableNameSuffixVar = "_capability_communities"
	// .TableNameSuffixVar = "_coaching_rejections"
	coachingRelationship.TableNameSuffixVar = "_coaching_relationships"
	adaptiveCommunity.TableNameSuffixVar = "_communities"
	adaptiveCommunityUser.TableNameSuffixVar = "_community_users"
	dialogEntry.TableNameSuffixVar = "_dialog_content"
	contextAliasEntry.TableNameSuffixVar = "_dialog_content_alias"
	strategyInitiativeCommunity.TableNameSuffixVar = "_initiative_communities"
	// .TableNameSuffixVar = "_objective_type_dictionary"
	// .TableNameSuffixVar = "_partnership_rejections"
	// .TableNameSuffixVar = "_postponed_event"
	// .TableNameSuffixVar = "_slack_team"
	strategyCommunity.TableNameSuffixVar = "_strategy_communities"
	strategyInitiative.TableNameSuffixVar = "_strategy_initiatives"
	strategyObjective.TableNameSuffixVar = "_strategy_objectives"
	// .TableNameSuffixVar = "_user_engagement"
	// .TableNameSuffixVar = "_user_objective"
	userObjectiveProgress.TableNameSuffixVar = "_user_objectives_progress"
	visionMission.TableNameSuffixVar = "_vision"
	log.Println("New table names:")
	log.Println(adHocHoliday.TableNameSuffixVar)
	log.Println(clientPlatformToken.TableNameSuffixVar)
	log.Println(userFeedback.TableNameSuffixVar)
	log.Println(user.TableNameSuffixVar)
	log.Println(capabilityCommunity.TableNameSuffixVar)
	log.Println(coachingRelationship.TableNameSuffixVar)
	log.Println(adaptiveCommunity.TableNameSuffixVar)
	log.Println(adaptiveCommunityUser.TableNameSuffixVar)
	log.Println(dialogEntry.TableNameSuffixVar)
	log.Println(contextAliasEntry.TableNameSuffixVar)
	log.Println(strategyInitiativeCommunity.TableNameSuffixVar)
	log.Println(strategyCommunity.TableNameSuffixVar)
	log.Println(strategyInitiative.TableNameSuffixVar)
	log.Println(strategyObjective.TableNameSuffixVar)
	log.Println(userObjectiveProgress.TableNameSuffixVar)
	log.Println(visionMission.TableNameSuffixVar)
}
