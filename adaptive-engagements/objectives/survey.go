package objectives

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	ObjectiveNamePlaceholder        ui.PlainText = ebm.EmptyPlaceholder
	ObjectiveDescriptionPlaceholder ui.PlainText = ebm.EmptyPlaceholder
)

func AlignmentFromAlignedStrategyType(tpe models.AlignedStrategyType, id string) (alignment string) {
	switch tpe {
	case models.ObjectiveStrategyObjectiveAlignment:
		alignment = fmt.Sprintf("%s:%s", community.Capability, id)
	case models.ObjectiveStrategyInitiativeAlignment:
		alignment = fmt.Sprintf("%s:%s", community.Initiative, id)
	case models.ObjectiveCompetencyAlignment:
		alignment = fmt.Sprintf("%s:%s", models.ObjectiveCompetencyAlignment, id)
	case models.ObjectiveNoStrategyAlignment:
		alignment = "None"
	}
	return
}

func EditObjectiveSurveyElems2(obj *models.UserObjective, coaches, dates []models.KvPair,
	initiativesAndObjectives []ebm.AttachmentActionElementOptionGroup) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	alignment := ""
	if obj != nil {
		alignment = AlignmentFromAlignedStrategyType(obj.StrategyAlignmentEntityType, obj.StrategyAlignmentEntityID)
	}

	if obj == nil {
		obj = &models.UserObjective{
			Name:                  "",
			Description:           "",
			AccountabilityPartner: "",
			ExpectedEndDate:       "",
		}
	}

	op = []ebm.AttachmentActionTextElement{
		ebm.NewTextBox(ObjectiveName, "Name", ObjectiveNamePlaceholder, ui.PlainText(obj.Name)),
		ebm.NewTextArea(ObjectiveDescription, "Description", ObjectiveDescriptionPlaceholder, ui.PlainText(obj.Description)),
		{
			Label:    "Coach",
			Name:     ObjectiveAccountabilityPartner,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(coaches),
			Value:    obj.AccountabilityPartner,
		},
		{
			Label:    "Expected end date",
			Name:     ObjectiveEndDate,
			ElemType: models.MenuSelectType,
			Options:  utils.AttachActionElementOptions(dates),
			Value:    obj.ExpectedEndDate,
		},
	}
	if len(initiativesAndObjectives) > 0 {
		op = append(op,
			ebm.AttachmentActionTextElement{
				Label:        "Strategy Alignment",
				Name:         ObjectiveStrategyAlignment,
				ElemType:     models.MenuSelectType,
				OptionGroups: initiativesAndObjectives,
				// Options:  utils.AttachActionElementOptions(initiativesAndObjectives),
				Value: alignment,
			},
		)
	}
	return op
}
