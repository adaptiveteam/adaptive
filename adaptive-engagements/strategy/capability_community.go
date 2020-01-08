package strategy

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

const (
	CapabilityCommunityName        = "capability_community_name"
	CapabilityCommunityDescription = "capability_community_description"
	CapabilityCommunityCoordinator = "capability_community_coordinator"

	CapabilityCommunityNameLabel        = "Name"
	CapabilityCommunityDescriptionLabel = "Description"
	CapabilityCommunityCoordinatorLabel = "Coordinator"

	CapabilityCommunityAdhocEvent = "capability_community_adhoc"
)

func CapabilityCommunityViewAttachment(mc models.MessageCallback, newCc, oldCc *CapabilityCommunity,
	enableActions bool) []ebm.Attachment {
	editStatus := "created"
	var actions []ebm.AttachmentAction
	var kvs []models.KvPair

	if oldCc != nil {
		editStatus = "updated"
		kvs = []models.KvPair{
			{Key: CapabilityCommunityNameLabel, Value: NewAndOld(newCc.Name, oldCc.Name)},
			{Key: CapabilityCommunityDescriptionLabel, Value: NewAndOld(newCc.Description, oldCc.Description)},
			{Key: CapabilityCommunityCoordinatorLabel,
				Value: NewAndOld(common.TaggedUser(newCc.Advocate), common.TaggedUser(oldCc.Advocate))},
		}
	} else {
		kvs = []models.KvPair{
			{Key: CapabilityCommunityNameLabel, Value: newCc.Name},
			{Key: CapabilityCommunityDescriptionLabel, Value: newCc.Description},
			{Key: CapabilityCommunityCoordinatorLabel, Value: common.TaggedUser(newCc.Advocate)},
		}
	}
	if enableActions {
		var extraActions []ebm.AttachmentAction
		mc = *mc.WithTarget(newCc.ID)
		actions = append(actions, EditAttachActions(mc, newCc.ID, true, true, false,
			CapabilityCommunityAdhocEvent, extraActions...)...)
	}
	return EntityViewAttachment(common.AttachmentEntity{MC: mc, Title: communityEditMessage(community.Capability, editStatus),
		Actions: actions, Fields: kvs})
}
