package issues

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

// GetMainFields shows name and description
func (ViewIDO) GetMainFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(NameLabel, getName, newAndOldIssues),
		attachmentFieldNewOld(DescriptionLabel, getDescription, newAndOldIssues),
		// {Title: string("Type"), Value: "Initiative"},
	}
	return
}
// GetDetailsFields shows all information except name/description
func (v ViewIDO) GetDetailsFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(StrategyAssociationFieldLabel, v.GetAlignment, newAndOldIssues),
		attachmentFieldNewOld(TimelineLabel, renderObjectiveViewDate, newAndOldIssues),
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner, newAndOldIssues),
		// {Title: string("Type"), Value: "Individual"},
		
		attachmentFieldNewOld(StatusLabel, getStatus, newAndOldIssues),
		attachmentField(LastReportedProgressLabel, getLatestComments(newAndOldIssues.NewIssue.PrefetchedData.Progress)),
	}

	return
}
// GetProgressFields shows only progress summary
func (ViewIDO) GetProgressFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		userObjectiveProgressField(newAndOldIssues.NewIssue.PrefetchedData.Progress),
	}
	return
}

func (ViewIDO) GetAlignment(issue Issue) (alignment ui.PlainText) {
	switch issue.StrategyAlignmentEntityType {
	case userObjective.ObjectiveStrategyObjectiveAlignment:
		alignment = renderStrategyAssociations("Capability Objective", issue.AlignedCapabilityObjective.Name)
	case userObjective.ObjectiveStrategyInitiativeAlignment:
		alignment = renderStrategyAssociations("Initiative", issue.AlignedCapabilityInitiative.Name)
	case userObjective.ObjectiveCompetencyAlignment:
		alignment = ui.PlainText(ui.Sprintf("Competency: `%s`", issue.AlignedCompetency.Name))
	}
	return
}

func (ViewIDO) GetTextView(issue Issue) ui.RichText {
	return ui.Sprintf("*%s*: %s \n *%s*: %s",
		NameLabel, issue.UserObjective.Name,
		DescriptionLabel, issue.UserObjective.Description)
}
