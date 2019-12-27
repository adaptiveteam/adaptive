package EngagementBuilder

import "github.com/adaptiveteam/engagement-builder/model"

// Attachment builder pattern code
type AttachmentBuilder struct {
	attachment *model.Attachment
}

func NewAttachmentBuilder() *AttachmentBuilder {
	attachment := &model.Attachment{}
	b := &AttachmentBuilder{attachment: attachment}
	return b
}

func LoadAttachmentBuilder(attach *model.Attachment) *AttachmentBuilder {
	return &AttachmentBuilder{attachment: attach}
}

func (b *AttachmentBuilder) Title(title string) *AttachmentBuilder {
	b.attachment.Title = title
	return b
}

func (b *AttachmentBuilder) Text(text string) *AttachmentBuilder {
	b.attachment.Text = text
	return b
}

func (b *AttachmentBuilder) Pretext(text string) *AttachmentBuilder {
	b.attachment.Pretext = text
	return b
}

func (b *AttachmentBuilder) Fallback(fallback string) *AttachmentBuilder {
	b.attachment.Fallback = fallback
	return b
}

func (b *AttachmentBuilder) CallbackId(callback string) *AttachmentBuilder {
	b.attachment.CallbackId = callback
	return b
}

func (b *AttachmentBuilder) Identifier(identifier string) *AttachmentBuilder {
	b.attachment.Identifier = identifier
	return b
}

func (b *AttachmentBuilder) Color(color string) *AttachmentBuilder {
	b.attachment.Color = color
	return b
}

func (b *AttachmentBuilder) Fields(fields []model.AttachmentField) *AttachmentBuilder {
	b.attachment.Fields = fields
	return b
}

func (b *AttachmentBuilder) Actions(actions []model.AttachmentAction) *AttachmentBuilder {
	b.attachment.Actions = actions
	return b
}

func (b *AttachmentBuilder) AttachmentType(attachmentType string) *AttachmentBuilder {
	b.attachment.AttachmentType = attachmentType
	return b
}

func (b *AttachmentBuilder) Footer(footer model.AttachmentFooter) *AttachmentBuilder {
	b.attachment.Footer = footer
	return b
}

func (b *AttachmentBuilder) TitleLink(titleLink string) *AttachmentBuilder {
	b.attachment.TitleLink = titleLink
	return b
}

func (b *AttachmentBuilder) ImageUrl(imageUrl string) *AttachmentBuilder {
	b.attachment.ImageUrl = imageUrl
	return b
}

func (b *AttachmentBuilder) ThumbUrl(thumbUrl string) *AttachmentBuilder {
	b.attachment.ThumbUrl = thumbUrl
	return b
}

func (b *AttachmentBuilder) MarkDownIn(fields []model.MarkdownField) *AttachmentBuilder {
	b.attachment.MrkdwnIn = fields
	return b
}

func (b *AttachmentBuilder) Author(author model.AttachmentAuthor) *AttachmentBuilder {
	b.attachment.Author = author
	return b
}

// Build converts builder to Attachment.
// deprecated. Use ToAttachment
// This method is misleading as it always returns `nil` as the second 
func (b *AttachmentBuilder) Build() (*model.Attachment, error) {
	return b.attachment, nil
}

func (b AttachmentBuilder) ToAttachment() model.Attachment {
	return *b.attachment
}
