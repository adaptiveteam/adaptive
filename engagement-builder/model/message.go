package model

import "time"

type Message struct {
	Id          string       `json:"id"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// When replying to a parent message, this value is the value of the parent message to the thread
	ParentMessageId string `json:"parent_message_id,omitempty"`
	// For slack, this is `in_channel` or `ephemeral`
	ResponseType string `json:"response_type,omitempty"`
	// 	Used only when creating messages in response to a button action invocation. When set to true, the inciting message
	// 	will be replaced by this message you're providing. When false, the message you're providing is considered a brand new message
	ReplaceOriginal bool `json:"replace_original,omitempty"`
	// Used only when creating messages in response to a button action invocation. When set to true, the inciting message
	// will be deleted and if a message is provided, it will be posted as a brand new message.
	DeleteOriginal     bool      `json:"delete_original,omitempty"`
	EffectiveStartDate time.Time `json:"effective_start_date,omitempty"`
	EffectiveEndDate   time.Time `json:"effective_end_date,omitempty"`
	Tags               []string  `json:"tags,omitempty"`
	SecondsToAnswer    int       `json:"time_to_answer,omitempty"`
}
