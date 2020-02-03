package lambda

import (
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	VisionDescription     = "vision_description"
	VisionMissionAdvocate = "mission_vision_advocate"

	VisionLabel                ui.PlainText = "Vision"
	MissionVisionAdvocateLabel              = "Vision Advocate"
)

func EditVisionMissionSurveyElems(vm *models.VisionMission, advocates []models.KvPair) []ebm.AttachmentActionTextElement {
	var op []ebm.AttachmentActionTextElement
	// Survey box should consist of a menu option to select rating and a text area for a user to enter the feedback
	if vm == nil {
		op = []ebm.AttachmentActionTextElement{
			ebm.NewTextArea(VisionDescription, VisionLabel, ebm.EmptyPlaceholder, ""),
			{
				Label:    MissionVisionAdvocateLabel,
				Name:     VisionMissionAdvocate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
			},
		}
	} else {
		op = []ebm.AttachmentActionTextElement{
			ebm.NewTextArea(VisionDescription, VisionLabel, ebm.EmptyPlaceholder, ui.PlainText(vm.Vision)),
			{
				Label:    MissionVisionAdvocateLabel,
				Name:     VisionMissionAdvocate,
				ElemType: models.MenuSelectType,
				Options:  utils.AttachActionElementOptions(advocates),
				Value:    vm.Advocate,
			},
		}
	}
	return op
}
