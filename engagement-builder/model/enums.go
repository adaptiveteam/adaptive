package model

type MarkdownField string
type AttachmentActionType string
type AttachmentActionStyle string
type AttachmentActionDataSource string

const (
	MarkdownFieldText    MarkdownField = "text"
	MarkdownFieldPretext MarkdownField = "pretext"
	MarkdownFieldFields  MarkdownField = "fields"
)

const (
	AttachmentActionTypeButton AttachmentActionType = "button"
	AttachmentActionTypeSelect AttachmentActionType = "select"
)

const (
	AttachmentActionStyleDefault AttachmentActionStyle = "default"
	AttachmentActionStylePrimary AttachmentActionStyle = "primary"
	AttachmentActionStyleDanger  AttachmentActionStyle = "danger"
)

const (
	AttachmentActionDataSourceUsers         AttachmentActionDataSource = "users"
	AttachmentActionDataSourceChannels      AttachmentActionDataSource = "channels"
	AttachmentActionDataSourceConversations AttachmentActionDataSource = "conversations"
)
