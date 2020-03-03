package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	InitiativeNameLabel                ui.PlainText = "Initiative Name"
	InitiativeDescriptionLabel         ui.PlainText = "Initiative Description"
	InitiativeVictoryLabel             ui.PlainText = "Definition of Victory"
	InitiativeAdvocateLabel                         = "Advocate"
	InitiativeBudgetLabel                           = "Budget($) in the following format: 1234.56"
	InitiateEndDateLabel                            = "Time to work on this"
	InitiativeCapabilityObjectiveLabel              = "Related Capability Objective"

	InitiativeName                    = "initiative_name"
	InitiativeDescriptionName         = "initiative_description"
	InitiativeVictoryName             = "definition_of_victory"
	InitiativeAdvocateName            = "advocate"
	InitiativeBudgetName              = "initiative_budget_name"
	InitiateEndDateName               = "time_to_work_on_this"
	InitiativeCapabilityObjectiveName = "initiative_capability_objective"
)

func EditInitiativeSurveyElems(si *models.StrategyInitiative, advocates, dates,
	objectives []models.KvPair) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	if si == nil {
		op = []ebm.AttachmentActionTextElement{
			ebm.NewTextBox(InitiativeName, InitiativeNameLabel, ebm.EmptyPlaceholder, ""),
			ebm.NewTextArea(InitiativeDescriptionName, InitiativeDescriptionLabel, ebm.EmptyPlaceholder, ""),
			ebm.NewTextArea(InitiativeVictoryName, InitiativeVictoryLabel, ebm.EmptyPlaceholder, ""),
			{
				Label:    InitiativeAdvocateLabel,
				Name:     InitiativeAdvocateName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
			},
			{
				Label:       InitiativeBudgetLabel,
				Name:        InitiativeBudgetName,
				ElemType:    "text",
				ElemSubtype: ebm.AttachmentActionTextElementNumberType,
			},
			{
				Label:    InitiateEndDateLabel,
				Name:     InitiateEndDateName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(dates),
			},
			{
				Label:    InitiativeCapabilityObjectiveLabel,
				Name:     InitiativeCapabilityObjectiveName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(objectives),
			},
		}
	} else {
		op = []ebm.AttachmentActionTextElement{
			ebm.NewTextBox(InitiativeName, InitiativeNameLabel, ebm.EmptyPlaceholder, ui.PlainText(si.Name)),
			ebm.NewTextArea(InitiativeDescriptionName, InitiativeDescriptionLabel, ebm.EmptyPlaceholder, ui.PlainText(si.Description)),
			ebm.NewTextArea(InitiativeVictoryName, InitiativeVictoryLabel, ebm.EmptyPlaceholder, ui.PlainText(si.DefinitionOfVictory)),
			{
				Label:    InitiativeAdvocateLabel,
				Name:     InitiativeAdvocateName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
				Value:    si.Advocate,
			},
			{
				Label:       InitiativeBudgetLabel,
				Name:        InitiativeBudgetName,
				ElemType:    "text",
				ElemSubtype: ebm.AttachmentActionTextElementNumberType,
				Value:       si.Budget,
			},
			{
				Label:    InitiateEndDateLabel,
				Name:     InitiateEndDateName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(dates),
				Value:    si.ExpectedEndDate,
			},
			{
				Label:    InitiativeCapabilityObjectiveLabel,
				Name:     InitiativeCapabilityObjectiveName,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(objectives),
				Value:    si.CapabilityObjective,
			},
		}
	}
	return op
}

func EditInitiativeCommunitySurveyElems(teamID models.TeamID, cc *strategy.StrategyInitiativeCommunity,
	capabilityComms []models.KvPair) (op []ebm.AttachmentActionTextElement) {
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	if cc == nil {
		cc = &strategy.StrategyInitiativeCommunity{}
	}
	op = []ebm.AttachmentActionTextElement{
		{
			Label:    strategy.InitiativeCommunityNameLabel,
			Name:     strategy.InitiativeCommunityName,
			ElemType: "text",
			Value:    cc.Name,
		},
		{
			Label:    strategy.InitiativeCommunityDescriptionLabel,
			Name:     strategy.InitiativeCommunityDescription,
			ElemType: string(ebm.ElemTypeTextArea),
			Value:    cc.Description,
		},
		{
			Label:    strategy.InitiativeCommunityCoordinatorLabel,
			Name:     strategy.InitiativeCommunityCoordinator,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(allUsersInAnyStrategyCommunities(teamID)),
			Value:    cc.Advocate,
		},
		{
			Label:    strategy.InitiativeCommunityCapabilityCommunityLabel,
			Name:     strategy.InitiativeCommunityCapabilityCommunity,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(capabilityComms),
			Value:    cc.CapabilityCommunityID,
		},
	}
	return
}
