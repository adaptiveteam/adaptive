package issues

import (
	"fmt"
	"log"

	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	objectives "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	strategy "github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/strategyCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiativeCommunity"
	"github.com/adaptiveteam/adaptive/daos/user"
	userObjectiveProgress "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	wfCommon "github.com/adaptiveteam/adaptive/workflows/common"
	"github.com/pkg/errors"
)

// This file contains implementation of interfaces from `environment.go`

// DynamoDBConnection has just what is needed for connecting to Dynamo
type DynamoDBConnection = daosCommon.DynamoDBConnection

func CreateWorkflowImpl(logger alog.AdaptiveLogger) func(conn DynamoDBConnection) workflowImpl {
	return func(conn DynamoDBConnection) workflowImpl {
		if conn.ClientID == "" {
			panic(errors.New("CreateWorkflowImpl: clientID == ''"))
		}
		impl := workflowImpl{
			WorkflowContext:    wfCommon.WorkflowContext{
				AdaptiveLogger: logger,
				DynamoDBConnection: conn,
			},
			DialogFetcherDAO: dialogFetcher.NewDAO(conn.Dynamo, dialogContentTableName(conn.ClientID)),
		}
		if impl.ClientID == "" {
			panic(errors.New("CreateWorkflowImpl 2: clientID == ''"))
		}
		return impl
	}
}

// type IssueProgressDynamoDBConnection DynamoDBConnection

// ReadAll reads at most `limit` progress elements in descending order.
// Set limit to -1 to retrieve all the updates
func IssueProgressReadAll(issueID string, limit int) func(conn DynamoDBConnection) (res []userObjectiveProgress.UserObjectiveProgress, err error) {
	return func(conn DynamoDBConnection) (res []userObjectiveProgress.UserObjectiveProgress, err error) {
		// With scan forward to true, dynamo returns list in the ascending order of the range key
		scanForward := false
		err = conn.Dynamo.QueryTableWithIndex(
			userObjectivesProgressTableName(conn.ClientID),
			awsutils.DynamoIndexExpression{
				Condition: "id = :i",
				Attributes: map[string]interface{}{
					":i": issueID,
				},
			}, map[string]string{}, scanForward, limit, &res)
		err = errors.Wrapf(err, "IssueProgressDynamoDBConnection) ReadAll(issueID=%s)", issueID)
		return
	}
}

func IssueProgressRead(issueProgressID IssueProgressID) func(conn DynamoDBConnection) (res userObjectiveProgress.UserObjectiveProgress, err error) {
	return func(conn DynamoDBConnection) (res userObjectiveProgress.UserObjectiveProgress, err error) {
		var ops []userObjectiveProgress.UserObjectiveProgress
		ops, err = userObjectiveProgress.ReadOrEmpty(issueProgressID.IssueID, issueProgressID.Date)(conn)
		if err == nil {
			if len(ops) > 0 {
				res = ops[0]
			} else {
				err = errors.New("UserObjectiveProgress " + issueProgressID.IssueID + " d: " + issueProgressID.Date + " not found")
			}
		}
		err = errors.Wrapf(err, "IssueProgressDynamoDBConnection) Read(issueProgressID=%s)", issueProgressID)
		return
	}
}

func UserObjectiveProgressSave(issueProgress userObjectiveProgress.UserObjectiveProgress) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		err = userObjectiveProgress.CreateOrUpdate(issueProgress)(conn)
		err = errors.Wrapf(err, "IssueProgressDynamoDBConnection) Read(issueProgress.ID=%s)", issueProgress.ID)
		return
	}
}

type AdaptiveCommunityDynamoDBConnection = DynamoDBConnection

func AdaptiveCommunityReadByID(communityID daosCommon.AdaptiveCommunityID) func(conn DynamoDBConnection) (comm models.AdaptiveCommunity, err error) {
	return func(conn DynamoDBConnection) (comm models.AdaptiveCommunity, err error) {
		comm = community.CommunityById(string(communityID), models.ParseTeamID(conn.PlatformID), communitiesTableName(conn.ClientID))
		// dao := adaptiveCommunity.NewDAO(conn.Dynamo, "AdaptiveCommunityDynamoDBConnection", conn.ClientID)
		// comm, err = dao.Read(communityID)
		if comm.ID != string(communityID) {
			err = fmt.Errorf("couldn't find CommunityById(communityID=%s). Instead got Id=%s", communityID, comm.ID)
		}
		return
	}
}

func SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated() func(conn AdaptiveCommunityDynamoDBConnection) (res []strategy.CapabilityCommunity, err error) {
	return func(conn AdaptiveCommunityDynamoDBConnection) (res []strategy.CapabilityCommunity, err error) {
		teamID := models.ParseTeamID(conn.PlatformID)
		defer core.RecoverToErrorVar("SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated", &err)
		capComms := strategy.AllCapabilityCommunitiesWhereChannelExists(models.ParseTeamID(conn.PlatformID))
		for _, each := range capComms {
			var comms[] strategy.StrategyCommunity
			comms, err = strategyCommunity.ReadOrEmpty(each.ID)(conn)
			if err != nil {
				err = errors.Wrapf(err, "AdaptiveCommunityDynamoDBConnection) SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(conn.PlatformID=%s)", conn.PlatformID)
				return
			}
			if len(comms) == 0 {
				log.Printf("Not found StrategyCommunityByID %s", each.ID)
			} else {
				var created bool
				created, err = strategy.IsChannelCreated(teamID, community.CapabilityPrefix, each.ID)(conn)
				if created {
					res = append(res, each)
				}
			}
		}
		return
	}
}

// func handleMenuCreateInitiative(userID, channelID string, ,
// 	message slack.InteractionCallback, deleteOriginal bool) {
// 	logger.Infof("In handleMenuCreateInitiative for user %s with platform %s", userID, conn.PlatformID)
// 	// Query all the Strategy Initiative communities
// 	var initComms []strategy.StrategyInitiativeCommunity
// 	if isMemberInCommunity(userID, community.Strategy) {
// 		initComms = strategy.AllStrategyInitiativeCommunities(conn.PlatformID, strategyInitiativeCommunitiesTable, strategyInitiativeCommunitiesPlatformIndex, strategyCommunitiesTable)
// 	} else {
// 		initComms = StrategyInitiativeCommunitiesForUserID(userID, models.TeamID(conn.PlatformID))
// 	}

// 	var adaptiveAssociatedInitComms []strategy.StrategyInitiativeCommunity
// 	// Get a list of Adaptive associated Initiative communities
// 	for _, each := range initComms {
// 		eachStrategyComms := StrategyCommunityByID(each.ID)
// 		if eachStrategyComms.ChannelCreated == 1 {
// 			adaptiveAssociatedInitComms = append(adaptiveAssociatedInitComms, each)
// 		}
// 	}
// 	logger.Infof("Adaptive associated Initiative Communities for platform %s: %s", conn.PlatformID, adaptiveAssociatedInitComms)
// 	if len(adaptiveAssociatedInitComms) > 0 {
// 		logger.Infof("Initiatives communities exist for user %s with platform %s", userID, conn.PlatformID)
// 		mc := models.MessageCallback{Module: string(community.Strategy), Source: userID, Topic: InitiativeSelectCommunityEvent,
// 			Action: string(strategy.Create)}
// 		handleMenuEvent("Select an initiative community", userID, mc,
// 			initiativeCreateSurveyInitiativeOptions(adaptiveAssociatedInitComms, conn.PlatformID))
// 		if deleteOriginal {
// 			DeleteOriginalEng(userID, channelID, message.MessageTs)
// 		}
// 	} else {
// 		handleCreateEvent(InitiativeCommunityEvent, "There are no Adaptive associated Initiative Communities. If you have already created an Initiative Community, please ask the coordinator to create a *_private_* channel, invite Adaptive and associate with the community.",
// 			userID, channelID, conn.PlatformID, message, false)
// 	}
// }
func SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated(userID string) func(conn DynamoDBConnection) (res []strategy.StrategyInitiativeCommunity, err error) {
	return func(conn DynamoDBConnection) (res []strategy.StrategyInitiativeCommunity, err error) {
		teamID := models.ParseTeamID(conn.PlatformID)
		defer core.RecoverToErrorVar("SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated", &err)
		// Query all the Strategy Initiative communities
		var initComms []strategy.StrategyInitiativeCommunity
		if isMemberInCommunity(conn, userID, community.Strategy) {
			initComms = strategy.AllStrategyInitiativeCommunitiesWhereChannelExists(models.ParseTeamID(conn.PlatformID))
		} else {
			initComms = strategy.UserStrategyInitiativeCommunities(userID,
				communityUsersTableName(conn.ClientID), communityUsersUserCommunityIndex, communityUsersUserIndex,
				strategyInitiativeCommunitiesTableName(conn.ClientID), strategyInitiativeCommunitiesPlatformIndex,
				strategyCommunitiesTableName(conn.ClientID), models.ParseTeamID(conn.PlatformID))
		}
		existingIds := map[string]struct{}{}
		// Get a list of Adaptive associated Initiative communities
		for _, each := range initComms {
			if _, ok := existingIds[each.ID]; !ok {
				existingIds[each.ID] = struct{}{}
				var comms[] strategy.StrategyCommunity
				comms, err = strategyCommunity.ReadOrEmpty(each.ID)(conn)
				if err != nil {
					err = errors.Wrapf(err, "AdaptiveCommunityDynamoDBConnection) SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated(conn.PlatformID=%s)", conn.PlatformID)
					return
				}
				if len(comms) == 0 {
					log.Printf("Not found StrategyCommunityByID %s", each.ID)
				} else {
					var created bool
					created, err = strategy.IsChannelCreated(teamID, community.InitiativePrefix, each.ID)(conn)
					if created {
						res = append(res, each)
					}
				}
			}
		}
		// 	res = removeDuplicates(res2)
		return
	}
}

func StrategyCommunityByID(id string) func(conn AdaptiveCommunityDynamoDBConnection) (comm strategyCommunity.StrategyCommunity, found bool, err error) {
	return func(conn AdaptiveCommunityDynamoDBConnection) (comm strategyCommunity.StrategyCommunity, found bool, err error) {
		var comms []strategyCommunity.StrategyCommunity
		comms, err = strategyCommunity.ReadOrEmpty(id)(conn)
		found = len(comms) > 0
		if found {
			comm = comms[0]
		}
		return
	}
}

type UserDynamoDBConnection = DynamoDBConnection

func UserRead(userID string) func(conn DynamoDBConnection) (users []models.User, err error) {
	return func(conn DynamoDBConnection) (users []models.User, err error) {
		if utilsUser.IsSpecialOrEmptyUserID(userID) {
			err = errors.Errorf("Cannot read nonexisting userID %s", userID)
		} else {
			users, err = user.ReadOrEmpty(conn.PlatformID, userID)(conn)
		}
		err = errors.Wrapf(err, "UserDynamoDBConnection) Read(userID=%s)", userID)
		return
	}
}

func mapAdaptiveCommunityUsersToUserID(users []adaptiveCommunityUser.AdaptiveCommunityUser) (userIDs []string) {
	for _, each := range users {
		userIDs = append(userIDs, each.UserID)
	}
	return
}

var requested = models.KvPair{Key: string(objectives.RequestACoachOption), Value: utilsUser.UserID_Requested}
var none = models.KvPair{Key: string(objectives.CoachNotNeededOption), Value: utilsUser.UserID_None}

// IDOCoaches returns Key-Value pairs with user id and user display name
// The set of users and the format are suitable for IDO dialog coach field.
func IDOCoaches(userID string, oldCoachIDOptional string) func(conn DynamoDBConnection) (res []models.KvPair, err error) {
	return func(conn DynamoDBConnection) (res []models.KvPair, err error) {
		defer core.RecoverToErrorVar("IDOCoaches", &err)
		var coachingMembers []adaptiveCommunityUser.AdaptiveCommunityUser
		coachingMembers, err = adaptiveCommunityUser.ReadByPlatformIDCommunityID(conn.PlatformID, string(community.Coaching))(conn)

		isOldCoachIDPresent := oldCoachIDOptional == ""
		for _, u := range coachingMembers {
			if u.UserID == oldCoachIDOptional {
				isOldCoachIDPresent = true
				break
			}
		}
		//  community.CommunityMembers(communityUsersTableName(conn.ClientID), string(community.Coaching), 
		// 	models.ParseTeamID(conn.PlatformID))
		userIDs := mapAdaptiveCommunityUsersToUserID(coachingMembers)
		if !isOldCoachIDPresent {
			userIDs = append(userIDs, oldCoachIDOptional)
		}
		res = []models.KvPair{none}
		if len(coachingMembers) > 0 { // Does this include adaptive bot name?
			res = append(res, requested)
		}
		for _, id := range userIDs {
			var users [] models.User
			users, err = user.ReadOrEmpty(conn.PlatformID, id)(conn)
			err = errors.Wrapf(err, "UserDynamoDBConnection) IDOCoaches(userID=%s)", userID)
			if err == nil {
				for _, user := range users {
					if (user.ID != userID || userID == oldCoachIDOptional) &&
					!user.IsAdaptiveBot &&
					user.DisplayName != "" {
						res = append(res, models.KvPair{Key: user.DisplayName, Value: id})
					}
				}
			} else {
				return
			}
		}

		return
	}
}

var CompetencyRead = adaptiveValue.Read

func CompetencyReadAll() func(conn DynamoDBConnection) (res []adaptiveValue.AdaptiveValue, err error) {
	return func(conn DynamoDBConnection) (res []adaptiveValue.AdaptiveValue, err error) {
		res, err = adaptiveValue.ReadByPlatformID(conn.PlatformID)(conn)
		err = errors.Wrapf(err, "CompetencyDynamoDBConnection) ReadAll(conn.PlatformID=%s)", conn.PlatformID)
		return
	}
}

func StrategyObjectiveCreateOrUpdate(so models.StrategyObjective) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		if so.ID == "" {
			err = errors.New("ID is empty")
		} else if so.PlatformID == "" {
			err = fmt.Errorf("PlatformID is empty for ID=%s", so.ID)
		} else if so.CapabilityCommunityIDs == nil {
			err = fmt.Errorf("CapabilityCommunityIDs is empty for ID=%s", so.ID)
		}
		if err == nil {
			err = conn.Dynamo.PutTableEntry(so, strategyObjectivesTableName(conn.ClientID))
		}
		err = errors.Wrapf(err, "StrategyObjectiveDynamoDBConnection) CreateOrUpdate(so.ID=%s)", so.ID)

		return
	}
}

// type StrategyCommunityDynamoDBConnection DynamoDBConnection

func StrategyCommunityRead(id string) func(conn DynamoDBConnection) (res strategy.StrategyCommunity, err error) {
	return func(conn DynamoDBConnection) (res strategy.StrategyCommunity, err error) {
		defer core.RecoverToErrorVar("StrategyCommunityDynamoDBConnection.Read", &err)
		res = strategy.StrategyCommunityByID(id, strategyCommunityTableName(conn.ClientID))
		if res.ID != id {
			err = fmt.Errorf("couldn't find StrategyCommunityByID(id=%s). Instead got ID=%s", id, res.ID)
		}
		return
	}
}

// type StrategyInitiativeDynamoDBConnection DynamoDBConnection

func StrategyInitiativeRead(id string) func(conn DynamoDBConnection) (res models.StrategyInitiative, err error) {
	return func(conn DynamoDBConnection) (res models.StrategyInitiative, err error) {
		defer core.RecoverToErrorVar("StrategyInitiativeDynamoDBConnection.Read", &err)
		res = strategy.StrategyInitiativeByID(
			models.ParseTeamID(conn.PlatformID), id, 
			strategyInitiativesTableName(conn.ClientID))
		if res.ID != id {
			err = fmt.Errorf("couldn't find StrategyInitiativeByID(id=%s). Instead got ID=%s", id, res.ID)
		}
		return
	}
}
func StrategyInitiativeCreateOrUpdate(si models.StrategyInitiative) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		err = conn.Dynamo.PutTableEntry(si, strategyInitiativesTableName(conn.ClientID))
		err = errors.Wrapf(err, "StrategyObjectiveDynamoDBConnection) CreateOrUpdate(si.ID=%s)", si.ID)
		return
	}
}

// StrategyInitiativeCommunityRead -
func StrategyInitiativeCommunityRead(id string) func(conn DynamoDBConnection) (res models.StrategyInitiativeCommunity, err error) {
	return func(conn DynamoDBConnection) (res models.StrategyInitiativeCommunity, err error) {
		res, err = strategyInitiativeCommunity.Read(id, conn.PlatformID)(conn)
		err = errors.Wrapf(err, "StrategyInitiativeCommunityRead(id=%s)", id)
		return
	}
}

// type QueriesDynamoDBConnection DynamoDBConnection

// Queries contains a few queries that are being used by the workflow
// SelectFromInitiativesJoinUserCommunityWhereUserID
// reads all initiatives that are associated with
// the initiative communities that the user is part of.
func SelectFromInitiativesJoinUserCommunityWhereUserID(userID string) func(conn DynamoDBConnection) (res []models.StrategyInitiative, err error) {
	return func(conn DynamoDBConnection) (res []models.StrategyInitiative, err error) {
		defer core.RecoverToErrorVar("SelectFromInitiativesJoinUserCommunityWhereUserID", &err)
		res = strategy.UserInitiativeCommunityInitiatives(userID, strategyInitiativesTableName(conn.ClientID),
			strategyInitiativesInitiativeCommunityIndex, communityUsersTableName(conn.ClientID), communityUsersUserIndex)

		return
	}
}

// SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID
// returns all open objectives associated with a user.
// If user is in strategy community, we return all objectives.
// Else we return those objectives associated with capability communities
// that the user is part of.
// See strategy.UserStrategyObjectives (utils.go)
func SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID(userID string) func(conn DynamoDBConnection) (res []models.StrategyObjective, err error) {
	return func(conn DynamoDBConnection) (res []models.StrategyObjective, err error) {
		defer core.RecoverToErrorVar("SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID", &err)
		objs := strategy.UserStrategyObjectives(userID, strategyObjectivesTableName(conn.ClientID),
			strategyObjectivesPlatformIndex, userObjectivesTableName(conn.ClientID),
			communityUsersTableName(conn.ClientID), communityUsersUserCommunityIndex,
			conn,
		)
		uniqueIDs := make(map[string]struct{})
		for _, obj := range objs {
			if _, ok := uniqueIDs[obj.ID]; !ok {
				uniqueIDs[obj.ID] = struct{}{}
				res = append(res, obj)
			}
		}
		return
	}
}

func SelectKvPairsFromCommunityJoinUsers(communityID community.AdaptiveCommunity) func(DynamoDBConnection) ([]models.KvPair, error) {
	return func(conn DynamoDBConnection) (members []models.KvPair, err error) {
		defer core.RecoverToErrorVar("SelectKvPairsFromCommunityJoinUsers", &err)
		log.Printf("Before calling CommunityMembers")
		commMembers := community.CommunityMembers(communityUsersTableName(conn.ClientID), string(communityID), 
			models.ParseTeamID(conn.PlatformID) )//, string(adaptiveCommunityUser.PlatformIDCommunityIDIndex)) // communityUsersCommunityIndex)
		log.Printf("After calling CommunityMembers")
		log.Printf("Found %d in community %s", len(commMembers), communityID)
		for _, each := range commMembers {
			// Self user checking
			us := user.ReadOrEmptyUnsafe(conn.PlatformID, each.UserId)(conn)
			if len(us) == 0 {
				log.Printf("Not found user %s", each.UserId)
			}
			for _, u := range us {
				if u.DisplayName != "" && !u.IsAdaptiveBot {
					members = append(members, models.KvPair{Key: u.DisplayName, Value: each.UserId})
				} else {
					log.Printf("Ignoring user %s with display name %s", each.UserId, u.DisplayName)
				}
			}
		}
		return
	}
}

func closeoutLabel(userObjID string) ui.PlainText {
	return ui.PlainText("Responsibility Closeout")
}

func isMemberInCommunity(conn DynamoDBConnection, userID string, comm community.AdaptiveCommunity) bool {
	defer RecoverToLog("DynamoDBConnection) isMemberInCommunity")
	return community.IsUserInCommunity(userID, communityUsersTableName(conn.ClientID), communityUsersUserCommunityIndex, comm)
}

func UserHasWriteAccessToIssuesImpl(conn DynamoDBConnection) func(userID string, itype IssueType) bool {
	return func(userID string, itype IssueType) (allow bool) {
		switch itype {
		case IDO:
			allow = true
		case SObjective, Initiative:
			allow = isMemberInCommunity(conn, userID, community.Strategy)
		}
		return
	}
}

func LoadObjectives(userID string) func(conn DynamoDBConnection) (objKVs []models.KvPair) {
	return func(conn DynamoDBConnection) (objKVs []models.KvPair) {
		isStrategyUser := isMemberInCommunity(conn, userID, community.Strategy)
		var objs []models.StrategyObjective
		if isStrategyUser {
			objs = strategy.UserStrategyObjectives(userID,
				strategyObjectivesTableName(conn.ClientID), strategyObjectivesPlatformIndex,
				userObjectivesTableName(conn.ClientID),
				communityUsersTableName(conn.ClientID), communityUsersUserIndex,
				conn,
			)
		} else {
			objs = strategy.UserCommunityObjectives(userID,
				strategyObjectivesTableName(conn.ClientID), strategyObjectivesPlatformIndex,
				userObjectivesTableName(conn.ClientID),
				communityUsersTableName(conn.ClientID), communityUsersUserIndex,
				conn,
			)
		}
		for _, eachObj := range objs {
			objKVs = append(objKVs, models.KvPair{Key: eachObj.Name, Value: eachObj.ID})
		}
		return
	}
}

func removeDuplicates(kvPairs []models.KvPair) (res []models.KvPair) {
	values := make(map[string]struct{})
	for _, p := range kvPairs {
		if _, ok := values[p.Value]; !ok {
			res = append(res, p)
			values[p.Value] = struct{}{}
		}
	}
	return
}
