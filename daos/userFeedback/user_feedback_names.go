package userFeedback
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	ID FieldName = "id"
	Source FieldName = "source"
	Target FieldName = "target"
	ValueID FieldName = "value_id"
	ConfidenceFactor FieldName = "confidence_factor"
	Feedback FieldName = "feedback"
	QuarterYear FieldName = "quarter_year"
	ChannelID FieldName = "channel"
	MsgTimestamp FieldName = "msg_timestamp"
	PlatformID FieldName = "platform_id"
)

type IndexName string
const (
	QuarterYearSourceIndex IndexName = "QuarterYearSourceIndex"
	QuarterYearTargetIndex IndexName = "QuarterYearTargetIndex"
)