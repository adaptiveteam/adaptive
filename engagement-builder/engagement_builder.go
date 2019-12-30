package EngagementBuilder

import (
	"encoding/json"
	"github.com/adaptiveteam/engagement-builder/model"
	"time"
)

type Engagement interface {
	Message() model.Message
	ToJson() ([]byte, error)
}

type EngagementBuilder interface {
	Id(string) EngagementBuilder
	Build() Engagement
	ToEngagement() Engagement
	Text(string) EngagementBuilder
	WithAttachment(*model.Attachment) EngagementBuilder
	WithResponseType(string) EngagementBuilder
	WithEffectiveStartDate(time.Time) EngagementBuilder
	WithEffectiveEndDate(time.Time) EngagementBuilder
	WithTags([]string) EngagementBuilder
	WithSecondsToAnswer(int) EngagementBuilder
}

type engagement struct {
	id                 string
	text               string
	attachments        []model.Attachment
	responseType       string
	tags               []string
	secondsToAnswer    int
	effectiveStartDate time.Time
	effectiveEndDate   time.Time
}

func (e *engagement) Id(id string) EngagementBuilder {
	e.id = id
	return e
}

func (e *engagement) Text(text string) EngagementBuilder {
	e.text = text
	return e
}

func (e *engagement) WithResponseType(rType string) EngagementBuilder {
	e.responseType = rType
	return e
}

func (e *engagement) WithAttachment(attach *model.Attachment) EngagementBuilder {
	e.attachments = append(e.attachments, *attach)
	return e
}

func (e *engagement) WithEffectiveStartDate(date time.Time) EngagementBuilder {
	e.effectiveStartDate = date
	return e
}

func (e *engagement) WithEffectiveEndDate(date time.Time) EngagementBuilder {
	e.effectiveEndDate = date
	return e
}

func (e *engagement) WithTags(tags []string) EngagementBuilder {
	e.tags = tags
	return e
}

func (e *engagement) WithSecondsToAnswer(time int) EngagementBuilder {
	e.secondsToAnswer = time
	return e
}

func (e *engagement) Build() Engagement {
	return e
}

func (e *engagement) ToEngagement() Engagement {
	return e
}

func (e *engagement) Message() model.Message {
	return model.Message{Id: e.id, Text: e.text, Attachments: e.attachments, ResponseType: e.responseType,
		EffectiveStartDate: e.effectiveStartDate, EffectiveEndDate: e.effectiveEndDate, Tags: e.tags,
		SecondsToAnswer: e.secondsToAnswer}
}

func (e *engagement) ToJson() ([]byte, error) {
	return json.Marshal(e.Message())
}

func NewEngagementBuilder() EngagementBuilder {
	return &engagement{}
}
