package EngagementBuilder

import (
	"encoding/json"
	"github.com/adaptiveteam/engagement-builder/model"
	"github.com/go-test/deep"
	"io/ioutil"
	"testing"
	"time"
)

func TestEngagementBuilder(t *testing.T) {
	optionGroup, _ := NewMenuOptionGroupBuilder().Text("Doggone bot antics").Options([]model.MenuOption{
		{
			Text:  "Unexpected sentience",
			Value: "AI-2323",
		},
	}).Build()
	menuOption1, _ := NewMenuOptionBuilder().Description("test").Text("Hearts").Value("hearts").Build()
	menuOption2, _ := NewMenuOptionBuilder().Description("test").Text("Bridge").Value("bridge").Build()
	fields := []model.AttachmentField{
		{
			Title: "Priority",
			Value: "High",
			Short: false,
		},
	}
	// attachment action builder
	attachAction, _ := NewAttachmentActionBuilder().
		Name("recommend").
		Text("Recommend").
		ActionType(model.AttachmentActionTypeButton).
		Value("recommend").
		Confirm(model.AttachmentActionConfirm{Text: "Really?"}).
		Style(model.AttachmentActionStyleDefault).
		Options([]model.MenuOption{*menuOption1, *menuOption2}).
		SelectedOptions(*menuOption1).
		OptionGroups([]model.MenuOptionGroup{*optionGroup}).
		Url("http://example.com/path/to/thumb.png").
		DataSource(model.AttachmentActionDataSourceUsers).
		Build()

	attach, _ := NewAttachmentBuilder().
		Fallback("Required plain-text summary of the attachment.").
		Color("#36a64f").
		Pretext("Optional text that appears above the attachment block").
		Author(model.AttachmentAuthor{
			Name: "Bobby Tables",
			Link: "http://flickr.com/bobby/",
			Icon: "http://flickr.com/icons/bobby.jpg",
		}).
		CallbackId("test_attach").
		Identifier("test_id").
		Title("Slack API Documentation").
		TitleLink("https://api.slack.com/").
		Text("Optional text that appears within the attachment").
		ImageUrl("http://my-website.com/path/to/image.jpg").
		ThumbUrl("http://example.com/path/to/thumb.png").
		Footer(model.AttachmentFooter{
			Text:      "Slack API",
			Icon:      "https://platform.slack-edge.com/img/default_application_icon.png",
			Timestamp: 123456789,
		}).
		Fields(fields).
		Actions([]model.AttachmentAction{*attachAction}).
		AttachmentType("default").
		MarkDownIn([]model.MarkdownField{model.MarkdownFieldText}).
		Build()

	engagement := NewEngagementBuilder().
		Id("test_engage").
		Text("Would you like to play a game?").
		WithAttachment(attach).
		WithResponseType("in_channel").
		WithEffectiveStartDate(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)).
		WithTags([]string{"test"}).
		WithSecondsToAnswer(1).
		Build()

	origBytes, err := engagement.ToJson()
	if err != nil {
		t.Fail()
	}
	var infMsg model.Message
	err = json.Unmarshal(origBytes, &infMsg)
	if err != nil {
		t.Fail()
	}

	// Comparing original struct with struct constructed from JSON
	if diff := deep.Equal(infMsg, engagement.Message()); diff != nil {
		t.Error(diff)
	}

	// Comparing with hand written JSON
	b, err := ioutil.ReadFile("test/engagement.json")
	if err != nil {
		t.Error(err)
	}

	var parsedJson model.Message
	err = json.Unmarshal(b, &parsedJson)
	if err != nil {
		t.Fail()
	}

	// Comparing original struct with struct constructed from JSON
	if diff := deep.Equal(parsedJson, engagement.Message()); diff != nil {
		t.Error(diff)
	}

	// Loading attachment builder
	newAttach, _ := LoadAttachmentBuilder(attach).Color("red").Build()
	if newAttach.Color != "red" {
		t.Fail()
	}
}
