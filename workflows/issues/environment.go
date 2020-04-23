package issues

import (
	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	strategy "github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	userObjectiveProgress "github.com/adaptiveteam/adaptive/daos/userObjectiveProgress"
)

// UserHasWriteAccessToIssues is an access policy.
// It might eventually evolve to interface
type UserHasWriteAccessToIssues = func(userID string, itype IssueType) bool

// SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated is an implementation of a query
// type SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated = func(models.TeamID) (out []strategy.CapabilityCommunity, err error)
// SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated -
// type SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated = func(models.TeamID) (out []strategy.StrategyInitiativeCommunity, err error)

type CommunityById = func(issueID string) (models.AdaptiveCommunity, error)
type PropertyName = string

// IssueDAO is an interface to read/write issue to/from database
// It should read/write all needed tables at once.
type IssueDAO interface {
	// reads all issues of the given type accessible by userID
	SelectFromIssuesWhereTypeAndUserID(userID string, issueType IssueType, completed int) ([]Issue, error)
	Read(issueType IssueType, issueID string) (Issue, error)
	Save(issue Issue) (err error)
	SetCancelled(issueID string) (err error)
	SetCompleted(issueID string) (err error)
	// SetPropertyValue updates a single field in the entity
	// SetPropertyValue( issueID string, propertyName PropertyName, value interface{}) (err error)
}

type IssueProgressDAO interface {
	// ReadAll reads at most `limit` progress elements in descending order.
	// Set limit to -1 to retrieve all the updates
	ReadAll(issueID string, limit int) ([]userObjectiveProgress.UserObjectiveProgress, error)

	Read(issueProgressID IssueProgressID) (userObjectiveProgress.UserObjectiveProgress, error)
	Save(issueProgress userObjectiveProgress.UserObjectiveProgress) (err error)
}

type AdaptiveCommunityDAO interface {
	Read(communityID community.AdaptiveCommunity) (comm models.AdaptiveCommunity, err error)
	SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated(models.TeamID) (out []strategy.CapabilityCommunity, err error)
	SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated(models.TeamID, string) (out []strategy.StrategyInitiativeCommunity, err error)
	// ReadMembers( communityID community.AdaptiveCommunity) (users []models.AdaptiveCommunityUser3, err error)
}

type CompetencyDAO interface {
	Read(id string) (adaptiveValue.AdaptiveValue, error)
	ReadAll(teamID models.TeamID) ([]adaptiveValue.AdaptiveValue, error)
}

// See strategy.StrategyObjectiveByID(teamID, each, strategyObjectivesTableName)
type StrategyObjectiveDAO interface {
	Read(id string) (models.StrategyObjective, error)
	CreateOrUpdate(so models.StrategyObjective) error
}

type StrategyCommunityDAO interface {
	Read(id string) (strategy.StrategyCommunity, error)
}

type CapabilityCommunityDAO interface {
	Read(id string) (models.CapabilityCommunity, error)
}
type StrategyInitiativeDAO interface {
	Read(id string) (models.StrategyInitiative, error)
	CreateOrUpdate(so models.StrategyInitiative) error
}
type StrategyInitiativeCommunityDAO interface {
	Read(id string) (models.StrategyInitiativeCommunity, error)
}

// SelectFromIssuesWhereTypeAndUserIDStrategyObjectives reads by the list of identifiers
func SelectFromIssuesWhereTypeAndUserIDStrategyObjectives(ids []string) func(conn DynamoDBConnection) (objectives []models.StrategyObjective, err error) {
	return func(conn DynamoDBConnection) (objectives []models.StrategyObjective, err error) {
		for _, each := range ids {
			var sos []models.StrategyObjective
			sos, err = utilsIssues.StrategyObjectiveReadOrEmpty(each)(conn)
			if err != nil {
				return
			}
			objectives = append(objectives, sos...)
		}
		return
	}
}

// Queries contains a few queries that are being used by the workflow
type Queries interface {
	// SelectFromInitiativesJoinUserCommunityWhereUserID
	// reads all initiatives that are associated with
	// the initiative communities that the user is part of.
	SelectFromInitiativesJoinUserCommunityWhereUserID(
		userID string) ([]models.StrategyInitiative, error)
	// SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID
	// returns all open objectives associated with a user.
	// If user is in strategy community, we return all objectives.
	// Else we return those objectives associated with capability communities
	// that the user is part of.
	// See strategy.UserStrategyObjectives (utils.go)
	SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID(
		userID string) ([]models.StrategyObjective, error)
	// SelectFromObjectivesWhereUserID( userID string) ([]models.StrategyObjective, error)
	/*
		func communityMembersIncludingStrategyMembers(commID string, teamID models.TeamID) []models.KvPair {
			// Strategy Community members
			strategyCommMembers := communityMembers(string(community.Strategy), teamID)
			commMembers := communityMembers(commID, teamID)
			return models.DistinctKvPairs(append(strategyCommMembers, commMembers...))
		}
	*/
	// SelectKvPairsFromCommunityUnionSelectAllFromStrategy( communityID string) (members []models.KvPair, err error)

	// SelectKvPairsFromCommunityJoinUsers loads members from community, then
	// for each member id loads UserToken and extracts display name
	/*
		func communityMembers(commID string, teamID models.TeamID) []models.KvPair {
			// Get coaching community members
			commMembers := community.CommunityMembers(communityUsersTable, commID, teamID.ToString(), communityUsersCommunityIndex)
			logger.Infof("Members in %s Community for %s platform: %s", commID, teamID, commMembers)
			var users []models.KvPair
			for _, each := range commMembers {
				// Self user checking
				ut := userTokenUnsafe(each.UserId)
				if ut.DisplayName != "" && ut.DisplayName != adaptiveBotRealName {
					users = append(users, models.KvPair{Key: ut.DisplayName, Value: each.UserId})
				}
			}
			logger.Infof("KvPairs from communities for %s community for %s platform: %s", commID, teamID, users)
			return users
		}

	*/
	SelectKvPairsFromCommunityJoinUsers(communityID community.AdaptiveCommunity) (members []models.KvPair, err error)
}
