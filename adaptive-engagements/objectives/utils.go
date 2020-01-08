package objectives

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	business_time "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	utc, _ = time.LoadLocation("UTC")
	today  = business_time.Today(utc)

	CoachNotNeededOption ui.PlainText = "Coach not needed"
	RequestACoachOption  ui.PlainText = "Request a coach"
)

// USER_OBJECTIVES_USER_ID_INDEX
// AllUserObjectives retrieves all user objectives based on the objective type: stragegy/individual
// For strategy type, it list all capability objecives and initiatives that the user is a part of
func AllUserObjectives(userID string, userObjectivesTable, userObjectivesUserIdIndex string,
	typ models.DevelopmentObjectiveType, completed int) []models.UserObjective {
	var ops []models.UserObjective
	var filteredOps []models.UserObjective
	dns := common.DeprecatedGetGlobalDns()
	err := dns.Dynamo.QueryTableWithIndex(userObjectivesTable, awsutils.DynamoIndexExpression{
		IndexName: userObjectivesUserIdIndex,
		Condition: "user_id = :u AND completed = :c",
		Attributes: map[string]interface{}{
			":u": userID,
			":c": completed,
		},
	}, map[string]string{}, true, -1, &ops)
	core.ErrorHandler(err, dns.Namespace, fmt.Sprintf("Could not query %s index", userObjectivesUserIdIndex))
	for _, each := range ops {
		// Filtering to show only the objectives that match the `typ` provided
		if each.Type == typ {
			filteredOps = append(filteredOps, each)
		}
	}
	return filteredOps
}

// USER_OBJECTIVES_USER_ID_INDEX
// Used
func AllUserObjectivesWithProgress(userID string, userObjectivesTable, userObjectivesUserIdIndex string,
	userObjectivesProgressTable, userObjectivesProgressIdIndex string,
	typ models.DevelopmentObjectiveType, completed int) []models.UserObjectiveWithProgress {
	var res []models.UserObjectiveWithProgress
	objs := AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIdIndex, typ, completed)
	for _, each := range objs {
		ops := ObjectiveLatestProgress(userObjectivesProgressTable, each.ID, userObjectivesProgressIdIndex)
		res = append(res, models.UserObjectiveWithProgress{Objective: each, Progress: ops})
	}
	return res
}

// USER_OBJECTIVES_USER_ID_INDEX
func AllUserObjectivesWithProgressWithinPeriod(userID, userObjectivesTable, userObjectivesUserIdIndex, userObjectivesProgressTable string,
	start, end string, typ models.DevelopmentObjectiveType,
	completed int) []models.UserObjectiveWithProgress {
	var res []models.UserObjectiveWithProgress
	objs := AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIdIndex, typ, completed)
	for _, each := range objs {
		ops := ObjectiveProgressInPeriod(userObjectivesProgressTable, each.ID, start, end)
		res = append(res, models.UserObjectiveWithProgress{Objective: each, Progress: ops})
	}
	return res
}

// USER_OBJECTIVES_PROGRESS_ID_INDEX
func ObjectiveLatestProgress(table string, id, objProgressIdIndex string) []models.UserObjectiveProgress {
	var ops []models.UserObjectiveProgress
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: objProgressIdIndex,
		Condition: "id = :i",
		Attributes: map[string]interface{}{
			":i": id,
		},
		// Take the last updated progress
	}, map[string]string{}, false, 1, &ops)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s index", objProgressIdIndex))
	return ops
}

func ObjectiveProgressInPeriod(table string, id, start, end string) []models.UserObjectiveProgress {
	var ops []models.UserObjectiveProgress
	queryExpr := "id = :i AND created_on BETWEEN :t1 AND :t2"
	params := map[string]*dynamodb.AttributeValue{
		":i": {
			S: aws.String(id),
		},
		":t1": {
			S: aws.String(start),
		},
		":t2": {
			S: aws.String(end),
		},
	}
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithExpr(table, queryExpr, map[string]string{}, params, true, -1, &ops)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query for objective progress within a period"))
	return ops
}

func OpenObjectives(table, userId string, objCompletedIndex string, dns common.DynamoNamespace) (ops []models.UserObjective) {
	err := dns.Dynamo.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: objCompletedIndex,
		Condition: "user_id = :u AND completed = :c",
		Attributes: map[string]interface{}{
			":u": aws.String(userId),
			":c": aws.Bool(false),
		},
	}, map[string]string{}, true, -1, &ops)
	core.ErrorHandler(err, dns.Namespace, fmt.Sprintf("Could not query %s index", objCompletedIndex))
	return ops
}

func idParams(id string) map[string]*dynamodb.AttributeValue {
	params := map[string]*dynamodb.AttributeValue{
		"id": dynString(id),
	}
	return params
}

func dynString(str string) *dynamodb.AttributeValue {
	attr := dynamodb.AttributeValue{S: aws.String(str)}
	return &attr
}

// UserObjectiveById returns user objective by ID.
// returns nil if not found.
func UserObjectiveById(userObjectiveTableName string, ID string, dns common.DynamoNamespace) (ref *models.UserObjective) {
	op := models.UserObjective{}
	params := idParams(ID)
	err2 := dns.Dynamo.GetItemFromTable(userObjectiveTableName, params, &op)
	if err2 == nil {
		ref = &op
	} else {
		log.Printf("UserObjectiveById: Could not GetItemFromTable(table=%s, ids=%v): %v", userObjectiveTableName, params, err2)
		ref = nil
	}
	return
}

// Used
// Get user objectives by type, individual or strategy
func UserObjectivesByType(userID string, userObjectivesTable, userObjectivesTypeIndex string,
	typ models.DevelopmentObjectiveType, completed int) []models.UserObjective {
	var op []models.UserObjective
	var rels []models.UserObjective
	err := common.DeprecatedGetGlobalDns().Dynamo.QueryTableWithIndex(userObjectivesTable, awsutils.DynamoIndexExpression{
		IndexName: userObjectivesTypeIndex,
		Condition: "user_id = :u and #type = :t",
		Attributes: map[string]interface{}{
			":u": userID,
			":t": string(typ),
		},
	}, map[string]string{"#type": "type"}, true, -1, &rels)
	core.ErrorHandler(err, common.DeprecatedGetGlobalDns().Namespace, fmt.Sprintf("Could not query %s table on %s index",
		userObjectivesTable, userObjectivesTypeIndex))
	for _, each := range rels {
		if each.Completed == completed {
			op = append(op, each)
		}
	}
	return op
}

func readableDate(dateStr, prefix, namespace string) string {
	res, err := common.DateFormat.ChangeLayout(dateStr, core.USDateLayout)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not convert %s to USDateLayout", dateStr))
	return fmt.Sprintf("%s (%s)", prefix, res)
}

func currentSelectionDate(namespace, selected string) (dates []models.KvPair) {
	// Adding the current selection when present
	if selected != "" {
		_, err := common.DateFormat.Parse(selected)
		if err == nil {
			dates = append(dates, models.KvPair{
				Key:   readableDate(selected, "Current Selection", namespace),
				Value: selected,
			})
		}
	}
	return
}

func commonDates(namespace string, selected string) []models.KvPair {
	dates := make([]models.KvPair, 0)

	ninetyDays := today.AddTime(0, 0, 90).DateToString(string(common.DateFormat))
	dates = append(dates, models.KvPair{
		Key:   readableDate(ninetyDays, "Ninety Days", namespace),
		Value: ninetyDays,
	})

	sixtyDays := today.AddTime(0, 0, 60).DateToString(string(common.DateFormat))
	dates = append(dates, models.KvPair{
		Key:   readableDate(sixtyDays, "Sixty Days", namespace),
		Value: sixtyDays,
	})

	FortyFiveDays := today.AddTime(0, 0, 45).DateToString(string(common.DateFormat))
	dates = append(dates, models.KvPair{
		Key:   readableDate(FortyFiveDays, "Forty Five Days", namespace),
		Value: FortyFiveDays,
	})

	thirtyDays := today.AddTime(0, 0, 30).DateToString(string(common.DateFormat))
	dates = append(dates, models.KvPair{
		Key:   readableDate(thirtyDays, "Thirty Days", namespace),
		Value: thirtyDays,
	})
	return dates
}

func endOfQuarterDate(namespace string) models.KvPair {
	endOfQuarter := strings.Split(today.GetLastDayOfQuarter().DateToString(string(common.DateFormat)), "T")[0]
	return models.KvPair{
		Key:   readableDate(endOfQuarter, "End of Quarter", namespace),
		Value: endOfQuarter,
	}
}
func endOfYearDate(namespace string) models.KvPair {
	endOfYear := business_time.NewDate(today.GetYear()+1, 1, 1).AddTime(0, 0, -1).
		DateToString(string(common.DateFormat))
	return models.KvPair{
		Key:   readableDate(endOfYear, "End of Year", namespace),
		Value: endOfYear,
	}
}

func yearFromNowDate(namespace string) models.KvPair {
	nextYear := today.AddTime(1, 0, 0).DateToString(string(common.DateFormat))
	return models.KvPair{
		Key:   readableDate(nextYear, "A year from now", namespace),
		Value: nextYear,
	}
}

func DevelopmentObjectiveDates(namespace, currentSelection string) []models.KvPair {
	var dates = currentSelectionDate(namespace, currentSelection)
	dates = append(dates, endOfQuarterDate(namespace))
	dates = append(dates, []models.KvPair{
		yearFromNowDate(namespace),
		endOfYearDate(namespace),
	}...)
	dates = append(dates, commonDates(namespace, currentSelection)...)
	return dates
}

func StrategyObjectiveDates(namespace, currentSelection string) []models.KvPair {
	var dates = currentSelectionDate(namespace, currentSelection)
	dates = append(dates, []models.KvPair{
		yearFromNowDate(namespace),
		endOfYearDate(namespace),
	}...)
	dates = append(dates, commonDates(namespace, currentSelection)...)
	return dates
}

func StrategyObjectiveDatesWithIndefiniteOption(namespace, currentSelection string) []models.KvPair {
	var dates = currentSelectionDate(namespace, currentSelection)
	dates = append(dates, []models.KvPair{
		endOfYearDate(namespace),
		yearFromNowDate(namespace),
	}...)
	dates = append(dates, models.KvPair{
		Key:   common.StrategyIndefiniteDateKey,
		Value: common.StrategyIndefiniteDateValue,
	})
	dates = append(dates, commonDates(namespace, currentSelection)...)
	return dates
}

func mapAdaptiveCommunityUsersToUserID(users []models.AdaptiveCommunityUser2) (userIDs []string) {
	for _, each := range users {
		userIDs = append(userIDs, each.UserId)
	}
	return
}

// IDOCoaches retrieves available coaches for a user
// It shows 'none' option indicating no coach is required
// It shows 'requested' option when a user wants to request uber-coach
func IDOCoaches(userID, platformID string,
	communityUsersTable, communityUsersCommunityIndex string,
	fetchUsers common.UserIDsToDisplayNames) []models.KvPair {
	// Get coaching community members
	commMembers := community.CommunityMembers(communityUsersTable, string(community.Coaching), platformID, communityUsersCommunityIndex)
	kvs := []models.KvPair{{Key: string(CoachNotNeededOption), Value: "none"}}
	// Showing "Request a Coach" option only when there is a coaching community
	if len(commMembers) > 0 { // Does this include adaptive bot name?
		kvs = append(kvs, models.KvPair{
			Key:   string(RequestACoachOption),
			Value: "requested"})
	}
	userIDs := mapAdaptiveCommunityUsersToUserID(commMembers)
	userIDs = filterStrings(userIDs, func(id string) bool { return userID != id })
	userOptions := fetchUsers(userIDs)
	kvs = append(kvs, userOptions...)
	return kvs
}

func filterStrings(in []string, predicate func(string) bool) (out []string) {
	for _, each := range in {
		if predicate(each) {
			out = append(out, each)
		}
	}
	return
}

// ObjectivesDueInDuration retrieves all user objectives that are due within a period based on the objective type: strategy/individual
func ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex string,
	typ models.DevelopmentObjectiveType, alignmentType models.AlignedStrategyType, ipDate business_time.Date,
	days int) []models.UserObjective {
	var op []models.UserObjective
	objs := AllUserObjectives(userID, userObjectivesTable, userObjectivesUserIndex, typ, 0)
	for _, each := range objs {
		endDate, err := business_time.DateFromYMDString(each.ExpectedEndDate)
		core.ErrorHandler(err, "adaptive-checks", fmt.Sprintf("Could not parse %s date", each.ExpectedEndDate))
		duration := endDate.DaysBetween(ipDate)
		if each.StrategyAlignmentEntityType == alignmentType && duration == days {
			op = append(op, each)
		}
	}
	return op
}

// StaleObjectivesDueInDuration retrieves all user objectives for which progress hasn't been added,
// based on the objective type: strategy/individual
func StaleObjectivesInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
	userObjectivesProgressTable string, fDay1, lDay1 business_time.Date, typ models.DevelopmentObjectiveType,
	alignmentTypes ...models.AlignedStrategyType) []models.UserObjective {
	fDay := fDay1.DateToString(string(common.DateFormat))
	lDay := lDay1.DateToString(string(common.DateFormat))
	objsWithProgress := AllUserObjectivesWithProgressWithinPeriod(userID, userObjectivesTable,
		userObjectivesUserIndex, userObjectivesProgressTable, fDay, lDay, typ, 0)
	var op []models.UserObjective

	for _, each := range objsWithProgress {
		containsType := false
		for _, eachType := range alignmentTypes {
			if each.Objective.StrategyAlignmentEntityType == eachType {
				containsType = true
			}
		}
		objCreatedDate, err := common.DateFormat.Parse(each.Objective.CreatedDate)
		if err == nil {
			if containsType && len(each.Progress) == 0 && objCreatedDate.Before(fDay1.DateToTimeMidnight()) {
				op = append(op, each.Objective)
			}
		}
	}
	return op
}

// Are there any open IDO's that exist for the user that are due in exactly 7 days
func IDOsDueInAWeek(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, models.ObjectiveNoStrategyAlignment, ipDate, 7)
}

// Are there any open IDO's that exist for the user that are
// NOT due within the week but are due exactly in 30 days
func IDOsDueInAMonth(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, models.ObjectiveNoStrategyAlignment, ipDate, 30)
}

// Are there any open IDO's that exist for the user that are
// NOT due within the month but are due in exactly 90 days
func IDOsDueInAQuarter(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex string) []models.UserObjective {
	return ObjectivesDueInDuration(userID, userObjectivesTable, userObjectivesUserIndex,
		models.IndividualDevelopmentObjective, models.ObjectiveNoStrategyAlignment, ipDate, 90)
}

// UserIDOsWithNoProgressInLastWeek returns IDOs that exist for the user that haven't been updated in last 7 days
func UserIDOsWithNoProgressInAWeek(userID string, ipDate business_time.Date, userObjectivesTable,
	userObjectivesUserIndex, userObjectivesProgressTable string) []models.UserObjective {
	aWeekBefore := ipDate.AddTime(0, 0, -7)
	fDay := aWeekBefore
	lDay := ipDate
	return StaleObjectivesInDuration(userID,
		userObjectivesTable, userObjectivesUserIndex, userObjectivesProgressTable,
		fDay, lDay, models.IndividualDevelopmentObjective, models.ObjectiveNoStrategyAlignment,
		models.ObjectiveStrategyObjectiveAlignment, models.ObjectiveStrategyInitiativeAlignment,
		models.ObjectiveCompetencyAlignment)
}
