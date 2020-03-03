package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

func EditCapabilityCommunitySurveyElems(teamID models.TeamID, cc *strategy.CapabilityCommunity) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	if cc == nil {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    strategy.CapabilityCommunityNameLabel,
				Name:     strategy.CapabilityCommunityName,
				ElemType: string(ebm.ElemTypeTextBox),
			},
			{
				Label:    strategy.CapabilityCommunityDescriptionLabel,
				Name:     strategy.CapabilityCommunityDescription,
				ElemType: string(ebm.ElemTypeTextArea),
			},
			{
				Label:    strategy.CapabilityCommunityCoordinatorLabel,
				Name:     strategy.CapabilityCommunityCoordinator,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(allUsersInAnyStrategyCommunities(teamID)),
			},
		}
	} else {
		op = []ebm.AttachmentActionTextElement{
			{
				Label:    strategy.CapabilityCommunityNameLabel,
				Name:     strategy.CapabilityCommunityName,
				ElemType: string(ebm.ElemTypeTextBox),
				Value:    cc.Name,
			},
			{
				Label:    strategy.CapabilityCommunityDescriptionLabel,
				Name:     strategy.CapabilityCommunityDescription,
				ElemType: string(ebm.ElemTypeTextArea),
				Value:    cc.Description,
			},
			{
				Label:    strategy.CapabilityCommunityCoordinatorLabel,
				Name:     strategy.CapabilityCommunityCoordinator,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(allUsersInAnyStrategyCommunities(teamID)),
				Value:    cc.Advocate,
			},
		}
	}
	return op
}
