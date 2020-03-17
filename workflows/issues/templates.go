package issues

import (
	"fmt"
	"sort"

	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"

	community "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	models "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// formatCompetenciesAsOptionGroup converts the list of values to option group if non empty
func formatCompetenciesAsOptionGroup(competencies []adaptiveValue.AdaptiveValue) (res []ebm.AttachmentActionElementOptionGroup) {
	if len(competencies) > 0 {
		grp := ebm.AttachmentActionElementOptionGroup{Label: "Competencies"}
		options := grp.Options
		for _, each := range competencies {
			if each.Name == "" {
				fmt.Printf("Problematic competency: %v\n", each)
				each.Name = "Empty name (value_name?): " + each.ID
			}

			options = append(options,
				ebm.AttachmentActionElementOption{
					Label: each.Name,
					Value: fmt.Sprintf("%s:%s", community.Competency, each.ID),
				})
		}
		grp.Options = options
		res = append(res, grp)
	}
	return
}

// initiativesGroup formats one option group with initiatives
func formatInitiativesAsGroup(initiativesForUser []models.StrategyInitiative) (res []ebm.AttachmentActionElementOptionGroup) {
	if len(initiativesForUser) != 0 {
		grp := ebm.AttachmentActionElementOptionGroup{}
		options := grp.Options
		for _, each := range initiativesForUser {
			options = append(options,
				ebm.AttachmentActionElementOption{
					Label: core.ClipString(each.Name, 30, "..."), // get first
					Value: fmt.Sprintf("%s:%s", community.Initiative, each.ID),
				})
		}
		sort.Sort(MenuOptionLabelSorter(options))
		grp.Options = options
		grp.Label = ui.PlainText("Initiatives")
		res = append(res, grp)
	}
	return
}

// MenuOptionLabelSorter is a type that is only used to sort menu options by label
type MenuOptionLabelSorter []ebm.AttachmentActionElementOption

func (a MenuOptionLabelSorter) Len() int           { return len(a) }
func (a MenuOptionLabelSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuOptionLabelSorter) Less(i, j int) bool { return a[i].Label < a[j].Label }

func mapInitiativesToString(items []models.StrategyInitiative, f func(models.StrategyInitiative) string) (res []string) {
	for _, i := range items {
		res = append(res, f(i))
	}
	return
}
func mapObjectivesToString(items []models.StrategyObjective, f func(models.StrategyObjective) string) (res []string) {
	for _, i := range items {
		res = append(res, f(i))
	}
	return
}
func mapObjectivesToAttachmentActionElementOption(items []models.StrategyObjective, f func(models.StrategyObjective) ebm.AttachmentActionElementOption) (res []ebm.AttachmentActionElementOption) {
	for _, i := range items {
		res = append(res, f(i))
	}
	return
}
func mapCommunityUsersToPlainTextOption(items []models.AdaptiveCommunityUser3, f func(models.AdaptiveCommunityUser3) ebm.AttachmentActionElementPlainTextOption) (res []ebm.AttachmentActionElementPlainTextOption) {
	for _, i := range items {
		res = append(res, f(i))
	}
	return
}

// objectives formats one option group with objectives
func formatObjectivesGroup(capabilityObjectives []models.StrategyObjective) (res []ebm.AttachmentActionElementOptionGroup) {
	// Add each Capability Objective from retrieved Capability Objectives in options
	options := mapObjectivesToAttachmentActionElementOption(capabilityObjectives,
		func(each models.StrategyObjective) ebm.AttachmentActionElementOption {
			return ebm.AttachmentActionElementOption{
				Label: core.ClipString(each.Name, 30, "…"),
				Value: fmt.Sprintf("%s:%s", community.Capability, each.ID),
			}
		})
	sort.Sort(MenuOptionLabelSorter(options))

	// adding options to group only when they exist
	// reference error: Element 2 field `options` must have at least one option
	if len(options) > 0 {
		grp := ebm.AttachmentActionElementOptionGroup{
			Label:   "Objectives",
			Options: options,
		}
		res = append(res, grp)
	}
	return
}

type GlobalDialogContext = string
type GlobalDialogContextByType = map[IssueType]GlobalDialogContext

var contexts = map[DialogSituationIDWithoutIssueType]GlobalDialogContextByType{
	DescriptionContext: {
		IDO:        "dialog/ido/language-coaching/description",
		SObjective: "",
		Initiative: "",
	},
	CloseoutAgreementContext: {
		IDO:        "dialog/ido/language-coaching/close-out-agreement",
		SObjective: "dialog/strategy/language-coaching/objective/close-out-agreement",
		Initiative: "dialog/strategy/language-coaching/initiative/close-out-agreement",
	},
	CloseoutDisagreementContext: {
		IDO:        "dialog/ido/language-coaching/close-out-disagreement",
		SObjective: "dialog/strategy/language-coaching/objective/close-out-disagreement",
		Initiative: "dialog/strategy/language-coaching/initiative/close-out-disagreement",
	},
	UpdateContext: {
		IDO:        "dialog/ido/language-coaching/update",
		SObjective: "dialog/strategy/language-coaching/objective/update",
		Initiative: "dialog/strategy/language-coaching/initiative/update",	
	},
	UpdateResponseContext: {
		IDO:        "dialog/ido/language-coaching/update-response",
		SObjective: "dialog/strategy/language-coaching/objective/update-response",
		Initiative: "dialog/strategy/language-coaching/initiative/update-response",
	},
	CoachingRequestRejectionContext: {
		IDO:        "dialog/ido/language-coaching/coaching-request-rejection",
		SObjective: "",
		Initiative: "",
	},
	ProgressUpdateContext: { // TODO: provide progress update contexts
		IDO:        "dialog/ido/language-coaching/update",
		SObjective: "dialog/strategy/language-coaching/objective/update",
		Initiative: "dialog/strategy/language-coaching/initiative/update",	
	},
}
// GetDialogContext returns dialog context. 
// In case of errors logs, and returns empty context.
func (w workflowImpl) GetDialogContext(dialogSituationID DialogSituationIDWithoutIssueType, itype IssueType) (context string) {
	contextsForSituation, ok := contexts[dialogSituationID]
	if ok {
		context, ok = contextsForSituation[itype]
		if !ok { w.AdaptiveLogger.
			WithField("dialogSituationID", string(dialogSituationID)).
			WithField("itype", string(itype)).
			Info("Missing dialog context")
		}
	} else { w.AdaptiveLogger.
		WithField("dialogSituationID", string(dialogSituationID)).
		Info("Missing dialog situation")
	}
	return
}
const (
	BlueDiamondEmoji = ":small_blue_diamond:"
)

func ObjectiveCreatedUpdatedStatusTemplate(issueType IssueType, updated bool, userID string) (text ui.RichText) {
	var verb = ""
	if updated {
		verb = "updated"
	} else {
		verb = "created"
	}
	text = ui.Sprintf("Below %s has been %s by <@%s>", issueType.Template(), verb, userID)
	return
}

func ObjectiveProgressCreatedUpdatedStatusTemplate(issueType IssueType, progressUpdated bool, userID string) (text ui.RichText) {
	var verb = ""
	if progressUpdated {
		verb = "updated"
	} else {
		verb = "added"
	}
	text = ui.Sprintf("Below %s's progress comment has been %s by <@%s>", issueType.Template(), verb, userID)
	return
}

