package strategy

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"time"
)

const (
	InitiativeCommunityName                = "initiative_community_name"
	InitiativeCommunityDescription         = "initiative_community_description"
	InitiativeCommunityCoordinator         = "initiative_community_coordinator"
	InitiativeCommunityCapabilityCommunity = "initiative_community_capability_community"

	InitiativeCommunityNameLabel                = "Name"
	InitiativeCommunityDescriptionLabel         = "Description"
	InitiativeCommunityCoordinatorLabel         = "Coordinator"
	InitiativeCommunityCapabilityCommunityLabel = "Capability Community"

	Create models.AttachActionName = "create"
	Delete models.AttachActionName = "delete"
	Update models.AttachActionName = "update"

	InitiativeEvent                                 = "initiative"
	AssociateInitiativeWithInitiativeCommunityEvent = "associate_initiative_with_initiative_community"
	InitiativeCommunityAdhocEvent                   = "initiative_community_adhoc"
)

var (
	CreatePrefix = fmt.Sprintf("%s%s", string(Create), core.Underscore)
	DeletePrefix = fmt.Sprintf("%s%s", string(Delete), core.Underscore)
	UpdatePrefix = fmt.Sprintf("%s%s", string(Update), core.Underscore)
)

func communityEditStatus(si *StrategyInitiativeCommunity) string {
	return core.IfThenElse(si != nil, "updated", "created").(string)
}

func initiativeCommunityAttachmentFields(mc models.MessageCallback, oldSi, newSi *StrategyInitiativeCommunity,
	capabilityCommunitiesTable string) ([]models.KvPair, ui.RichText) {
	var kvs []models.KvPair
	dn := common.TaggedUser(newSi.Advocate)
	platformID := UserIDToPlatformID(userDAO())(mc.Source)
	if oldSi != nil {
		oldDn := common.TaggedUser(oldSi.Advocate)

		newCapComm := CapabilityCommunityByID(platformID, newSi.CapabilityCommunityID, capabilityCommunitiesTable)
		oldCapComm := CapabilityCommunityByID(platformID, oldSi.CapabilityCommunityID, capabilityCommunitiesTable)
		kvs = []models.KvPair{
			{Key: InitiativeCommunityNameLabel, Value: NewAndOld(newSi.Name, oldSi.Name)},
			{Key: InitiativeCommunityDescriptionLabel, Value: NewAndOld(newSi.Description, oldSi.Description)},
			{Key: InitiativeCommunityCoordinatorLabel, Value: NewAndOld(dn, oldDn)},
			{Key: InitiativeCommunityCapabilityCommunityLabel, Value: NewAndOld(newCapComm.Name, oldCapComm.Name)},
		}
	} else {
		newCapComm := CapabilityCommunityByID(platformID, newSi.CapabilityCommunityID, capabilityCommunitiesTable)
		kvs = []models.KvPair{
			{Key: InitiativeCommunityNameLabel, Value: newSi.Name},
			{Key: InitiativeCommunityDescriptionLabel, Value: newSi.Description},
			{Key: InitiativeCommunityCoordinatorLabel, Value: dn},
			{Key: InitiativeCommunityCapabilityCommunityLabel, Value: newCapComm.Name},
		}
	}
	return kvs, communityEditMessage(community.Initiative, communityEditStatus(oldSi))
}

func InitiativeCommunityViewAttachmentReadOnly(mc models.MessageCallback, newSi, oldSi *StrategyInitiativeCommunity,
	capabilityCommunitiesTable string) []ebm.Attachment {
	kvs, title := initiativeCommunityAttachmentFields(mc, oldSi, newSi, capabilityCommunitiesTable)
	return EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: title, Fields: kvs})
}

func initiativeCommunityEditActions(initCommID string, mc models.MessageCallback, strategyInitiativesTable, strategyInitiativesPlatformIndex string) []ebm.AttachmentAction {
	var actions []ebm.AttachmentAction
	platformID := UserIDToPlatformID(userDAO())(mc.Source)
	allInits := AllStrategyInitiatives(platformID, strategyInitiativesTable, strategyInitiativesPlatformIndex)
	mc = *mc.WithTarget(initCommID)
	if len(allInits) > 0 {
		// Show allocation button when there are initiatives
		actions = append(actions, AllocateInitiativeForCommunityAttachAction(mc, initCommID))
	}
	// actions = append(actions, AddInitiativeAttachAction(mc))
	return actions
}

func InitiativeCommunityViewAttachmentEditable(mc models.MessageCallback, newSi, oldSi *StrategyInitiativeCommunity,
	capabilityCommunitiesTable string, strategyInitiativesTable, strategyInitiativesPlatformIndex string) []ebm.Attachment {
	kvs, title := initiativeCommunityAttachmentFields(mc, oldSi, newSi, capabilityCommunitiesTable)

	editActions := initiativeCommunityEditActions(newSi.ID, mc, strategyInitiativesTable, strategyInitiativesPlatformIndex)
	actions := EditAttachActions(mc, newSi.ID, true, true, false, InitiativeCommunityAdhocEvent, editActions...)
	return EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: title, Actions: actions, Fields: kvs})
}

func AllocateInitiativeForCommunityAttachAction(mc models.MessageCallback, initCommID string) ebm.AttachmentAction {
	return *models.SimpleAttachAction(
		*mc.WithAction(string(Create)).
			WithTopic(AssociateInitiativeWithInitiativeCommunityEvent).
			WithTarget(initCommID),
		models.Now,
		"Allocate Initiative",
	)
}

func AddInitiativeAttachAction(mc models.MessageCallback) ebm.AttachmentAction {
	return *models.SimpleAttachAction(
		*mc.WithAction(string(Create)).WithTopic(InitiativeEvent),
		models.Now,
		"Create Initiative",
	)
}

func EditAttachActions(mc models.MessageCallback, id string, addAction, editAction, deleteAction bool, addNewTopic string,
	additionalActions ...ebm.AttachmentAction) []ebm.AttachmentAction {
	var actions []ebm.AttachmentAction
	if editAction {
		actions = append(actions, []ebm.AttachmentAction{
			*models.SimpleAttachAction(
				*mc.WithAction(string(Update)).WithTarget(id),
				models.Now,
				"Edit",
			),
		}...)
	}
	if addAction {
		actions = append(actions,
			*models.SimpleAttachAction(
				*mc.WithTopic(addNewTopic).WithAction(string(Create)).WithTarget(""),
				models.Now,
				"Add another?",
			))
	}
	if deleteAction {
		actions = append(actions,
			*models.GenAttachAction(
				*mc.WithAction(string(Delete)).WithTarget(id),
				models.Now,
				"Remove", models.EmptyActionConfirm(), true))
	}
	if addAction {
		actions = append(actions, additionalActions...)
	}
	return actions
}

func EntityViewAttachment(ae common.AttachmentEntity) []ebm.Attachment {
	attach := utils.ChatAttachment(string(ae.Title), string(ae.Text), string(ae.Fallback), ae.MC.ToCallbackID(), ae.Actions,
		models.AttachmentFields(ae.Fields), time.Now().Unix())
	return []ebm.Attachment{*attach}
}
