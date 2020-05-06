package community
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	PlatformID FieldName = "platform_id"
	ID FieldName = "id"
	ChannelID FieldName = "channel_id"
	CommunityKind FieldName = "community_kind"
	ParentCommunityID FieldName = "parent_community_id"
	Name FieldName = "name"
	Description FieldName = "description"
	Advocate FieldName = "advocate"
	AccountabilityPartner FieldName = "accountability_partner"
	CreatedBy FieldName = "created_by"
	ModifiedBy FieldName = "modified_by"
	RequestedBy FieldName = "requested_by"
)

type IndexName string
const (
	ChannelIDPlatformIDIndex IndexName = "ChannelIDPlatformIDIndex"
	PlatformIDCommunityKindIndex IndexName = "PlatformIDCommunityKindIndex"
)
