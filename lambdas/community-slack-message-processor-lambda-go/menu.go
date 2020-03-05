package lambda

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

func CreateCommunityMenu(callbackId, userID string, teamID models.TeamID, comms []models.AdaptiveCommunity) []ebm.Attachment {
	var menuOptGroups []ebm.MenuOptionGroup
	// Check if there are any available communities to subscribe
	availComms := availableCommunities(teamID)

	if len(availComms) > 0 || len(comms) > 0 || len(availableStrategyCommunities(teamID, userID)) > 0 {
		var opts []ebm.MenuOption
		// one of subscribe or unsubscribe is true
		if len(availComms) > 0 || len(availableStrategyCommunities(teamID, userID)) > 0 {
			// subscribe option
			opts = append(opts, option(CommunitySubscribeAction, CommunitySubscribeText))
		}
		if len(comms) > 0 {
			// unsubscribe option
			opts = append(opts, option(CommunityUnsubscribeAction, CommunityUnsubscribeText))
		}
		menuOptGroups = append(menuOptGroups, optionGroup(CommunitiesMenuTitle, opts...))
	}

	for _, each := range comms {
		id := community.AdaptiveCommunity(each.ID)
		switch id {
		case community.Admin:
			menuOptGroups = append(menuOptGroups, simulationMenu)
		case community.Coaching:
		case community.Competency:
		case community.User:
		case community.Strategy:
		default:
		}
	}

	attachAction1, _ := eb.NewAttachmentActionBuilder().
		Name(CommunityMenuActionName).
		Text(string(MainMenuEmbeddedPrompt)).
		ActionType(ebm.AttachmentActionTypeSelect).
		OptionGroups(menuOptGroups).
		Build()

	attach, _ := eb.NewAttachmentBuilder().
		Title(string(MenuPrompt)).
		Fallback(string(MenuPromptFallback)).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		CallbackId(callbackId).
		Actions([]ebm.AttachmentAction{*attachAction1}).
		Build()
	return []ebm.Attachment{*attach}
}

// func getCommunitiesList() []string {

// }

// lifts a function string -> MenuOption to be a function from slice to slice
func liftStringToOption(f func(string) ebm.MenuOption) func([]string) []ebm.MenuOption {
	return func(opts []string) (menuOptions []ebm.MenuOption) {
		for _, v := range opts {
			menuOptions = append(menuOptions, simpleOption(ui.PlainText(v)))
		}
		return
	}
}

// lifts a function string -> MenuOption to be a function from slice to slice
func liftKvPairToOption(f func(models.KvPair) ebm.MenuOption) func([]models.KvPair) []ebm.MenuOption {
	return func(opts []models.KvPair) (menuOptions []ebm.MenuOption) {
		for _, each := range opts {
			menuOptions = append(menuOptions, f(each))
		}
		return
	}
}

func kvPairToMenuOption(kvPair models.KvPair) ebm.MenuOption {
	return option(kvPair.Value, ui.PlainText(kvPair.Key))
}

var (
	simulationMenu = optionGroup(
		SimulateMenuTitle,
		option(SimulateCurrentQuarterAction, SimulateCurrentQuarterText),
		option(SimulateNextQuarterAction, SimulateNextQuarterText),
	)
)
