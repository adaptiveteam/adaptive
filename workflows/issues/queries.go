package issues

import (
	"strings"
	"log"
	"github.com/adaptiveteam/adaptive/business-time"
	"fmt"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

// StaleObjectivesQueryIDO queries stale IDOs
func StaleObjectivesQueryIDO(userID string) IssueQuery { 
	return func (conn DynamoDBConnection) (issues []Issue, err error) {
		defer recoverToErrorVar("StaleObjectivesQuery IDO", &err)
		date := business_time.Today(time.UTC)
		fmt.Printf("UserIDOsWithNoProgressInAWeek(%s, %s, %s, %s, %s)\n",
			userID, date,
			userObjectivesTableName(conn.ClientID),
			userObjectivesIDIndex,
			userObjectivesProgressTableName(conn.ClientID))
//UserIDOsWithNoProgressInAWeek(UJ0SX0G9X, {%!s(int=2019) %!s(int=11) %!s(int=28)}, ivan_user_objective, IDIndex, ivan_user_objectives_progress)[31mERROR:[adaptive.dynamo.user-objectives] ValidationException: Query key condition not supported
		userObjs := objectives.UserIDOsWithNoProgressInAWeek(
			userID,
			date,
			userObjectivesTableName(conn.ClientID),
			userObjectivesUserIDIndex,
			userObjectivesProgressTableName(conn.ClientID))
		for _, uo := range userObjs {
			fmt.Printf("Converting UserObjective to Issue: id=%s, name=%s\n", uo.ID, uo.Name)
			issues = append(issues, Issue{UserObjective: convertModelsUserObjectiveToUserObjective(uo)})
		}
		return
	}
}

// StaleObjectivesQuerySObjectives queries stale SObjective-issues
func StaleObjectivesQuerySObjectives(userID string) IssueQuery { 
	return func (conn DynamoDBConnection) (issues []Issue, err error) {
		defer recoverToErrorVar("StaleObjectivesQuery SObjectives", &err)
		date := business_time.Today(time.UTC)
		userObjs := strategy.UserCapabilityObjectivesWithNoProgressInAMonth(
			userID,
			date,
			userObjectivesTableName(conn.ClientID),
			userObjectivesUserIDIndex,
			userObjectivesProgressTableName(conn.ClientID), 0)
		for _, uo := range userObjs {
			var issue Issue
			// log.Printf("found strategy objective. uo.ID=%s\n", uo.ID)
			issue, err = FillInSObjective(convertModelsUserObjectiveToUserObjective(uo))(conn)
			if err != nil { return }
			issues = append(issues, issue)
		}
		return
	}
}

func FillInSObjective(userObj userObjective.UserObjective) func (conn DynamoDBConnection) (issue Issue, err error) {
	return func (conn DynamoDBConnection) (issue Issue, err error) {
		id := userObj.ID
		i := strings.Index(id, "_")
		if i >= 0 {
			log.Printf("WARN: ID has '_': %s\n", id)
			id = id[0:i]
		}
		var so models.StrategyObjective
		so, err = StrategyObjectiveRead(id)(conn)

		if err == nil {
			issue = Issue{
				UserObjective: userObj,
				StrategyObjective: so,
			}
		}
		return
	}
}

func FillInInitiative(userObj userObjective.UserObjective) func (conn DynamoDBConnection) (issue Issue, err error) {
	return func (conn DynamoDBConnection) (issue Issue, err error) {
		id := userObj.ID
		i := strings.Index(id, "_")
		if i >= 0 {
			log.Printf("WARN: ID has '_': %s\n", id)
			id = id[0:i]
		}
		var ini models.StrategyInitiative
		ini, err = StrategyInitiativeRead(id)(conn)
		if err == nil {
			issue = Issue{
				UserObjective: userObj,
				StrategyInitiative: ini,
			}
		}
		return
	}
}

// StaleObjectivesQueryInitiative queries stale Initiative-issues
func StaleObjectivesQueryInitiative(userID string) IssueQuery { 
	return func (conn DynamoDBConnection) (issues []Issue, err error) {
		defer recoverToErrorVar("StaleObjectivesQuery Initiative", &err)
		date := business_time.Today(time.UTC)
		userObjs := strategy.UserInitiativesWithNoProgressInAWeek(
			userID,
			date,
			userObjectivesTableName(conn.ClientID),
			userObjectivesUserIDIndex,
			userObjectivesProgressTableName(conn.ClientID), 0)
		for _, uo := range userObjs {
			var issue Issue
			issue, err = FillInInitiative(convertModelsUserObjectiveToUserObjective(uo))(conn)
			if err != nil { return }
			issues = append(issues, issue)
		}
		return
	}
}

// StaleObjectivesQuery queries stale issues
func StaleObjectivesQuery(ctx wf.EventHandlingContext) (query IssueQuery) { 
	itype := getIssueTypeFromContext(ctx)
	switch itype {
	case IDO:
		query = StaleObjectivesQueryIDO(ctx.Request.User.ID)
	case SObjective:
		query = StaleObjectivesQuerySObjectives(ctx.Request.User.ID)
	case Initiative:
		query = StaleObjectivesQueryInitiative(ctx.Request.User.ID)
	}
	return
}

func convertModelsUserObjectiveToUserObjective(muo models.UserObjective) userObjective.UserObjective {
	createdAtTime, err2 := core.ISODateLayout.Parse(muo.CreatedDate)
	if err2 != nil {
		fmt.Printf("INVALID: Couldn't convert date %+v\n", err2)
		createdAtTime = time.Now()
	}
	createdAt := core.ISODateLayout.Format(createdAtTime)
	return userObjective.UserObjective{
		ID:                            muo.ID,
		UserID:                        muo.UserID,
		Name:                          muo.Name,
		PlatformID:                    muo.PlatformID,
		Description:                   muo.Description,
		AccountabilityPartner:         muo.AccountabilityPartner,
		Accepted:                      muo.Accepted,
		ObjectiveType:                 userObjective.DevelopmentObjectiveType(muo.ObjectiveType),
		StrategyAlignmentEntityID:     muo.StrategyAlignmentEntityID,
		StrategyAlignmentEntityType:   userObjective.AlignedStrategyType(muo.StrategyAlignmentEntityType),
		Quarter:                       muo.Quarter,
		Year:                          muo.Year,
		CreatedDate:                   muo.CreatedDate,
		ExpectedEndDate:               muo.ExpectedEndDate,
		Completed:                     muo.Completed,
		PartnerVerifiedCompletion:     muo.PartnerVerifiedCompletion,
		CompletedDate:                 muo.CompletedDate,
		PartnerVerifiedCompletionDate: muo.PartnerVerifiedCompletionDate,
		Comments:                      muo.Comments,
		Cancelled:                     muo.Cancelled,
		CreatedAt:                     createdAt,
		ModifiedAt:                    "",
	}
}

// AdvocacyIssuesQuery queries issues for which the current user is the advocate.
func AdvocacyIssuesQuery(ctx wf.EventHandlingContext) (query IssueQuery) { 
	itype := getIssueTypeFromContext(ctx)
	switch itype {
	case IDO:
		query = AdvocacyIssuesQueryIDO(ctx.Request.User.ID)
	case SObjective:
		query = AdvocacyIssuesQuerySObjectives(ctx.Request.User.ID)
	case Initiative:
		query = AdvocacyIssuesQueryInitiative(ctx.Request.User.ID)
	}
	return
}

// AdvocacyIssuesQueryIDO queries IDO-issues for which the current user is the advocate
func AdvocacyIssuesQueryIDO(userID string) IssueQuery { 
	return func (conn DynamoDBConnection)(issues []Issue, err error) {
		defer recoverToErrorVar("AdvocacyIssuesQueryIDO IDO", &err)
		// NB! This is implemented in View Coachee IDOs in coaching lambda.
		log.Println("AdvocacyIssuesQueryIDO is implemented in View Coachee IDOs in coaching lambda.")
		return
	}
}

// AdvocacyIssuesQuerySObjectives queries SObjectives-issues for which the current user is the advocate
func AdvocacyIssuesQuerySObjectives(userID string) IssueQuery { 
	return func (conn DynamoDBConnection)(issues []Issue, err error) {
		defer recoverToErrorVar("AdvocacyIssuesQuery SObjectives", &err)
		userObjs := strategy.UserAdvocacyObjectives(
			userID,
			userObjectivesTableName(conn.ClientID),
			userObjectivesTypeIndex, 0)
		for _, uo := range userObjs {
			var issue Issue
			issue, err = FillInSObjective(convertModelsUserObjectiveToUserObjective(uo))(conn)
			if err != nil { return }
			issues = append(issues, issue)
		}
		return
	}
}

// AdvocacyIssuesQueryInitiative queries Initiative-issues for which the current user is the advocate
func AdvocacyIssuesQueryInitiative(userID string) IssueQuery { 
	return func (conn DynamoDBConnection)(issues []Issue, err error) {
		defer recoverToErrorVar("AdvocacyIssuesQuery Initiative", &err)
		fmt.Printf("UserAdvocacyInitiatives(%s, %s, %s, %d)\n",
			userID,
			userObjectivesTableName(conn.ClientID), //"UserIDTypeIndex"
			userObjectivesTypeIndex, 0) // UserAdvocacyInitiatives(UJ0SX0G9X, ivan_user_objective, UserIDTypeIndex, 0)
										// UserAdvocacyInitiatives(UJ0SX0G9X, ivan_user_objective, UserIDTypeIndex, 0)
		userObjs := strategy.UserAdvocacyInitiatives(
			userID,
			userObjectivesTableName(conn.ClientID),
			userObjectivesTypeIndex, 0)
		fmt.Printf("len(userObjs): %d\n", len(userObjs))
		for _, uo := range userObjs {
			var issue Issue
			issue, err = FillInInitiative(
				convertModelsUserObjectiveToUserObjective(uo))(conn)
			if err != nil { return }
			issues = append(issues, issue)
		}
		fmt.Printf("len(issues): %d\n", len(issues))
		return
	}
}

/*

	case strategy.ViewCommunityAdvocateObjectives:
		// List objectives for which the user is an advocate for
		stratObjs := objectives.AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIdIndex,
			models.StrategyDevelopmentObjective, 0)

*/

