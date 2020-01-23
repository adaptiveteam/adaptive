package issues

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// GetMainFields shows name and description
func (ViewSObjective) GetMainFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(NameLabel, getName, newAndOldIssues),
		attachmentFieldNewOld(DescriptionLabel, getDescription, newAndOldIssues),
		// {Title: string("Type"), Value: "Initiative"},
	}
	return
}

// GetDetailsFields shows all information except name/description
func (v ViewSObjective) GetDetailsFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	// For ViewMore action, we only need the latest comment
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(SObjectiveTypeLabel, getObjectiveType, newAndOldIssues),
		attachmentFieldNewOld(TimelineLabel, renderObjectiveViewDate, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveEndDateLabel, getObjectiveExpectedEndDate, newAndOldIssues),
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveAdvocateLabel, getObjectiveAdvocate, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveMeasuresLabel, getAsMeasuredBy, newAndOldIssues),
		attachmentFieldNewOld(SObjectiveTargetsLabel, getTargets, newAndOldIssues),
		attachmentFieldNewOld(StatusLabel, getStatus, newAndOldIssues),
		attachmentFieldNewOld(StrategyAssociationFieldLabel, v.GetAlignment, newAndOldIssues),
		attachmentField(LastReportedProgressLabel, getLatestComments(newAndOldIssues.NewIssue.PrefetchedData.Progress)),
	}
	return
}

// GetProgressFields shows only progress summary
func (ViewSObjective) GetProgressFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		userObjectiveProgressField(newAndOldIssues.NewIssue.PrefetchedData.Progress),
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

func getObjectiveExpectedEndDate(issue Issue) ui.PlainText {
	so := issue.StrategyObjective
	newExpectedEndDate := formatDate(so.ExpectedEndDate, core.ISODateLayout, core.USDateLayout)
	return ui.PlainText(newExpectedEndDate)
}

func (ViewSObjective) GetAlignment(issue Issue) (alignment ui.PlainText) {
	alignment = ui.PlainText(ui.Sprintf("`%s Objective` : `%s`\n", issue.StrategyObjective.ObjectiveType, issue.StrategyObjective.Name))
	return
}

func (ViewSObjective) GetTextView(issue Issue) ui.RichText {
	return ui.Sprintf("*%s*: %s \n *%s*: %s",
		NameLabel, issue.UserObjective.Name,
		DescriptionLabel, issue.UserObjective.Description)
}
