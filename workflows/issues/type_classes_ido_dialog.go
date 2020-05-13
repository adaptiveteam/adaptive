package issues

import (
	"fmt"
	"time"

	objectives "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	adaptiveValue "github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// CreateDialog loads data and creates IDO dialog
func (IDOImpl) CreateDialog(w workflowImpl, ctx wf.EventHandlingContext, issue Issue) (survey ebm.AttachmentActionSurvey, err error) {
	defer w.recoverToErrorVar("IDOImpl.CreateDialog", &err) // because there are panics downstream
	namespace := "(IDOImpl)CreateDialog"
	allDates := objectives.DevelopmentObjectiveDates(namespace, issue.UserObjective.ExpectedEndDate)
	dates := convertKvPairToPlainTextOption(allDates)

	var allMembers []models.KvPair
	allMembers, err = IDOCoaches(ctx.Request.User.ID, issue.UserObjective.AccountabilityPartner)(w.DynamoDBConnection)
	if err != nil {
		return
	}
	w.AdaptiveLogger.Infof("Coaches of the team: %s (%d)", w.DynamoDBConnection.PlatformID, len(allMembers))
	for _, kv := range allMembers {
		w.AdaptiveLogger.Infof("Coach: %s -> %s", kv.Key, kv.Value)
	}
	coaches := convertKvPairToPlainTextOption(allMembers)

	var initiativesObjectiveAndCompetencies []ebm.AttachmentActionElementOptionGroup
	initiativesObjectiveAndCompetencies, err = LoadAndFormatInitiativesObjectiveAndCompetencies(w, ctx.Request.User.ID, 
		ctx.TeamID, issue.UserObjective)
	survey = IDOObjectiveSurvey(issue.UserObjective, coaches, dates, initiativesObjectiveAndCompetencies)
	return
}

// LoadAndFormatInitiativesObjectiveAndCompetencies loads InitiativesObjectiveAndCompetencies
func LoadAndFormatInitiativesObjectiveAndCompetencies(w workflowImpl,
	userID string,
	teamID models.TeamID,
	item userObjective.UserObjective) (initiativesObjectiveAndCompetencies []ebm.AttachmentActionElementOptionGroup, err error) {
	var competencies []adaptiveValue.AdaptiveValue
	competencies, err = CompetencyReadAll()(w.DynamoDBConnection)
	if err == nil {
		var initiativesForUser []models.StrategyInitiative
		var capabilityObjectives []models.StrategyObjective
		initiativesForUser, capabilityObjectives, err = LoadInitsAndObjectives(w, userID, teamID)
		if err == nil {
			i := formatInitiativesAsGroup(initiativesForUser)
			initiativesObjectiveAndCompetencies = append(initiativesObjectiveAndCompetencies, i...)
			o := formatObjectivesGroup(capabilityObjectives)
			initiativesObjectiveAndCompetencies = append(initiativesObjectiveAndCompetencies, o...)
			c := formatCompetenciesAsOptionGroup(competencies)
			w.AdaptiveLogger.Infof("Retrieved competencies for %s platform: %v", teamID, c)
			initiativesObjectiveAndCompetencies = append(initiativesObjectiveAndCompetencies, c...)
		}
	}
	return
}

// IDOObjectiveSurvey shows a form to create or modify an objective
func IDOObjectiveSurvey(item userObjective.UserObjective,
	coaches, dates []ebm.AttachmentActionElementPlainTextOption,
	initiativesAndObjectives []ebm.AttachmentActionElementOptionGroup) ebm.AttachmentActionSurvey {
	alignmentID := objectives.AlignmentIDFromAlignedStrategyType(item.StrategyAlignmentEntityType, item.StrategyAlignmentEntityID)
	return ebm.AttachmentActionSurvey{
		Title: "Objective",
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewTextBox(objectives.ObjectiveName, "Name", objectives.ObjectiveNamePlaceholder, ui.PlainText(item.Name)),
			ebm.NewTextArea(objectives.ObjectiveDescription, "Description", objectives.ObjectiveDescriptionPlaceholder, ui.PlainText(item.Description)),
			ebm.NewSimpleOptionsSelect(objectives.ObjectiveAccountabilityPartner, "Coach", ebm.EmptyPlaceholder, string(item.AccountabilityPartner), coaches...),
			ebm.NewSimpleOptionsSelect(objectives.ObjectiveEndDate, "Expected end date", ebm.EmptyPlaceholder, item.ExpectedEndDate, dates...),
			ebm.NewSimpleOptionGroupsSelect(objectives.ObjectiveStrategyAlignment, "Strategy Alignment", ebm.EmptyPlaceholder, alignmentID,
				initiativesAndObjectives...),
		},
	}
}

func convertKvPairToPlainTextOption(pairs []models.KvPair) (out []ebm.AttachmentActionElementPlainTextOption) {
	for _, p := range pairs {
		out = append(out, ebm.AttachmentActionElementPlainTextOption{Value: p.Value, Label: ui.PlainText(p.Key)})
	}
	return
}

// LoadInitsAndObjectives returns initiatives and objectives
func LoadInitsAndObjectives(
	w workflowImpl,
	userID string,
	teamID models.TeamID) (
	initiativesForUser []models.StrategyInitiative,
	capabilityObjectives []models.StrategyObjective,
	err error) {
	initiativesForUser, err = SelectFromInitiativesJoinUserCommunityWhereUserID(userID)(w.DynamoDBConnection)
	if err == nil {
		capabilityObjectives, err = SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID(userID)(w.DynamoDBConnection)
		if err == nil {

			initiativeRelatedCapabilityObjectiveIDs := mapInitiativesToString(initiativesForUser, func(i models.StrategyInitiative) string { return i.CapabilityObjective })
			capabilityObjectiveIDs := mapObjectivesToString(capabilityObjectives, func(i models.StrategyObjective) string { return i.ID })
			objectivesIDsFromInitiativesNotInOptions := core.InBButNotA(capabilityObjectiveIDs, initiativeRelatedCapabilityObjectiveIDs)

			fmt.Println(capabilityObjectiveIDs)
			fmt.Println(initiativeRelatedCapabilityObjectiveIDs)
			fmt.Printf("### objectivesIDsFromInitiativesNotInOptions: %v\n", objectivesIDsFromInitiativesNotInOptions)
			var iObjs []models.StrategyObjective
			iObjs, err = SelectFromIssuesWhereTypeAndUserIDStrategyObjectives(objectivesIDsFromInitiativesNotInOptions)(w.DynamoDBConnection)
			if err == nil {
				capabilityObjectives = append(capabilityObjectives, iObjs...)
			}
		}
	}
	return
}

// ExtractFromContext extracts UserObjective from the context.
func (IDOImpl) ExtractFromContext(ctx wf.EventHandlingContext, id string, updated bool, oldIssue Issue) (newIssue Issue) {
	form := ctx.Request.DialogSubmissionCallback.Submission
	var issueID string
	issueID, updated = ctx.Data[issueIDKey]
	if !updated {
		issueID = core.Uuid()
	}
	userID := ctx.Request.User.ID
	objName := form[objectives.ObjectiveName]
	objDescription := form[objectives.ObjectiveDescription]
	partner := form[objectives.ObjectiveAccountabilityPartner]
	endDate := form[objectives.ObjectiveEndDate]
	strategyEntityID := form[objectives.ObjectiveStrategyAlignment]
	// Get the alignment type for the aligned objective
	alignment, alignmentID := getAlignedStrategyTypeFromStrategyEntityID(strategyEntityID)

	newIssue.UserObjective = userObjective.UserObjective{
		ID:                          issueID,
		UserID:                      userID,
		CreatedBy:                   userID,
		Name:                        objName,
		Description:                 objDescription,
		AccountabilityPartner:       partner,
		ObjectiveType:               userObjective.IndividualDevelopmentObjective,
		StrategyAlignmentEntityID:   alignmentID,
		StrategyAlignmentEntityType: alignment,
		ExpectedEndDate:             endDate,
		PlatformID:                  ctx.TeamID.ToPlatformID(),
	}
	if updated {
		newIssue.UserObjective.Year = oldIssue.UserObjective.Year
		newIssue.UserObjective.Quarter = oldIssue.UserObjective.Quarter
		newIssue.UserObjective.CreatedDate = oldIssue.UserObjective.CreatedDate
		newIssue.UserObjective.CreatedBy = oldIssue.UserObjective.CreatedBy
	} else {
		year, quarter := core.CurrentYearQuarter()
		newIssue.UserObjective.Year = year
		newIssue.UserObjective.Quarter = quarter
		newIssue.UserObjective.CreatedDate = core.ISODateLayout.Format(time.Now())
	}
	newIssue.UserObjective.ModifiedBy = userID

	return
}

func (i IDOImpl) View(w workflowImpl, isShowingDetails, isShowingProgress bool,
	newAndOldIssues NewAndOldIssues,
) (fields []ebm.AttachmentField) {
	viewState := engIssues.ViewState{IsShowingDetails: isShowingDetails, IsShowingProgress: isShowingProgress}
	fields = engIssues.OrdinaryViewFields(newAndOldIssues, viewState)
	return
}
