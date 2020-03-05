package userObjective
// This file has been automatically generated by `adaptive/scripts`
// The changes will be overridden by the next automatic generation.

type FieldName string
const (
	ID FieldName = "id"
	PlatformID FieldName = "platform_id"
	UserID FieldName = "user_id"
	Name FieldName = "name"
	Description FieldName = "description"
	AccountabilityPartner FieldName = "accountability_partner"
	Accepted FieldName = "accepted"
	ObjectiveType FieldName = "type"
	StrategyAlignmentEntityID FieldName = "strategy_alignment_entity_id"
	StrategyAlignmentEntityType FieldName = "strategy_alignment_entity_type"
	Quarter FieldName = "quarter"
	Year FieldName = "year"
	CreatedDate FieldName = "created_date"
	ExpectedEndDate FieldName = "expected_end_date"
	Completed FieldName = "completed"
	PartnerVerifiedCompletion FieldName = "partner_verified_completion"
	CompletedDate FieldName = "completed_date"
	PartnerVerifiedCompletionDate FieldName = "partner_verified_completion_date"
	Comments FieldName = "comments"
	Cancelled FieldName = "cancelled"
)

type IndexName string
const (
	UserIDCompletedIndex IndexName = "UserIDCompletedIndex"
	AcceptedIndex IndexName = "AcceptedIndex"
	AccountabilityPartnerIndex IndexName = "AccountabilityPartnerIndex"
	UserIDTypeIndex IndexName = "UserIDTypeIndex"
)