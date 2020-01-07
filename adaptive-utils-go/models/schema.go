package models

// Schema describes tables and indices in DynamoDB
type Schema struct {
	Holidays HolidaysTableSchema
	AdaptiveValues AdaptiveValuesTableSchema
	ClientPlatformTokens ClientPlatformTokenTableSchema 
	AdaptiveUsers AdaptiveUsersTableSchema
	UserFeedback UserFeedbackTableSchema
	CommunityUsers CommunityUsersTableSchema
	AdaptiveCommunity AdaptiveCommunityTableSchema
}

// SchemaForClientID creates comprehensive schema for a given clientID
func SchemaForClientID(clientID string) Schema {
	return Schema{
		Holidays: HolidaysTableSchemaForClientID(clientID),
		AdaptiveValues: AdaptiveValuesTableSchemaForClientID(clientID),
		ClientPlatformTokens: ClientPlatformTokenTableSchemaForClientID(clientID),
		AdaptiveUsers: AdaptiveUsersTableSchemaForClientID(clientID),
		UserFeedback: UserFeedbackTableSchemaForClientID(clientID),
		CommunityUsers: CommunityUsersTableSchemaForClientID(clientID),
		AdaptiveCommunity: AdaptiveCommunityTableSchemaForClientID(clientID),
	}
}
