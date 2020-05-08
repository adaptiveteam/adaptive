package issues

import (
	"fmt"
	"log"
	"time"

	engIssues "github.com/adaptiveteam/adaptive/adaptive-engagements/issues"
	objectives "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"

	// userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

func (i InitiativeImpl) View(w workflowImpl, isShowingDetails, isShowingProgress bool,
	newAndOldIssues NewAndOldIssues,
) (fields []ebm.AttachmentField) {
	viewState := engIssues.ViewState{IsShowingDetails: isShowingDetails, IsShowingProgress: isShowingProgress}
	fields = engIssues.OrdinaryViewFields(newAndOldIssues, viewState)
	return
}

const (
	InitiativeCommunityLabel ui.PlainText = "Initiative Community"
	RelatedObjectiveLabel    ui.PlainText = "Related Objective"
	DefinitionOfVictoryLabel ui.PlainText = "Definition of Victory"
	BudgetLabel              ui.PlainText = "Budget"

	InitiativeNameLabel        ui.PlainText = "Initiative Name"
	InitiativeDescriptionLabel ui.PlainText = "Initiative Description"
	InitiativeVictoryLabel     ui.PlainText = "Definition of Victory"

	InitiativeAdvocateLabel            = "Advocate"
	InitiativeBudgetLabel              = "Budget($) in the following format: 1234.56"
	InitiateEndDateLabel               = "Time to work on this"
	InitiativeCapabilityObjectiveLabel = "Related Capability Objective"

	InitiativeName                    = "initiative_name"
	InitiativeDescriptionName         = "initiative_description"
	InitiativeVictoryName             = "definition_of_victory"
	InitiativeAdvocateName            = "advocate"
	InitiativeBudgetName              = "initiative_budget_name"
	InitiateEndDateName               = "time_to_work_on_this"
	InitiativeCapabilityObjectiveName = "initiative_capability_objective"
)

// CreateDialog returns a survey for the given initiative issue.
func (InitiativeImpl) CreateDialog(w workflowImpl, ctx wf.EventHandlingContext, issue Issue) (survey ebm.AttachmentActionSurvey, err error) {
	defer w.recoverToErrorVar("InitiativeImpl.CreateDialog", &err) // because there are panics downstream
	capCommID := ctx.Data[capCommIDKey]
	commID := community.AdaptiveCommunity(fmt.Sprintf("%s:%s", community.Capability, capCommID))
	userID := ctx.Request.User.ID

	var commMembers, strMembers []models.KvPair
	var objs []models.StrategyObjective

	objs, err = SelectFromStrategyObjectivesWhenUserIsInStrategyUnionSelectFromStrategyObjectivesJoinCapabilityCommunitiesWhereUserID(userID)(w.DynamoDBConnection)
	if err == nil {
		strMembers, err = SelectKvPairsFromCommunityJoinUsers(community.Strategy)(w.DynamoDBConnection)
		if err == nil {
			commMembers, err = SelectKvPairsFromCommunityJoinUsers(commID)(w.DynamoDBConnection)
			if err == nil {
				allMembers := append(strMembers, commMembers...)
				allMembers = removeDuplicates(allMembers)
				allDates := objectives.StrategyObjectiveDatesWithIndefiniteOption("InitiativeImpl CreateDialog", issue.UserObjective.ExpectedEndDate)

				availableObjectives := convertStrategyObjectivesToPlainTextOption(objs)
				advocates := convertKvPairToPlainTextOption(allMembers)
				dates := convertKvPairToPlainTextOption(allDates)

				survey = InitiativeSurvey(issue.StrategyInitiative, availableObjectives, advocates, dates)
			}
		}
	}
	return
}

// InitiativeSurvey shows a form to create or modify an objective
func InitiativeSurvey(item models.StrategyInitiative,
	availableObjectives, advocates, dates []ebm.AttachmentActionElementPlainTextOption) ebm.AttachmentActionSurvey {
	return ebm.AttachmentActionSurvey{
		Title: "Initiative",
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewTextBox(InitiativeName, InitiativeNameLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Name)),
			ebm.NewTextArea(InitiativeDescriptionName, InitiativeDescriptionLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Description)),
			ebm.NewTextArea(InitiativeVictoryName, InitiativeVictoryLabel, ebm.EmptyPlaceholder, ui.PlainText(item.DefinitionOfVictory)),
			ebm.NewSimpleOptionsSelect(InitiativeAdvocateName, InitiativeAdvocateLabel, ebm.EmptyPlaceholder, string(item.Advocate), advocates...),
			// TODO: Budget: ElemSubtype: ebm.AttachmentActionTextElementNumberType,
			ebm.NewTextArea(InitiativeBudgetName, InitiativeBudgetLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Budget)),
			ebm.NewSimpleOptionsSelect(InitiateEndDateName, InitiateEndDateLabel, ebm.EmptyPlaceholder, item.ExpectedEndDate, dates...),
			ebm.NewSimpleOptionsSelect(InitiativeCapabilityObjectiveName, InitiativeCapabilityObjectiveLabel, ebm.EmptyPlaceholder, string(item.CapabilityObjective), availableObjectives...),
		},
	}
}

// ExtractFromContext extracts information about initiative from context.
func (i InitiativeImpl) ExtractFromContext(ctx wf.EventHandlingContext, id string, updated bool, oldIssue Issue) (newIssue Issue) {
	form := ctx.Request.DialogSubmissionCallback.Submission
	newIssue.StrategyInitiative.Name = form[InitiativeName]
	// newIssue.StrategyInitiative.Type = models.StrategyInitiativeType(form[SObjectiveType])
	newIssue.StrategyInitiative.Description = form[InitiativeDescriptionName]
	newIssue.StrategyInitiative.DefinitionOfVictory = form[InitiativeVictoryName]
	newIssue.StrategyInitiative.Budget = form[InitiativeBudgetName]
	newIssue.StrategyInitiative.CapabilityObjective = form[InitiativeCapabilityObjectiveName]
	// newIssue.StrategyInitiative.Targets = form[SObjectiveTargets]
	newIssue.StrategyInitiative.Advocate = form[InitiativeAdvocateName]
	newIssue.StrategyInitiative.ExpectedEndDate = form[InitiateEndDateName]
	newIssue.StrategyInitiative.PlatformID = ctx.TeamID.ToPlatformID()
	now := core.TimestampLayout.Format(time.Now())
	newIssue.UserObjective.ModifiedAt = now
	if updated {
		newIssue.StrategyInitiative.ID = oldIssue.StrategyInitiative.ID
		newIssue.StrategyInitiative.CreatedBy = oldIssue.StrategyInitiative.CreatedBy
		newIssue.StrategyInitiative.CreatedAt = oldIssue.StrategyInitiative.CreatedAt
		newIssue.StrategyInitiative.InitiativeCommunityID = oldIssue.StrategyInitiative.InitiativeCommunityID
	} else {
		newIssue.StrategyInitiative.ID = core.Uuid()
		newIssue.StrategyInitiative.CreatedBy = ctx.Request.User.ID
		newIssue.StrategyInitiative.CreatedAt = now
		initCommID := ctx.Data[initCommIDKey]
		if initCommID == "" {
			log.Printf("!Invalid state - there is no initCommID in the context")
		}
		newIssue.StrategyInitiative.InitiativeCommunityID = initCommID

	}
	newIssue.UserObjective.ID = newIssue.StrategyInitiative.ID

	newIssue.UserObjective.Name = newIssue.StrategyInitiative.Name
	newIssue.UserObjective.Description = newIssue.StrategyInitiative.Description
	newIssue.UserObjective.UserID = newIssue.StrategyInitiative.Advocate
	newIssue.UserObjective.AccountabilityPartner = newIssue.StrategyInitiative.CreatedBy
	newIssue.UserObjective.Accepted = 1 // since it is created by the same person
	newIssue.UserObjective.Name = newIssue.StrategyInitiative.Name
	newIssue.UserObjective.PlatformID = ctx.TeamID.ToPlatformID()
	newIssue.UserObjective.CreatedAt = newIssue.StrategyInitiative.CreatedAt
	newIssue.UserObjective.CreatedDate = core.NormalizeDate(newIssue.StrategyInitiative.CreatedAt)
	newIssue.UserObjective.ObjectiveType = userObjective.StrategyDevelopmentInitiative
	newIssue.UserObjective.StrategyAlignmentEntityType = userObjective.ObjectiveStrategyInitiativeAlignment
	newIssue.UserObjective.ExpectedEndDate = newIssue.StrategyInitiative.ExpectedEndDate
	newIssue.UserObjective.CreatedBy = newIssue.StrategyInitiative.CreatedBy
	newIssue.UserObjective.ModifiedBy = ctx.Request.User.ID

	return
}

func convertStrategyObjectivesToPlainTextOption(objs []models.StrategyObjective) (options []ebm.AttachmentActionElementPlainTextOption) {
	for _, o := range objs {
		options = append(options, ebm.AttachmentActionElementPlainTextOption{
			Value: o.ID,
			Label: ui.PlainText(o.Name),
		})
	}
	return
}
