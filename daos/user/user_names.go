package user
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	PlatformID FieldName = "platform_id"
	ID FieldName = "id"
	DisplayName FieldName = "display_name"
	FirstName FieldName = "first_name"
	LastName FieldName = "last_name"
	Timezone FieldName = "timezone"
	IsAdaptiveBot FieldName = "is_adaptive_bot"
	TimezoneOffset FieldName = "timezone_offset"
	AdaptiveScheduledTime FieldName = "adaptive_scheduled_time"
	AdaptiveScheduledTimeInUTC FieldName = "adaptive_scheduled_time_in_utc"
	PlatformOrg FieldName = "platform_org"
	IsAdmin FieldName = "is_admin"
	IsShared FieldName = "is_shared"
)

type IndexName string
const (
	PlatformIDIndex IndexName = "PlatformIDIndex"
	PlatformIDTimezoneOffsetIndex IndexName = "PlatformIDTimezoneOffsetIndex"
	PlatformIDAdaptiveScheduledTimeInUTCIndex IndexName = "PlatformIDAdaptiveScheduledTimeInUTCIndex"
)
