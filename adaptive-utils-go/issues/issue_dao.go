package issues

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	strategy "github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsUser "github.com/adaptiveteam/adaptive/adaptive-utils-go/user"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	"github.com/adaptiveteam/adaptive/daos/capabilityCommunity"
	"github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiative"
	"github.com/adaptiveteam/adaptive/daos/strategyInitiativeCommunity"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/pkg/errors"
)

type DynamoDBConnection = common.DynamoDBConnection

var (
	dialogContentTableName                     = func(clientID string) string { return clientID + "_dialog_content" }
	strategyObjectiveTableName                 = func(clientID string) string { return clientID + "_strategy_objectives" }
	strategyInitiativeTableName                = func(clientID string) string { return clientID + "_strategy_initiatives" }
	strategyInitiativeInitiativeCommunityIndex = "InitiativeCommunityIDIndex"
	userObjectiveTableName                     = func(clientID string) string { return clientID + "_user_objective" }
	userObjectiveIDIndex                       = "IDIndex"
	userObjectiveUserIDIndex                   = "UserIDCompletedIndex"
	userObjectiveTypeIndex                     = "UserIDTypeIndex"
	userObjectiveProgressTableName             = func(clientID string) string { return clientID + "_user_objectives_progress" }
	adaptiveCommunityUserTableName             = func(clientID string) string { return clientID + "_community_users" }
	communityTableName                         = func(clientID string) string { return clientID + "_communities" }
	competencyTableName                        = func(clientID string) string { return clientID + "_adaptive_value" }
	strategyInitiativeCommunityTableName       = func(clientID string) string { return clientID + "_initiative_communities" }
	strategyInitiativeCommunityPlatformIndex   = "PlatformIDIndex"
	strategyCommunityTableName                 = func(clientID string) string { return clientID + "_strategy_communities" }
	visionMissionTableName                     = func(clientID string) string { return clientID + "_vision" }
	capabilityCommunityTableName               = func(clientID string) string { return clientID + "_capability_communities" }
	capabilityCommunityPlatformIndex           = "PlatformIDIndex"
	adaptiveUserTableName                      = func(clientID string) string { return clientID + "_adaptive_users" }
)

const communityUsersUserCommunityIndex = "UserIDCommunityIDIndex"
const strategyObjectivesPlatformIndex = "PlatformIDIndex"
const strategyInitiativesPlatformIndex = "PlatformIDIndex"
const communityUsersUserIndex = "UserIDIndex"

// SelectFromIssuesWhereTypeAndUserID reads all issues of the given type accessible by userID
func SelectFromIssuesWhereTypeAndUserID(userID string, issueType IssueType, completed int) func(conn common.DynamoDBConnection) (res []Issue, err error) {
	switch issueType {
	case IDO:
		return selectFromIssuesWhereTypeAndUserIDIDO(userID, completed)
	case SObjective:
		return selectFromIssuesWhereTypeAndUserIDSObjective(userID, completed)
	case Initiative:
		return selectFromIssuesWhereTypeAndUserIDInitiative(userID, completed)
	}
	// should never happen
	return selectFromIssuesWhereTypeAndUserIDIDO(userID, completed)
}

func selectFromIssuesWhereTypeAndUserIDIDO(userID string, completed int) func(conn DynamoDBConnection) (res []Issue, err error) {
	return func(conn DynamoDBConnection) (res []Issue, err error) {
		var objs []userObjective.UserObjective
		objs, err = userObjective.ReadByUserIDCompleted(userID, completed)(conn)
		if err == nil {
			for _, o := range objs {
				if o.ObjectiveType == userObjective.IndividualDevelopmentObjective { // o.Completed == completed { // this should be automatic
					res = append(res, Issue{UserObjective: o})
				}
			}
		}
		err = errors.Wrapf(err, "selectFromIssuesWhereTypeAndUserIDIDO(userID=%s)", userID)
		return
	}
}

func platformIndexExpr(index string, teamID models.TeamID) awsutils.DynamoIndexExpression {
	return awsutils.DynamoIndexExpression{
		IndexName: index,
		Condition: "platform_id = :p",
		Attributes: map[string]interface{}{
			":p": teamID.ToString(),
		},
	}
}

func selectFromIssuesWhereTypeAndUserIDSObjective(userID string, completed int) func(conn common.DynamoDBConnection) (res []Issue, err error) {
	return func(conn common.DynamoDBConnection) (res []Issue, err error) {
		defer core.RecoverToErrorVar("selectFromIssuesWhereTypeAndUserIDSObjective", &err)

		var allObjs []models.StrategyObjective
		err = conn.Dynamo.QueryTableWithIndex(
			strategyObjectiveTableName(conn.ClientID),
			platformIndexExpr(strategyObjectivesPlatformIndex, models.ParseTeamID(conn.PlatformID)),
			map[string]string{}, true, -1, &allObjs)
		log.Printf("AllStrategyObjectives: len(allObjs)=%d\n", len(allObjs))

		for _, each := range allObjs {
			// there has to be at least one objective community id
			// TODO: This presents a tricky scenario when original objective community is updated. Think about this.
			// Customer and financial objectives have no capability communities associated with them. For them,we only use the ID
			id := each.ID
			var objs []userObjective.UserObjective
			objs, err = userObjective.ReadOrEmpty(id)(conn)
			if err != nil {
				err = errors.Wrapf(err, "DynamoDBConnection) selectFromIssuesWhereTypeAndUserIDSObjective/userObjectiveDao.ReadOrEmpty")
				return
			}
			if len(objs) == 0 && len(each.CapabilityCommunityIDs) > 0 {
				id = fmt.Sprintf("%s_%s", each.ID, each.CapabilityCommunityIDs[0])
				objs, err = userObjective.ReadOrEmpty(id)(conn)
				if err != nil {
					err = errors.Wrapf(err, "DynamoDBConnection) selectFromIssuesWhereTypeAndUserIDSObjective/userObjectiveDao.ReadOrEmpty 1_2")
					return
				}
			}
			var issue Issue
			if len(objs) > 0 {
				issue = Issue{
					StrategyObjective: each,
					UserObjective:     objs[0],
				}
			} else {
				// err = errors.New("UserObjective " + each.ID + " or " + id + " not found")
				log.Printf("selectFromIssuesWhereTypeAndUserIDSObjective: Not found user objective for %s or %s\n", each.ID, id) // err)
				// err = nil
				var uo userObjective.UserObjective
				uo, err = UserObjectiveFromStrategyObjective(each)(conn)
				if err != nil {
					err = errors.Wrapf(err, "DynamoDBConnection) selectFromIssuesWhereTypeAndUserIDSObjective 3: Couldn't convert strategy objective %s to user objective", each.ID)
					return
				}
				issue = Issue{
					StrategyObjective: each,
					UserObjective:     uo,
				}
			}
			if issue.UserObjective.Completed == 0 {
				res = append(res, issue)
			}
		}

		// itemsl := strategy.AllOpenStrategyObjectives(
		// 	conn.PlatformID,
		// 	strategyObjectivesTableName(conn.ClientID), strategyObjectivesPlatformIndex,
		// 	userObjectivesTableName(conn.ClientID), userObjectivesIDIndex)
		// for _, so := range items {
		// 	log.Printf("found strategy objective: so.ID=%s\n",so.ID)
		// 	issue := Issue{StrategyObjective: so}
		// 	issue.UserObjective, err = common.DynamoDBConnection(conn).UserObjectiveFromStrategyObjective(so, conn.PlatformID)
		// 	if err != nil {
		// 		err = errors.Wrapf(err, "selectFromIssuesWhereTypeAndUserIDSObjective(userID=%s)", userID)
		// 		return
		// 	}
		// 	res = append(res, issue)
		// }
		return
	}
}

func UserObjectiveFromStrategyObjective(so models.StrategyObjective) func(conn common.DynamoDBConnection) (uObj userObjective.UserObjective, err error) {
	return func(conn common.DynamoDBConnection) (uObj userObjective.UserObjective, err error) {
		defer core.RecoverToErrorVar("UserObjectiveFromStrategyObjective", &err)
		var commID string
		if len(so.CapabilityCommunityIDs) > 0 {
			commID = so.CapabilityCommunityIDs[0]
		} else {
			commID = ""
			log.Printf("UserObjectiveFromStrategyObjective: CapabilityCommunityIDs is empty")
		}
		vision := strategy.StrategyVision(models.ParseTeamID(conn.PlatformID), visionMissionTableName(conn.ClientID))
		// We are using _ here because `:` will create issues with callback
		id := so.ID // core.IfThenElse(commID != core.EmptyString, fmt.Sprintf("%s_%s", so.ID, commID), so.ID).(string)
		// log.Printf("UserObjectiveFromStrategyObjective: id=%s; so.ID=%s, commID=%s, so=%v\n",  id, so.ID, commID, so)

		if id == "" {
			log.Printf("UserObjectiveFromStrategyObjective: id is empty; so.ID=%s, commID=%s, so=%v\n", so.ID, commID, so)
		}
		createdDate := core.NormalizeDate(so.CreatedAt)
		uObj = userObjective.UserObjective{
			ID:              id,
			Name:            so.Name,
			Description:     so.Description,
			UserID:          vision.Advocate, // TODO: why?
			Accepted:        1,
			ObjectiveType:   userObjective.StrategyDevelopmentObjective,
			PlatformID:      conn.PlatformID,
			CreatedDate:     createdDate,
			ExpectedEndDate: so.ExpectedEndDate,
			CreatedAt:       so.CreatedAt,
			// ModifiedAt:      so.ModifiedAt,

			AccountabilityPartner:       so.Advocate,
			StrategyAlignmentEntityID:   commID,
			StrategyAlignmentEntityType: userObjective.ObjectiveStrategyObjectiveAlignment,
		}
		err = errors.Wrapf(err, "UserObjectiveFromStrategyObjective(so.ID=%s)", so.ID)
		return
	}
}

func IsMemberInCommunity(userID string, comm community.AdaptiveCommunity) func(conn DynamoDBConnection) bool {
	return func(conn DynamoDBConnection) bool {
		defer core.RecoverAsLogError("issues_dao.go: IsMemberInCommunity")
		return community.IsUserInCommunity(userID, adaptiveCommunityUserTableName(conn.ClientID), communityUsersUserCommunityIndex, comm)
	}
}
func selectFromIssuesWhereTypeAndUserIDInitiative(userID string, completed int) func(conn DynamoDBConnection) (res []Issue, err error) {
	return func(conn DynamoDBConnection) (res []Issue, err error) {
		inStrategyCommunity := IsMemberInCommunity(userID, community.Strategy)(conn)
		var res2 []Issue
		if inStrategyCommunity {
			log.Printf("selectFromIssuesWhereTypeAndUserIDInitiative, inStrategyCommunity=true")
			// User is in Strategy community, show all Initiatives
			res2, err = IssuesFromAllStrategyInitiatives()(conn)
			log.Printf("selectFromIssuesWhereTypeAndUserIDInitiative, AllStrategyInitiatives.count=%d", len(res))
			if err != nil {
				return
			}
		}
		var res1 []Issue
		res1, err = IssuesFromCapabilityCommunityInitiatives(userID)(conn)
		if err == nil {
			log.Printf("selectFromIssuesWhereTypeAndUserIDInitiative, CapabilityCommunityInitiatives.count=%d", len(res1))

			res2 = append(res2, res1...)
			res1, err = IssuesFromInitiativeCommunityInitiatives(userID)(conn)
			log.Printf("selectFromIssuesWhereTypeAndUserIDInitiative, IssuesFromInitiativeCommunityInitiatives.count=%d", len(res1))
			if err == nil {
				res2 = append(res2, res1...)
			}
		}
		res = removeDuplicates(res2)
		err = errors.Wrapf(err, "selectFromIssuesWhereTypeAndUserIDInitiative(userID=%s)", userID)
		return
	}
}

func removeDuplicates(issues []Issue) (res []Issue) {
	existingIDs := map[string]struct{}{}
	for _, each := range issues {
		if _, ok := existingIDs[each.UserObjective.ID]; !ok {
			existingIDs[each.UserObjective.ID] = struct{}{}
			res = append(res, each)
		} else {
			log.Printf("Found duplicate %s", each.UserObjective.ID)
		}
	}
	return
}

func IssuesFromAllStrategyInitiatives() func(conn DynamoDBConnection) (res []Issue, err error) {
	return func(conn DynamoDBConnection) (res []Issue, err error) {
		defer core.RecoverToErrorVar("AllStrategyInitiatives", &err)
		inits := strategy.AllOpenStrategyInitiatives(models.ParseTeamID(conn.PlatformID),
			strategyInitiativeTableName(conn.ClientID),
			strategyInitiativesPlatformIndex,
			userObjectiveTableName(conn.ClientID))
		res, err = IssuesFromGivenStrategyInitiatives(inits)(conn)
		err = errors.Wrapf(err, "AllStrategyInitiatives(conn.PlatformID=%s)", conn.PlatformID)
		return
	}
}

func IssuesFromGivenStrategyInitiatives(inits []models.StrategyInitiative) func(conn DynamoDBConnection) (issues []Issue, err error) {
	return func(conn DynamoDBConnection) (issues []Issue, err error) {
		for _, si := range inits {
			var issue Issue
			issue, err = IssueFromStrategyInitiative(si)(conn)
			if err != nil {
				err = errors.Wrapf(err, "IssuesFromStrategyInitiatives(si.ID=%s)", si.ID)
				return
			}
			issues = append(issues, issue)
		}
		return
	}
}
func IssuesFromCapabilityCommunityInitiatives(userID string) func(conn DynamoDBConnection) (res []Issue, err error) {
	return func(conn DynamoDBConnection) (res []Issue, err error) {
		defer core.RecoverToErrorVar("CapabilityCommunityInitiatives", &err)
		strategyInitiativesInitiativeCommunityIndex := "InitiativeCommunityIDIndex"
		inits := strategy.UserCapabilityCommunityInitiatives(userID,
			strategyObjectiveTableName(conn.ClientID), strategyObjectivesPlatformIndex,
			strategyInitiativeTableName(conn.ClientID), strategyInitiativesInitiativeCommunityIndex,
			userObjectiveTableName(conn.ClientID),
			adaptiveCommunityUserTableName(conn.ClientID), communityUsersUserCommunityIndex,
			communityUsersUserIndex)
		res, err = IssuesFromGivenStrategyInitiatives(inits)(conn)
		err = errors.Wrapf(err, "CapabilityCommunityInitiatives(userID=%s)", userID)
		return
	}
}

func IssuesFromInitiativeCommunityInitiatives(userID string) func(conn DynamoDBConnection) (res []Issue, err error) {
	return func(conn DynamoDBConnection) (res []Issue, err error) {
		defer core.RecoverToErrorVar("IssuesFromInitiativeCommunityInitiatives", &err)
		var inits []models.StrategyInitiative
		inits = strategy.AllOpenStrategyInitiatives(models.ParseTeamID(conn.PlatformID),
			strategyInitiativeTableName(conn.ClientID), strategyInitiativesPlatformIndex,
			userObjectiveTableName(conn.ClientID))
		res, err = IssuesFromGivenStrategyInitiatives(inits)(conn)
		err = errors.Wrapf(err, "IssuesFromInitiativeCommunityInitiatives(userID=%s)", userID)
		return
	}
}

func IssueFromStrategyInitiative(si models.StrategyInitiative) func(conn common.DynamoDBConnection) (issue Issue, err error) {
	return func(conn common.DynamoDBConnection) (issue Issue, err error) {
		advocate := ""
		if si.CapabilityObjective != "" {
			capObj := strategy.StrategyObjectiveByID(models.ParseTeamID(conn.PlatformID), si.CapabilityObjective,
				strategyObjectiveTableName(conn.ClientID))
			advocate = capObj.Advocate
		}
		createdDate := core.NormalizeDate(si.CreatedAt)
		issue.StrategyInitiative = si
		issue.UserObjective = userObjective.UserObjective{
			UserID:                      si.Advocate,
			Name:                        si.Name,
			ID:                          si.ID,
			Description:                 si.Description,
			AccountabilityPartner:       advocate,
			Accepted:                    1,
			ObjectiveType:               userObjective.StrategyDevelopmentObjective,
			StrategyAlignmentEntityID:   "", //si.InitiativeCommunityID,
			StrategyAlignmentEntityType: userObjective.ObjectiveStrategyInitiativeAlignment,
			PlatformID:                  conn.PlatformID,
			CreatedDate:                 createdDate,
			ExpectedEndDate:             si.ExpectedEndDate,
		}
		err = errors.Wrapf(err, "IssueFromStrategyInitiative(si.ID=%s)", si.ID)
		return
	}
}

var StrategyInitiativeReadOrEmpty = strategyInitiative.ReadOrEmpty

func StrategyInitiativeCreateOrUpdate(si models.StrategyInitiative) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		err = conn.Dynamo.PutTableEntry(si, strategyInitiativeTableName(conn.ClientID))
		err = errors.Wrapf(err, "StrategyObjectiveDynamoDBConnection) CreateOrUpdate(si.ID=%s)", si.ID)
		return
	}
}

// StrategyObjectiveReadOrEmpty attempts to find the strategy objective.
// Since we have corrupted DB structure, we not always find it. Hence we have to
// skip not found objectives.
func StrategyObjectiveReadOrEmpty(id string) func(conn DynamoDBConnection) (res []models.StrategyObjective, err error) {
	return func(conn DynamoDBConnection) (res []models.StrategyObjective, err error) {
		defer core.RecoverToErrorVar("StrategyObjectiveDynamoDBConnection.Read", &err)
		id2 := id
		i := strings.Index(id2, "_")
		if i >= 0 {
			log.Printf("WARN: StrategyObjectiveDynamoDBConnection) Read: ID has '_': %s\n", id)
			id2 = id[0:i]
		}
		// log.Printf("StrategyObjectiveDynamoDBConnection) Read: reading id2=%s\n", id2)
		res, err = strategyObjective.ReadOrEmpty(conn.PlatformID, id2)(conn)
		if len(res) > 0 && res[0].ID != id2 {
			err = fmt.Errorf("couldn't find StrategyObjectiveByID(id2=%s, id=%s). Instead got ID=%s", id2, id, res[0].ID)
		}
		return
	}
}

// CapabilityCommunityReadOrEmpty -
func CapabilityCommunityReadOrEmpty(id string) func(conn DynamoDBConnection) (res []models.CapabilityCommunity, err error) {
	return func(conn DynamoDBConnection) (res []models.CapabilityCommunity, err error) {
		defer core.RecoverToErrorVar("CapabilityCommunityReadOrEmpty", &err)
		res, err = capabilityCommunity.ReadOrEmpty(conn.PlatformID, id)(conn)
		if len(res) > 0 && res[0].ID != id {
			err = fmt.Errorf("couldn't find CapabilityCommunityByID(id=%s). Instead got ID=%s", id, res[0].ID)
		}
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
			err = conn.Dynamo.PutTableEntry(so, strategyObjectiveTableName(conn.ClientID))
		}
		err = errors.Wrapf(err, "StrategyObjectiveDynamoDBConnection) CreateOrUpdate(so.ID=%s)", so.ID)

		return
	}
}

func Read(issueType IssueType, issueID string) func(conn DynamoDBConnection) (issue Issue, err error) {
	return func(conn DynamoDBConnection) (issue Issue, err error) {
		defer core.RecoverToErrorVar("DynamoDBConnection) Read", &err)
		if issueID == "" {
			err = fmt.Errorf("%s issue id is empty", issueType)
			return
		}
		var objs []userObjective.UserObjective
		objs, err = userObjective.ReadOrEmpty(issueID)(conn)
		if len(objs) > 0 {
			issue.UserObjective = objs[0]
		}
		switch issueType {
		case IDO:
			// We don't need to read anything additional. It's only UserObjective
			// Though it's required that we have found one objective
			if len(objs) == 0 {
				err = errors.New("UserObjective " + issueID + " not found")
			}
		case SObjective:
			// dao := strategyObjective.NewDAO(conn.Dynamo, "issues_dao", conn.ClientID)
			var sos []strategyObjective.StrategyObjective
			sos, err = StrategyObjectiveReadOrEmpty(issueID)(conn)
			if len(sos) > 0 {
				issue.StrategyObjective = sos[0]
				if err == nil && len(objs) == 0 {
					issue.UserObjective, err = UserObjectiveFromStrategyObjective(issue.StrategyObjective)(conn)
				}
			}
		case Initiative:
			var inis []strategyInitiative.StrategyInitiative
			inis, err = StrategyInitiativeReadOrEmpty(conn.PlatformID, issueID)(conn)
			if len(inis) > 0 {
				issue.StrategyInitiative = inis[0]
				if err == nil && len(objs) == 0 {
					issue, err = IssueFromStrategyInitiative(issue.StrategyInitiative)(conn)
				}	
			}			
		}
		issue.NormalizeIssueDateTimes()
		err = errors.Wrapf(err, "issue_dao.Read(issueType=%s, ID=%s)", issueType, issueID)
		return
	}
}

func Save(issue Issue) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		log.Printf("DynamoDBConnection) Save(uo.ID=%s, so.ID=%s, si.ID=%s, issue=%v)\n",
			issue.UserObjective.ID, issue.StrategyObjective.ID, issue.StrategyInitiative.ID, issue)
		err = userObjective.CreateOrUpdate(issue.UserObjective)(conn)
		if err == nil {
			switch issue.GetIssueType() {
			case IDO:
				// already saved above
			case SObjective:
				log.Printf("DynamoDBConnection) Save SObjective(so.ID=%s)\n", issue.StrategyObjective.ID)
				// sdao := strategyObjective.NewDAO(conn.Dynamo, "issues_dao", conn.ClientID)

				err = StrategyObjectiveCreateOrUpdate(issue.StrategyObjective)(conn)
				if err == nil {
					var sos []models.StrategyObjective
					sos, err = StrategyObjectiveReadOrEmpty(issue.StrategyObjective.ID)(conn)
					if err == nil && len(sos) > 0 {
						so := sos[0]
						if so.ID == issue.StrategyObjective.ID {
							log.Printf("DynamoDBConnection) Saved successfully SObjective(so.ID=%s)%+v\n", issue.StrategyObjective.ID, err)
						} else {
							var bytes []byte
							bytes, err = json.Marshal(issue.StrategyObjective)
							log.Printf("DynamoDBConnection) NOT Saved SObjective(so.ID=%s) without any error. Table name: '%s'. Value:\n%v\n", issue.StrategyObjective.ID, strategyObjectiveTableName(conn.ClientID),
								string(bytes))
						}
					}
				}

			case Initiative:
				// idao := strategyInitiative.NewDAO(conn.Dynamo, "issues_dao", conn.ClientID)
				err = StrategyInitiativeCreateOrUpdate(issue.StrategyInitiative)(conn)
			}
		}
		err = errors.Wrapf(err, "DynamoDBConnection) Save(issue.UserObjective.ID=%s)", issue.UserObjective.ID)
		return
	}
}

// SetCancelled updates a single field in the entity - Cancelled - to true
func SetCancelled(issueID string) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		var objs []userObjective.UserObjective
		objs, err = userObjective.ReadOrEmpty(issueID)(conn)
		if err == nil {
			if len(objs) > 0 {
				objs[0].Cancelled = 1
				objs[0].Completed = 1
				objs[0].CompletedDate = core.ISODateLayout.Format(time.Now())
				err = userObjective.CreateOrUpdate(objs[0])(conn)
			} else {
				err = errors.New("UserObjective " + issueID + " not found (SetCancelled)")
			}
		}
		err = errors.Wrapf(err, "DynamoDBConnection) SetCancelled(issueID=%s)", issueID)
		return
	}
}

// SetCompleted updates a single field in the entity - Completed - to true
func SetCompleted(issueID string) func(conn DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		var objs []userObjective.UserObjective
		objs, err = userObjective.ReadOrEmpty(issueID)(conn)
		if len(objs) > 0 {
			objs[0].Completed = 1
			objs[0].CompletedDate = core.ISODateLayout.Format(time.Now())
			err = userObjective.CreateOrUpdate(objs[0])(conn)
		} else {
			err = errors.New("UserObjective " + issueID + " not found (SetCompleted)")
		}
		err = errors.Wrapf(err, "DynamoDBConnection) SetCompleted(issueID=%s)", issueID)
		return
	}
}

// ReadNewAndOldIssuesAndPrefetch loads issue and prefetches all dictionaries
// NB! Only the new issue is loaded and prefetched!
func ReadNewAndOldIssuesAndPrefetch(issueType IssueType, issueID string, isShowingProgress bool) func(DynamoDBConnection) (newAndOldIssues NewAndOldIssues, err error) {
	return func(DynamoDBConnection DynamoDBConnection) (newAndOldIssues NewAndOldIssues, err error) {
		newAndOldIssues.NewIssue, err = Read(issueType, issueID)(DynamoDBConnection)
		if err != nil {
			err = errors.Wrapf(err, "ReadNewAndOldIssuesAndPrefetch/Read")
			return
		}
		if newAndOldIssues.NewIssue.GetIssueID() == "" {
			err = errors.New("newAndOldIssues.NewIssue.GetIssueID = ''")
			return
		}
		if newAndOldIssues.NewIssue.GetIssueID() != issueID {
			err = errors.Errorf(" newAndOldIssues.NewIssue.UserObjective.ID = %s != issueID = %s", newAndOldIssues.NewIssue.GetIssueID(), issueID)
			return
		}
		err = Prefetch(&newAndOldIssues.NewIssue, isShowingProgress)(DynamoDBConnection)
		if err != nil {
			err = errors.Wrapf(err, "getNewAndOldIssues/prefetch")
			return
		}
		newAndOldIssues.OldIssue = newAndOldIssues.NewIssue // we don't have the previous version of the entity
		err = errors.Wrap(err, "{ReadNewAndOldIssuesAndPrefetch}")
		return
	}
}

// Prefetch reads joined tables and puts related data into issue
func Prefetch(issueRef *Issue, isShowingProgress bool) func(DynamoDBConnection) (err error) {
	return func(DynamoDBConnection DynamoDBConnection) (err error) {
		if isShowingProgress {
			// 	objectiveProgress := LatestProgressUpdateByObjectiveID(issue.UserObjective.ID)
			issueRef.PrefetchedData.Progress, err = IssueProgressReadAll(issueRef.UserObjective.ID, 0)(DynamoDBConnection)
			log.Printf("Prefetch: len(Progress)==%d", len(issueRef.PrefetchedData.Progress))
			if err != nil {
				return
			}
		}
		return PrefetchIssueWithoutProgress(issueRef)(DynamoDBConnection)
	}
}

// PrefetchIssueWithoutProgress loads issue information ignoring context
func PrefetchIssueWithoutProgress(issueRef *Issue) func(DynamoDBConnection) (err error) {
	return func(conn DynamoDBConnection) (err error) {
		if !utilsUser.IsSpecialOrEmptyUserID(issueRef.UserObjective.AccountabilityPartner) {
			issueRef.PrefetchedData.AccountabilityPartner, err =
				utilsUser.DAOFromConnection(conn).
					Read(issueRef.UserObjective.AccountabilityPartner)
			if err != nil {
				return
			}
		}

		if issueRef.StrategyAlignmentEntityID != "" {
			switch issueRef.StrategyAlignmentEntityType {
			case userObjective.ObjectiveStrategyObjectiveAlignment:
				var sos []strategyObjective.StrategyObjective
				sos, err = StrategyObjectiveReadOrEmpty(issueRef.StrategyAlignmentEntityID)(conn)
				if len(sos) > 0 {
					issueRef.PrefetchedData.AlignedCapabilityObjective = sos[0]
				}
			case userObjective.ObjectiveStrategyInitiativeAlignment:
				var inis []strategyInitiative.StrategyInitiative
				inis, err = StrategyInitiativeReadOrEmpty(conn.PlatformID, issueRef.StrategyAlignmentEntityID)(conn)
				if len(inis) > 0 {
					issueRef.PrefetchedData.AlignedCapabilityInitiative = inis[0]
				}
			case userObjective.ObjectiveCompetencyAlignment:
				issueRef.PrefetchedData.AlignedCompetency, err = adaptiveValue.Read(issueRef.StrategyAlignmentEntityID)(conn)
			}
			if err != nil {
				log.Printf("IGNORE ERROR Couldn't prefetch %s id=%s due to an error: %+v", issueRef.StrategyAlignmentEntityType, issueRef.StrategyAlignmentEntityID, err)
				err = nil
			}
		}
		itype := issueRef.GetIssueType()
		switch itype {
		case IDO:
			// see above - prefetched data
		case SObjective:
			// already prefetched?
			if len(issueRef.StrategyObjective.CapabilityCommunityIDs) > 0 {
				capCommID := issueRef.StrategyObjective.CapabilityCommunityIDs[0]
				var comms []capabilityCommunity.CapabilityCommunity
				comms, err = CapabilityCommunityReadOrEmpty(capCommID)(conn)
				if len(comms) > 0 {
					issueRef.PrefetchedData.AlignedCapabilityCommunity = comms[0]
				}
			}
			// splits := strings.Split(issueRef.UserObjective.ID, "_")
			// if len(splits) == 2 {
			// 	soID := splits[0]
			// 	capCommID := splits[1]
			// 	issueRef.PrefetchedData.AlignedCapabilityObjective, err = StrategyObjectiveDAO.Read(teamID, soID)
			// 	if err != nil { return }
			// 	issueRef.PrefetchedData.AlignedCapabilityCommunity, err = CapabilityCommunityDAO.Read(teamID, capCommID)
			// } else {
			// 	issueRef.PrefetchedData.AlignedCapabilityObjective, err = StrategyObjectiveDAO.Read(teamID, issueRef.UserObjective.ID)
			// }
		case Initiative:
			initCommID := issueRef.StrategyInitiative.InitiativeCommunityID
			if initCommID != "" {
				var comms []strategyInitiativeCommunity.StrategyInitiativeCommunity
				comms, err = strategyInitiativeCommunity.ReadOrEmpty(initCommID, conn.PlatformID)(conn)
				for _, comm := range comms {
					issueRef.PrefetchedData.AlignedInitiativeCommunity = comm
				}
				if err != nil {
					return
				}
			}
			capObjID := issueRef.StrategyInitiative.CapabilityObjective
			if capObjID != "" {
				var sos []strategyObjective.StrategyObjective
				sos, err = StrategyObjectiveReadOrEmpty(capObjID)(conn)
				if len(sos) > 0 {
					issueRef.PrefetchedData.AlignedCapabilityObjective = sos[0]
				}
			}
		default:
		}
		err = errors.Wrap(err, "{prefetch}")
		return
	}
}

func PrefetchManyIssuesWithoutProgress(issues []Issue) func(DynamoDBConnection) (prefetchedIssues []Issue, err error) {
	return func(DynamoDBConnection DynamoDBConnection) (prefetchedIssues []Issue, err error) {
		for _, issue := range issues {
			err = PrefetchIssueWithoutProgress(&issue)(DynamoDBConnection)
			if err != nil {
				return
			}
			prefetchedIssues = append(prefetchedIssues, issue)
		}
		return
	}
}
