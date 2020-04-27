package issues

import (
	"github.com/adaptiveteam/adaptive/daos/strategyInitiative"
	"strings"
	"log"
	"github.com/adaptiveteam/adaptive/business-time"
	"fmt"
	"time"
	utilsIssues "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
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
			issues = append(issues, Issue{UserObjective: uo})
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
			var moreIssues []Issue
			// log.Printf("found strategy objective. uo.ID=%s\n", uo.ID)
			moreIssues, err = FillInSObjective(uo)(conn)
			if err != nil { 
				return 
			}
			if len(moreIssues) == 0 {
				log.Printf("ERROR: couldn't find StrategyObjective for ID=%s", uo.ID)
			}
			issues = append(issues, moreIssues...)
		}
		return
	}
}

func FillInSObjective(userObj userObjective.UserObjective) func (conn DynamoDBConnection) (issues [] Issue, err error) {
	return func (conn DynamoDBConnection) (issues [] Issue, err error) {
		id := userObj.ID
		i := strings.Index(id, "_")
		if i >= 0 {
			log.Printf("WARN: ID has '_': %s\n", id)
			id = id[0:i]
		}
		var sos []models.StrategyObjective
		sos, err = utilsIssues.StrategyObjectiveReadOrEmpty(id)(conn)

		if err == nil {
			if len(sos) > 0 {
				issue := Issue{
					UserObjective: userObj,
					StrategyObjective: sos[0],
				}
				issues = []Issue{issue}
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
		var inis []strategyInitiative.StrategyInitiative
		inis, err = utilsIssues.StrategyInitiativeReadOrEmpty(conn.PlatformID, id)(conn)
		if len(inis) > 0 && err == nil {
			issue = Issue{
				UserObjective: userObj,
				StrategyInitiative: inis[0],
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
			issue, err = FillInInitiative(uo)(conn)
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
			var moreIssues []Issue
			moreIssues, err = FillInSObjective(uo)(conn)
			if err != nil { return }
			if len(moreIssues) == 0 {
				log.Printf("WARN Couldn't FillInSObjective for uo.ID=%s\n",uo.ID)
			}
			issues = append(issues, moreIssues...)
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
			issue, err = FillInInitiative(uo)(conn)
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

