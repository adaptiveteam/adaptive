package models

type ScriptRequest struct {
	Username string `json:"username"`
}

type ScriptResponse struct {
	Username string           `json:"username"`
	Script   []ScriptQuestion `json:"script"`
	SetId    int              `json:"set_id"`
}

type ScriptQuestion struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

type ScriptAnswer struct {
	UserQuestion string `json:"user_question"`
	SetId        int    `json:"set_id"`
	Username     string `json:"username"`
	Question     string `json:"question"`
	Answered     bool   `json:"answered"`
	Answer       string `json:"answer,omitempty"`
}

type UserNotification struct {
	UserId string `json:"user_id"`
	Id     string `json:"id"`
	Text   string `json:"text"`
}

type UserFeedback struct {
	Id               string `json:"id"`
	Source           string `json:"source"`
	Target           string `json:"target"`
	ValueID          string `json:"value_id"`
	ConfidenceFactor string `json:"confidence_factor"`
	Feedback         string `json:"feedback"`
	QuarterYear      string `json:"quarter_year"`
	PlatformID       string `json:"platform_id"`
	// Channel, if any, to engage user in response to the feedback
	// This is useful to reply to an event with no knowledge of the previous context
	Channel string `json:"channel"`
	// A reference to the original timestamp that can be used to reply via threading
	MsgTimestamp string `json:"msg_timestamp"`
}

// UserFeedbackTableSchema is schema of Dynamo table with user info.
type UserFeedbackTableSchema struct {
	Name                           string
	FeedbackSourceQuarterYearIndex string
}

// UserFeedbackTableSchemaForClientID creates table schema given client id
func UserFeedbackTableSchemaForClientID(clientID string) UserFeedbackTableSchema {
	return UserFeedbackTableSchema{
		Name:                           clientID + "_adaptive_user_feedback",
		FeedbackSourceQuarterYearIndex: "SourceQuarterYear",
	}
}
