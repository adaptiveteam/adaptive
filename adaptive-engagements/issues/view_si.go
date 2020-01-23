package issues

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// GetMainFields shows name and description
func (ViewInitiative) GetMainFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(NameLabel, getName, newAndOldIssues),
		attachmentFieldNewOld(DescriptionLabel, getDescription, newAndOldIssues),
		// {Title: string("Type"), Value: "Initiative"},
	}
	return
}

// GetDetailsFields shows all information except name/description
func (v ViewInitiative) GetDetailsFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	// For ViewMore action, we only need the latest comment
	fields = []ebm.AttachmentField{
		attachmentFieldNewOld(TimelineLabel, renderObjectiveViewDate, newAndOldIssues),
		attachmentFieldNewOld(InitiativeCommunityLabel, getInitiativeCommunity, newAndOldIssues),
		attachmentFieldNewOld(RelatedObjectiveLabel, getRelatedObjective, newAndOldIssues),
		attachmentFieldNewOld(DefinitionOfVictoryLabel, getDefinitionOfVictory, newAndOldIssues),
		attachmentFieldNewOld(BudgetLabel, getBudget, newAndOldIssues),
		attachmentFieldNewOld(AccountabilityPartnerLabel, getAccountabilityPartner, newAndOldIssues),
		attachmentFieldNewOld(StatusLabel, getStatus, newAndOldIssues),
		attachmentFieldNewOld(StrategyAssociationFieldLabel, v.GetAlignment, newAndOldIssues),
		attachmentField(LastReportedProgressLabel, getLatestComments(newAndOldIssues.NewIssue.PrefetchedData.Progress)),
	}
	return
}

// GetProgressFields shows only progress summary
func (ViewInitiative) GetProgressFields(newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField) {
	fields = []ebm.AttachmentField{
		userObjectiveProgressField(newAndOldIssues.NewIssue.PrefetchedData.Progress),
	}
	return
}

func (ViewInitiative) GetAlignment(issue Issue) (alignment ui.PlainText) {
	alignment = ui.PlainText(ui.Sprintf("%s - %s",
		renderStrategyAssociations("Initiative Communities", issue.AlignedInitiativeCommunity.Name),
		renderStrategyAssociations("Capability Objectives", issue.AlignedCapabilityObjective.Name),
	))
	return
}

func (ViewInitiative) GetTextView(issue Issue) ui.RichText {
	return ui.Sprintf("*%s*: %s \n *%s*: %s",
		NameLabel, issue.UserObjective.Name,
		DescriptionLabel, issue.UserObjective.Description)
}
