package issues

import (
	"log"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
	"time"
	"fmt"
	
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	objectives "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// CreateDialog loads data and creates StrategyObjective dialog
func (SObjectiveImpl) CreateDialog(w workflowImpl, ctx wf.EventHandlingContext, issue Issue) (survey ebm.AttachmentActionSurvey, err error) {
	defer w.recoverToErrorVar("SObjectiveImpl.CreateDialog", &err) // because there are panics downstream
	capCommID := ctx.Data[capCommIDKey]
	commID := community.AdaptiveCommunity(fmt.Sprintf("%s:%s", community.Capability, capCommID))
	var types, advocates, dates []ebm.AttachmentActionElementPlainTextOption

	var commMembers, strMembers []models.KvPair
	strMembers, err = SelectKvPairsFromCommunityJoinUsers(community.Strategy)(w.DynamoDBConnection)
	if err == nil {
		commMembers, err = SelectKvPairsFromCommunityJoinUsers(commID)(w.DynamoDBConnection)
		if err == nil {
			allMembers := append(strMembers, commMembers...)
			allDates := objectives.StrategyObjectiveDatesWithIndefiniteOption("SObjectiveImpl CreateDialog", issue.UserObjective.ExpectedEndDate)
			advocates = convertKvPairToPlainTextOption(allMembers)
			dates = convertKvPairToPlainTextOption(allDates)
			types = convertKvPairToPlainTextOption(ObjectiveTypes())

			survey = SOObjectiveSurvey(issue.StrategyObjective, types, advocates, dates)
		}
	}
	return
}

// SOObjectiveSurvey shows a form to create or modify an objective
func SOObjectiveSurvey(item models.StrategyObjective,
	types, advocates, dates []ebm.AttachmentActionElementPlainTextOption) ebm.AttachmentActionSurvey {
	return ebm.AttachmentActionSurvey{
		Title: "Objective",
		Elements: []ebm.AttachmentActionTextElement{
			ebm.NewSimpleOptionsSelect(SObjectiveType, SObjectiveTypeLabel, ebm.EmptyPlaceholder, string(item.ObjectiveType), types...),
			ebm.NewTextBox(SObjectiveName, SObjectiveNameLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Name)),
			ebm.NewTextArea(SObjectiveDescription, SObjectiveDescriptionLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Description)),
			ebm.NewTextArea(SObjectiveMeasures, SObjectiveMeasuresLabel, ebm.EmptyPlaceholder, ui.PlainText(item.AsMeasuredBy)),
			ebm.NewTextArea(SObjectiveTargets, SObjectiveTargetsLabel, ebm.EmptyPlaceholder, ui.PlainText(item.Targets)),
			ebm.NewSimpleOptionsSelect(SObjectiveAdvocate, SObjectiveAdvocateLabel, ebm.EmptyPlaceholder, item.Advocate, advocates...),
			ebm.NewSimpleOptionsSelect(SObjectiveEndDate, SObjectiveEndDateLabel, ebm.EmptyPlaceholder, item.ExpectedEndDate, dates...),
		},
	}
}

const (
	SObjectiveName        = "s_objective_name"
	SObjectiveDescription = "s_objective_description"
	SObjectiveMeasures    = "s_objective_measures"
	SObjectiveTargets     = "s_objective_targets"
	SObjectiveType        = "s_objective_type"
	SObjectiveAdvocate    = "s_objective_advocate"
	SObjectiveEndDate     = "s_objective_end_Date"

	// labels
	SObjectiveNameLabel        ui.PlainText = "Name"
	SObjectiveDescriptionLabel ui.PlainText = "Description"
	SObjectiveMeasuresLabel    ui.PlainText = "Measures"
	SObjectiveTargetsLabel     ui.PlainText = "Targets"
	SObjectiveTypeLabel                     = "Type"
	SObjectiveAdvocateLabel                 = "Advocate"
	SObjectiveEndDateLabel                  = "Time to work on this"
)

const ObjectiveTypeDefaultValue = "No Type"

// ObjectiveTypes is the collection of objective types.
// should be saved to DB.
func ObjectiveTypes() []models.KvPair {
	return []models.KvPair{
		// {Key: "Customer Strategy Objective", Value: string(strategy.CustomerStrategyObjective)},
		// {Key: "Financial Strategy Objective", Value: string(strategy.FinancialStrategyObjective)},
		// {Key: "Capability Strategy Objective", Value: string(strategy.CapabilityStrategyObjective)},
		{Key: "No Type", Value: "No Type"},
		{Key: "Financial Performance", Value: "Financial Performance"},
		{Key: "Effective Resource Use", Value: "Effective Resource Use"},
		{Key: "Customer Value", Value: "Customer Value"},
		{Key: "Customer Satisfaction", Value: "Customer Satisfaction"},
		{Key: "Customer Retention", Value: "Customer Retention"},
		{Key: "Efficiency", Value: "Efficiency"},
		{Key: "Quality", Value: "Quality"},
		{Key: "People Development", Value: "People Development"},
		{Key: "Infrastructure", Value: "Infrastructure"},
		{Key: "Technology", Value: "Technology"},
		{Key: "Culture", Value: "Culture"},
	}
}

func (SObjectiveImpl) ExtractFromContext(ctx wf.EventHandlingContext, _ string, updated bool, oldIssue Issue) (newIssue Issue) {
	form := ctx.Request.DialogSubmissionCallback.Submission
	newIssue.StrategyObjective.Name = form[SObjectiveName]
	newIssue.StrategyObjective.ObjectiveType = strategyObjective.StrategyObjectiveType(form[SObjectiveType])
	newIssue.StrategyObjective.Description = form[SObjectiveDescription]
	newIssue.StrategyObjective.AsMeasuredBy = form[SObjectiveMeasures]
	newIssue.StrategyObjective.Targets = form[SObjectiveTargets]
	newIssue.StrategyObjective.Advocate = form[SObjectiveAdvocate]
	newIssue.StrategyObjective.ExpectedEndDate = form[SObjectiveEndDate]
	newIssue.StrategyObjective.PlatformID = ctx.PlatformID
	now := core.TimestampLayout.Format(time.Now())
	newIssue.UserObjective.ModifiedAt = now
	if updated {
		newIssue.StrategyObjective.ID = oldIssue.StrategyObjective.ID
		newIssue.StrategyObjective.CreatedBy = oldIssue.StrategyObjective.CreatedBy
		newIssue.StrategyObjective.CreatedAt = oldIssue.StrategyObjective.CreatedAt
		newIssue.StrategyObjective.CapabilityCommunityIDs = oldIssue.StrategyObjective.CapabilityCommunityIDs

	} else {
		newIssue.StrategyObjective.ID = core.Uuid()
		newIssue.StrategyObjective.CreatedBy = ctx.Request.User.ID
		newIssue.StrategyObjective.CreatedAt = now
		// time.Now().Format(string(TimestampFormat))
		capCommID := ctx.Data[capCommIDKey]
		if capCommID == "" {
			log.Printf("!Invalid state - there is no capCommID in the context")
		}
		newIssue.StrategyObjective.CapabilityCommunityIDs = []string{capCommID}
	}
	newIssue.UserObjective.ID = newIssue.StrategyObjective.ID
	newIssue.UserObjective.Name = newIssue.StrategyObjective.Name
	newIssue.UserObjective.Description = newIssue.StrategyObjective.Description
	newIssue.UserObjective.UserID = newIssue.StrategyObjective.Advocate
	newIssue.UserObjective.AccountabilityPartner = newIssue.StrategyObjective.CreatedBy
	newIssue.UserObjective.Name = newIssue.StrategyObjective.Name
	newIssue.UserObjective.PlatformID = ctx.PlatformID
	newIssue.UserObjective.CreatedAt = newIssue.StrategyObjective.CreatedAt
	newIssue.UserObjective.CreatedDate = core.NormalizeDate(newIssue.StrategyObjective.CreatedAt)
	newIssue.UserObjective.ObjectiveType = userObjective.StrategyDevelopmentObjective
	newIssue.UserObjective.StrategyAlignmentEntityType = userObjective.ObjectiveStrategyObjectiveAlignment
	newIssue.UserObjective.ExpectedEndDate = newIssue.StrategyObjective.ExpectedEndDate
	return
}

func (s SObjectiveImpl) View(w workflowImpl, isShowingDetails, isShowingProgress bool,
	newAndOldIssues NewAndOldIssues,
	) (fields []ebm.AttachmentField) {
	fields = s.ObjectiveToFields(w, newAndOldIssues)
	if isShowingDetails {
		fields = append(fields, s.ObjectiveToFieldDetails(w, newAndOldIssues)...)
	}
	if isShowingProgress {
		fields = append(fields, userObjectiveProgressField(newAndOldIssues.NewIssue.PrefetchedData.Progress))
	}
	return
}

func (s SObjectiveImpl)ObjectiveToFields(w workflowImpl, newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(SObjectiveNameLabel, getObjectiveName, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveDescriptionLabel, getObjectiveDescription, newAndOldIssues),
	}
	return
}

func getAsMeasuredBy(issue Issue) ui.PlainText {
	return ui.PlainText(issue.AsMeasuredBy)
}

func getTargets(issue Issue) ui.PlainText {
	return ui.PlainText(issue.Targets)
}

func getObjectiveAdvocate(issue Issue) ui.PlainText {
	return ui.PlainText(common.TaggedUser(issue.StrategyObjective.Advocate))
}

func getObjectiveName(issue Issue) ui.PlainText {
	return ui.PlainText(issue.StrategyObjective.Name)
}

func getObjectiveDescription(issue Issue) ui.PlainText {
	return ui.PlainText(issue.StrategyObjective.Description)
}

func getObjectiveType(issue Issue) ui.PlainText {
	return ui.PlainText(issue.StrategyObjective.ObjectiveType)
}

func (w workflowImpl)getObjectiveExpectedEndDate(issue Issue) ui.PlainText {
	so := issue.StrategyObjective
	newExpectedEndDate := formatDate(w, so.ExpectedEndDate, core.ISODateLayout, core.USDateLayout)
	return ui.PlainText(newExpectedEndDate)
}


func (s SObjectiveImpl)ObjectiveToFieldDetails(w workflowImpl, newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	// For ViewMore action, we only need the latest comment
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(SObjectiveTypeLabel, getObjectiveType, newAndOldIssues),
		attachmentFieldNewOld(TimelineLabel, renderObjectiveViewDate, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveEndDateLabel, w.getObjectiveExpectedEndDate, newAndOldIssues),
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveAdvocateLabel, getObjectiveAdvocate, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveMeasuresLabel, getAsMeasuredBy, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveTargetsLabel, getTargets, newAndOldIssues),
		attachmentFieldNewOld(StatusLabel, getStatus, newAndOldIssues),
		attachmentField(LastReportedProgressLabel, getLatestComments(newAndOldIssues.NewIssue.PrefetchedData.Progress)),
	}
	return
}

const (
	AccountabilityPartnerLabel ui.PlainText = "Accountability Partner"
	StatusLabel                ui.PlainText = "Status"
	LastReportedProgressLabel  ui.PlainText = "Last reported progress"
)

const (
	StatusCancelled                                ui.PlainText = "Cancelled"
	StatusPending                                  ui.PlainText = "Pending"
	StatusCompletedAndPartnerVerifiedCompletion    ui.PlainText = "Completed by you and closeout approved by your partner"
	StatusCompletedAndNotPartnerVerifiedCompletion ui.PlainText = "Completed by you and pending closeout approval from your partner"
)
