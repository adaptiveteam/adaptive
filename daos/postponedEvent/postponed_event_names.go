package postponedEvent
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	ID FieldName = "id"
	UserID FieldName = "user_id"
	PlatformID FieldName = "platform_id"
	ActionPath FieldName = "action_path"
	ValidThrough FieldName = "valid_through"
)

type IndexName string
const (
	PlatformIDUserIDIndex IndexName = "PlatformIDUserIDIndex"
	UserIDIndex IndexName = "UserIDIndex"
)