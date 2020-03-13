package issues

import (
	// alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	// "encoding/json"
	// "fmt"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	// "github.com/adaptiveteam/adaptive/adaptive-engagements/objectives"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/strategy"
	// "github.com/thoas/go-funk"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	// utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	// core "github.com/adaptiveteam/adaptive/core-utils-go"
	// ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	// mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	// "github.com/adaptiveteam/adaptive/daos/strategyObjective"
	// "log"
)

func (w workflowImpl)OnCreateItem(isFromMainMenu bool, itype IssueType) wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		ctx.Data[issueTypeKey] = string(itype)
		if UserHasWriteAccessToIssuesImpl(w.DynamoDBConnection)(ctx.Request.User.ID, itype) {
			tc := getTypeClass(itype)
			if tc.IsCapabilityCommunityNeeded() {
				out, err = w.requireCapabilityCommunity(ctx, tc)
			} else if tc.IsInitiativeCommunityNeeded() {
				out, err = w.requireInitiativeCommunity(ctx, tc)
			} else {
				out, err = w.showDialog(ctx, tc.Empty(), DescriptionContext)
				out.NextState = FormShownState
			}
		} else {
			// send a message that user is not authorized to create objectives
			out = ctx.Reply("You are not part of the Adaptive Strategy Community or an Objective Community, " +
			"you will not be able to create Capability Objectives.")
	
		}
		out.KeepOriginal = !isFromMainMenu
		return
	}
}
func (w workflowImpl)requireCapabilityCommunity(ctx wf.EventHandlingContext, tc IssueTypeClass) (out wf.EventOutput, err error) { 
	// check if the user is in strategy community
	var adaptiveAssociatedCapComms []strategy.CapabilityCommunity
	adaptiveAssociatedCapComms, err = SelectFromCapabilityCommunityJoinStrategyCommunityWhereChannelCreated()(w.DynamoDBConnection)
	if err != nil {return}
	switch len(adaptiveAssociatedCapComms) {
	case 0:
		out.NextState = DoneState
		out = ctx.Reply("There are no Adaptive associated Objective Communities. " +
				"If you have already created a Objective Community, " +
				"please ask the coordinator to create a *_private_* channel, " +
				"invite Adaptive and associate with the community.")
	case 1: // we already know the community. No need to ask.
		capCommID := adaptiveAssociatedCapComms[0].ID
		ctx.Data[capCommIDKey] = capCommID
		out, err = w.showDialog(ctx, tc.Empty(), DescriptionContext)
		out.NextState = FormShownState
	default:
		out.NextState = CommunitySelectingState
		opts := w.mapCapabilityCommunitiesToOptions(adaptiveAssociatedCapComms, ctx.TeamID)
		// Enable a user to create an objective if user is in strategy community and there are capability communities
		out.Interaction = wf.Buttons(
			"Select an objective community. You can assign the "+ ui.RichText(tc.IssueTypeName()) + " to other communities later but you need at least one for now.",
			wf.Selectors(wf.Selector{Event: CommunitySelectedEvent, Options: opts})...) // , wf.MenuOption("ignore", "Not now"))
	}
	return 
}

func (w workflowImpl)requireInitiativeCommunity(ctx wf.EventHandlingContext, tc IssueTypeClass) (out wf.EventOutput, err error) { 
	// check if the user is in strategy community
	var adaptiveAssociatedComms []strategy.StrategyInitiativeCommunity
	adaptiveAssociatedComms, err = SelectFromInitiativeCommunityJoinStrategyCommunityWhereChannelCreated(ctx.Request.User.ID)(w.DynamoDBConnection)
	if err != nil {return}
	switch len(adaptiveAssociatedComms) {
	case 0:
		out.NextState = DoneState
		out = ctx.Reply("There are no Adaptive associated Initiative Communities. " +
				"If you have already created an Initiative Community, " +
				"please ask the coordinator to create a *_private_* channel, " +
				"invite Adaptive and associate with the community.")
	case 1: // we already know the community. No need to ask.
		initCommID := adaptiveAssociatedComms[0].ID
		ctx.Data[initCommIDKey] = initCommID
		out, err = w.showDialog(ctx, tc.Empty(), DescriptionContext)
		out.NextState = FormShownState
	default:
		out.NextState = CommunitySelectingState
		opts := w.mapInitiativeCommunitiesToOptions(adaptiveAssociatedComms, ctx.TeamID)
		// Enable a user to create an objective if user is in strategy community and there are capability communities
		out.Interaction = wf.Buttons(
			"Select an initiative community. You can assign the "+ ui.RichText(tc.IssueTypeName()) + " to other communities later but you need at least one for now.",
			wf.Selectors(wf.Selector{Event: CommunitySelectedEvent, Options: opts})...) // , wf.MenuOption("ignore", "Not now"))
	}
	return 
}

func (w workflowImpl)mapCapabilityCommunitiesToOptions(comms []strategy.CapabilityCommunity, teamID models.TeamID) (opts []wf.SelectorOption) {
	for _, each := range comms {
		opts = append(opts, wf.SelectorOption{
			Label: ui.PlainText(each.Name), 
			Value: each.ID,
		})
	}
	return
}

func (w workflowImpl)mapInitiativeCommunitiesToOptions(comms []strategy.StrategyInitiativeCommunity, teamID models.TeamID) (opts []wf.SelectorOption) {
	for _, each := range comms {
		opts = append(opts, wf.SelectorOption{
			Label: ui.PlainText(each.Name), 
			Value: each.ID,
		})
	}
	return
}

func (w workflowImpl)OnCommunitySelected(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	itype := IssueType(ctx.Data[issueTypeKey])
	tc := getTypeClass(itype)
	var commID string
	commID, err = wf.SelectedValue(ctx.Request)
	if err == nil {
		switch itype {
		case SObjective:
			ctx.Data[capCommIDKey] = commID
		case Initiative:
			ctx.Data[initCommIDKey] = commID
		default: 
			w.AdaptiveLogger.Warnf("Unexpected itype with selected community")
		}
		out, err = w.showDialog(ctx, tc.Empty(), DescriptionContext)
	}
	return
}

