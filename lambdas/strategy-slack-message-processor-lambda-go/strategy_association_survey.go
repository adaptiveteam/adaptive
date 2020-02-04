package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	StrategyObjectiveAssociationCommunity            = "strategy_objective_association_community"
	StrategyObjectiveCommunityAssociationDescription = "strategy_objective_community_association_description"
)

func EditStrategyAssociation(so *models.StrategyObjective, targetList []models.KvPair,
	targetLabel, descriptionLabel ui.PlainText) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	if so == nil {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    string(targetLabel),
				Name:     StrategyObjectiveAssociationCommunity,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(targetList),
			},
			ebm.NewTextArea(StrategyObjectiveCommunityAssociationDescription, descriptionLabel, ebm.EmptyPlaceholder,
				""),
		}
	} else {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    string(targetLabel),
				Name:     StrategyObjectiveAssociationCommunity,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(targetList),
			},
			ebm.NewTextArea(StrategyObjectiveCommunityAssociationDescription, descriptionLabel, ebm.EmptyPlaceholder,
				ui.PlainText(so.Description)),
		}
	}
	return op
}
