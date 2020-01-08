package engagement_slack_mapper

import (
	"fmt"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/nlopes/slack"
	"testing"
)

const (
	token   = "xoxb-436528929141-492802537186-3tbN7QlbieTa27P6ROdsOoTj" // adaptive-team
	channel = "UE48A5TC0"
)

func postToSlack(options []slack.MsgOption) (string, error) {
	api := slack.New(token)
	_, ts, err := api.PostMessage(channel, options...)
	return ts, err
}

func checkPost(e eb.Engagement, t *testing.T) {
	se := SlackEngagement{Message: e.Message()}
	_, err := postToSlack(se.MsgOptions())
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}

//----------------------------- ATTACHMENTS -----------------------------//

func TestSimpleSlackMessage(t *testing.T) {
	attach, _ := eb.NewAttachmentBuilder().
		Text("And here’s an attachment!").
		Build()
	engagement := eb.NewEngagementBuilder().
		Text("I am a test message http://slack.com").
		WithAttachment(attach).
		Build()
	checkPost(engagement, t)
}

// From example here: https://api.slack.com/docs/messages/builder?msg=%7B%22attachments%22%3A%5B%7B%22fallback%22%3A%22Required%20plain-text%20summary%20of%20the%20attachment.%22%2C%22color%22%3A%22%2336a64f%22%2C%22pretext%22%3A%22Optional%20text%20that%20appears%20above%20the%20attachment%20block%22%2C%22author_name%22%3A%22Bobby%20Tables%22%2C%22author_link%22%3A%22http%3A%2F%2Fflickr.com%2Fbobby%2F%22%2C%22author_icon%22%3A%22http%3A%2F%2Fflickr.com%2Ficons%2Fbobby.jpg%22%2C%22title%22%3A%22Slack%20API%20Documentation%22%2C%22title_link%22%3A%22https%3A%2F%2Fapi.slack.com%2F%22%2C%22text%22%3A%22Optional%20text%20that%20appears%20within%20the%20attachment%22%2C%22fields%22%3A%5B%7B%22title%22%3A%22Priority%22%2C%22value%22%3A%22High%22%2C%22short%22%3Afalse%7D%5D%2C%22image_url%22%3A%22http%3A%2F%2Fmy-website.com%2Fpath%2Fto%2Fimage.jpg%22%2C%22thumb_url%22%3A%22http%3A%2F%2Fexample.com%2Fpath%2Fto%2Fthumb.png%22%2C%22footer%22%3A%22Slack%20API%22%2C%22footer_icon%22%3A%22https%3A%2F%2Fplatform.slack-edge.com%2Fimg%2Fdefault_application_icon.png%22%2C%22ts%22%3A123456789%7D%5D%7D
func TestSlackMessageAttachments(t *testing.T) {
	fields := []model.AttachmentField{
		{
			Title: "Priority",
			Value: "High",
			Short: false,
		},
	}

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Required plain-text summary of the attachment.").
		Color("#36a64f").
		Pretext("Optional text that appears above the attachment block").
		Author(model.AttachmentAuthor{
			Name: "Bobby Tables",
			Link: "http://flickr.com/bobby/",
			Icon: "http://flickr.com/icons/bobby.jpg",
		}).
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
		Build()
	engagement := eb.NewEngagementBuilder().
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

//----------------------------- BUTTONS -----------------------------//

// From example here: https://api.slack.com/docs/messages/builder?msg=%7B%22text%22%3A%22Would%20you%20like%20to%20play%20a%20game%3F%22%2C%22attachments%22%3A%5B%7B%22text%22%3A%22Choose%20a%20game%20to%20play%22%2C%22fallback%22%3A%22You%20are%20unable%20to%20choose%20a%20game%22%2C%22callback_id%22%3A%22wopr_game%22%2C%22color%22%3A%22%233AA3E3%22%2C%22attachment_type%22%3A%22default%22%2C%22actions%22%3A%5B%7B%22name%22%3A%22game%22%2C%22text%22%3A%22Chess%22%2C%22type%22%3A%22button%22%2C%22value%22%3A%22chess%22%7D%2C%7B%22name%22%3A%22game%22%2C%22text%22%3A%22Falken%27s%20Maze%22%2C%22type%22%3A%22button%22%2C%22value%22%3A%22maze%22%7D%2C%7B%22name%22%3A%22game%22%2C%22text%22%3A%22Thermonuclear%20War%22%2C%22style%22%3A%22danger%22%2C%22type%22%3A%22button%22%2C%22value%22%3A%22war%22%2C%22confirm%22%3A%7B%22title%22%3A%22Are%20you%20sure%3F%22%2C%22text%22%3A%22Wouldn%27t%20you%20prefer%20a%20good%20game%20of%20chess%3F%22%2C%22ok_text%22%3A%22Yes%22%2C%22dismiss_text%22%3A%22No%22%7D%7D%5D%7D%5D%7D
func TestSlackMessageButtons(t *testing.T) {
	attachAction1, _ := eb.NewAttachmentActionBuilder().
		Name("game").
		Text("Chess").
		ActionType(model.AttachmentActionTypeButton).
		Value("chess").
		Build()
	attachAction2, _ := eb.NewAttachmentActionBuilder().
		Name("game").
		Text("Falken's Maze").
		ActionType("button").
		Value("maze").
		Build()
	attachAction3, _ := eb.NewAttachmentActionBuilder().
		Name("game").
		Text("Thermonuclear War").
		ActionType("button").
		Value("war").
		Style(model.AttachmentActionStyleDanger).
		Confirm(model.AttachmentActionConfirm{
			Title:       "Are you sure?",
			Text:        "Wouldn't you prefer a good game of chess?",
			OkText:      "Yes",
			DismissText: "No",
		}).
		Build()
	attach1, _ := eb.NewAttachmentBuilder().
		Text("Choose a game to play").
		Fallback("You are unable to choose a game").
		Color("#3AA3E3").
		AttachmentType("default").
		Actions([]model.AttachmentAction{*attachAction1, *attachAction2, *attachAction3}).
		Build()
	engagement := eb.NewEngagementBuilder().
		Text("Would you like to play a game?").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

//----------------------------- FIELDS -----------------------------//

// Example from here: https://api.slack.com/docs/messages/builder?msg=%7B%22attachments%22%3A%5B%7B%22fallback%22%3A%22Required%20plain-text%20summary%20of%20the%20attachment.%22%2C%22color%22%3A%22%2336a64f%22%2C%22pretext%22%3A%22Optional%20text%20that%20appears%20above%20the%20attachment%20block%22%2C%22author_name%22%3A%22Bobby%20Tables%22%2C%22author_link%22%3A%22http%3A%2F%2Fflickr.com%2Fbobby%2F%22%2C%22author_icon%22%3A%22http%3A%2F%2Fflickr.com%2Ficons%2Fbobby.jpg%22%2C%22title%22%3A%22Slack%20API%20Documentation%22%2C%22title_link%22%3A%22https%3A%2F%2Fapi.slack.com%2F%22%2C%22text%22%3A%22Optional%20text%20that%20appears%20within%20the%20attachment%22%2C%22fields%22%3A%5B%7B%22title%22%3A%22Priority%22%2C%22value%22%3A%22High%22%2C%22short%22%3Afalse%7D%5D%2C%22image_url%22%3A%22http%3A%2F%2Fmy-website.com%2Fpath%2Fto%2Fimage.jpg%22%2C%22thumb_url%22%3A%22http%3A%2F%2Fexample.com%2Fpath%2Fto%2Fthumb.png%22%2C%22footer%22%3A%22Slack%20API%22%2C%22footer_icon%22%3A%22https%3A%2F%2Fplatform.slack-edge.com%2Fimg%2Fdefault_application_icon.png%22%2C%22ts%22%3A123456789%7D%5D%7D
func TestSlackMessageFields(t *testing.T) {
	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Required plain-text summary of the attachment.").
		Color("#36a64f").
		Pretext("Optional text that appears above the attachment block").
		Author(model.AttachmentAuthor{Name: "Bobby Tables", Link: "http://flickr.com/bobby/", Icon: "http://flickr.com/icons/bobby.jpg"}).
		Title("Slack API Documentation").
		TitleLink("https://api.slack.com/").
		Text("Optional text that appears within the attachment").
		Fields([]model.AttachmentField{
			{
				Title: "Priority",
				Value: "High",
				Short: false,
			},
		}).ImageUrl("http://my-website.com/path/to/image.jpg").
		ThumbUrl("http://example.com/path/to/thumb.png").
		Footer(model.AttachmentFooter{Text: "Slack API", Icon: "https://platform.slack-edge.com/img/default_application_icon.png", Timestamp: 123456789}).
		Build()

	engagement := eb.NewEngagementBuilder().
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

// From example here: https://api.slack.com/docs/messages/builder?msg=%7B%22text%22%3A%22New%20comic%20book%20alert!%22%2C%22attachments%22%3A%5B%7B%22title%22%3A%22The%20Further%20Adventures%20of%20Slackbot%22%2C%22fields%22%3A%5B%7B%22title%22%3A%22Volume%22%2C%22value%22%3A%221%22%2C%22short%22%3Atrue%7D%2C%7B%22title%22%3A%22Issue%22%2C%22value%22%3A%223%22%2C%22short%22%3Atrue%7D%5D%2C%22author_name%22%3A%22Stanford%20S.%20Strickland%22%2C%22author_icon%22%3A%22http%3A%2F%2Fa.slack-edge.com%2F7f18%2Fimg%2Fapi%2Fhomepage_custom_integrations-2x.png%22%2C%22image_url%22%3A%22http%3A%2F%2Fi.imgur.com%2FOJkaVOI.jpg%3F1%22%7D%2C%7B%22title%22%3A%22Synopsis%22%2C%22text%22%3A%22After%20%40episod%20pushed%20exciting%20changes%20to%20a%20devious%20new%20branch%20back%20in%20Issue%201%2C%20Slackbot%20notifies%20%40don%20about%20an%20unexpected%20deploy...%22%7D%2C%7B%22fallback%22%3A%22Would%20you%20recommend%20it%20to%20customers%3F%22%2C%22title%22%3A%22Would%20you%20recommend%20it%20to%20customers%3F%22%2C%22callback_id%22%3A%22comic_1234_xyz%22%2C%22color%22%3A%22%233AA3E3%22%2C%22attachment_type%22%3A%22default%22%2C%22actions%22%3A%5B%7B%22name%22%3A%22recommend%22%2C%22text%22%3A%22Recommend%22%2C%22type%22%3A%22button%22%2C%22value%22%3A%22recommend%22%7D%2C%7B%22name%22%3A%22no%22%2C%22text%22%3A%22No%22%2C%22type%22%3A%22button%22%2C%22value%22%3A%22bad%22%7D%5D%7D%5D%7D
func TestSlackMessageFieldsWithButtons(t *testing.T) {
	attach1, _ := eb.NewAttachmentBuilder().
		Text("The Further Adventures of Slackbot").
		Fields([]model.AttachmentField{
			{
				Title: "Volume",
				Value: "1",
				Short: true,
			},
			{
				Title: "Issue",
				Value: "3",
				Short: true,
			},
		}).
		Author(model.AttachmentAuthor{
			Name: "Stanford S. Strickland",
			Icon: "http://a.slack-edge.com/7f18/img/api/homepage_custom_integrations-2x.png",
		}).
		ImageUrl("http://i.imgur.com/OJkaVOI.jpg?1").
		Build()

	attach2, _ := eb.NewAttachmentBuilder().
		Title("Synopsis").
		Text("After @episod pushed exciting changes to a devious new branch back in Issue 1, Slackbot notifies @don about an unexpected deploy...").
		Build()

	attach3Action1, _ := eb.NewAttachmentActionBuilder().
		Name("recommend").
		Text("Recommend").
		ActionType(model.AttachmentActionTypeButton).
		Value("cherecommendss").
		Build()
	attach3Action2, _ := eb.NewAttachmentActionBuilder().
		Name("no").
		Text("No").
		ActionType(model.AttachmentActionTypeButton).
		Value("bad").
		Build()

	attach3, _ := eb.NewAttachmentBuilder().
		Fallback("Would you recommend it to customers?").
		Title("Would you recommend it to customers?").
		CallbackId("comic_1234_xyz").
		Color("#3AA3E3").
		AttachmentType("default").
		Actions([]model.AttachmentAction{*attach3Action1, *attach3Action2}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("New comic book alert!").
		WithAttachment(attach1).
		WithAttachment(attach2).
		WithAttachment(attach3).
		Build()
	checkPost(engagement, t)
}

//----------------------------- URLs -----------------------------//

// Example from here: https://api.slack.com/docs/messages/builder?msg=%7B%22attachments%22%3A%5B%7B%22fallback%22%3A%22Network%20traffic%20(kb%2Fs)%3A%20How%20does%20this%20look%3F%20%40slack-ops%20-%20Sent%20by%20Julie%20Dodd%20-%20https%3A%2F%2Fdatadog.com%2Fpath%2Fto%2Fevent%22%2C%22title%22%3A%22Network%20traffic%20(kb%2Fs)%22%2C%22title_link%22%3A%22https%3A%2F%2Fdatadog.com%2Fpath%2Fto%2Fevent%22%2C%22text%22%3A%22How%20does%20this%20look%3F%20%40slack-ops%20-%20Sent%20by%20Julie%20Dodd%22%2C%22image_url%22%3A%22https%3A%2F%2Fdatadoghq.com%2Fsnapshot%2Fpath%2Fto%2Fsnapshot.png%22%2C%22color%22%3A%22%23764FA5%22%7D%5D%7D
func TestSlackMessageWithImage(t *testing.T) {

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Network traffic (kb/s): How does this look? @slack-ops - Sent by Julie Dodd - https://datadog.com/path/to/event").
		Title("Network traffic (kb/s)").
		TitleLink("https://datadog.com/path/to/event").
		Text("How does this look? @slack-ops - Sent by Julie Dodd").
		ImageUrl("https://datadoghq.com/snapshot/path/to/snapshot.png").
		Build()

	engagement := eb.NewEngagementBuilder().
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

// Example from here: https://api.slack.com/docs/messages/builder?msg=%7B%22text%22%3A%22%3C%40W1A2BC3DD%3E%20approved%20your%20travel%20request.%20Book%20any%20airline%20you%20like%20by%20continuing%20below.%22%2C%22channel%22%3A%22C061EG9SL%22%2C%22attachments%22%3A%5B%7B%22fallback%22%3A%22Book%20your%20flights%20at%20https%3A%2F%2Fflights.example.com%2Fbook%2Fr123456%22%2C%22actions%22%3A%5B%7B%22type%22%3A%22button%22%2C%22text%22%3A%22Book%20flights%20%F0%9F%9B%AB%22%2C%22url%22%3A%22https%3A%2F%2Fflights.example.com%2Fbook%2Fr123456%22%7D%5D%7D%5D%7D
func TestSlackWithExternalUrl(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Text("Book flights 🛫").
		ActionType(model.AttachmentActionTypeButton).
		Url("https://flights.example.com/book/r123456").
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Book your flights at https://flights.example.com/book/r123456").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("<@W1A2BC3DD> approved your travel request. Book any airline you like by continuing below.").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

// Example from here: https://api.slack.com/docs/messages/builder?msg=%7B%22text%22%3A%22%3C%40W1A2BC3DD%3E%20approved%20your%20travel%20request.%20Book%20any%20airline%20you%20like%20by%20continuing%20below.%22%2C%22channel%22%3A%22C061EG9SL%22%2C%22attachments%22%3A%5B%7B%22fallback%22%3A%22Book%20your%20flights%20at%20https%3A%2F%2Fflights.example.com%2Fbook%2Fr123456%22%2C%22actions%22%3A%5B%7B%22type%22%3A%22button%22%2C%22name%22%3A%22travel_request_123456%22%2C%22text%22%3A%22Book%20flights%20%F0%9F%9B%AB%22%2C%22url%22%3A%22https%3A%2F%2Fflights.example.com%2Fbook%2Fr123456%22%2C%22style%22%3A%22primary%22%2C%22confirm%22%3A%22Really%3F%22%7D%2C%7B%22type%22%3A%22button%22%2C%22name%22%3A%22travel_cancel_123456%22%2C%22text%22%3A%22Cancel%20travel%20request%22%2C%22url%22%3A%22https%3A%2F%2Frequests.example.com%2Fcancel%2Fr123456%22%2C%22style%22%3A%22danger%22%7D%5D%7D%5D%7D
func TestSlackWithButtonsExternalUrl(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("travel_request_123456").
		Text("Book flights 🛫").
		ActionType(model.AttachmentActionTypeButton).
		Url("https://flights.example.com/book/r123456").
		Style(model.AttachmentActionStylePrimary).
		Confirm(model.AttachmentActionConfirm{Text: "Really?"}).
		Build()

	attach1Action2, _ := eb.NewAttachmentActionBuilder().
		Name("travel_cancel_123456").
		Text("Cancel travel request").
		ActionType(model.AttachmentActionTypeButton).
		Url("https://requests.example.com/cancel/r123456").
		Style(model.AttachmentActionStyleDanger).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Book your flights at https://flights.example.com/book/r123456").
		Actions([]model.AttachmentAction{*attach1Action1, *attach1Action2}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("<@W1A2BC3DD> approved your travel request. Book any airline you like by continuing below.").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

//----------------------------- MENUS -----------------------------//
// From examples here: https://api.slack.com/docs/message-menus

func TestSlackWithSimpleMenu(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("games_list").
		Text("Pick a game...").
		ActionType(model.AttachmentActionTypeSelect).
		Options([]model.MenuOption{
			{
				Text:  "Hearts",
				Value: "Hearts",
			},
			{
				Text:  "Bridge",
				Value: "bridge",
			},
			{
				Text:  "Checkers",
				Value: "checkers",
			},
			{
				Text:  "Chess",
				Value: "chess",
			},
			{
				Text:  "Poker",
				Value: "poker",
			},
			{
				Text:  "Falken's Maze",
				Value: "maze",
			},
			{
				Text:  "Global Thermonuclear War",
				Value: "war",
			},
		}).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Text("Choose a game to play").
		Fallback("If you could read this message, you'd be choosing something fun to do right now.").
		Color("#3AA3E3").
		AttachmentType("default").
		CallbackId("game_selection").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("Would you like to play a game?").
		WithResponseType("in_channel").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

func TestSlackMenuWithUsersSource(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("winners_list").
		Text("Who should win?").
		ActionType(model.AttachmentActionTypeSelect).
		DataSource(model.AttachmentActionDataSourceUsers).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Text("Who wins the lifetime supply of chocolate?").
		Fallback("You could be telling the computer exactly what it can do with a lifetime supply of chocolate.").
		Color("#3AA3E3").
		AttachmentType("default").
		CallbackId("select_simple_1234").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("I hope the tour went well, Mr. Wonka.").
		WithResponseType("in_channel").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

func TestSlackMenuWithChannelSource(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("channels_list").
		Text("Which channel changed your life this week?").
		ActionType(model.AttachmentActionTypeSelect).
		DataSource(model.AttachmentActionDataSourceChannels).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Upgrade your Slack client to use messages like these.").
		Color("#3AA3E3").
		AttachmentType("default").
		CallbackId("select_simple_1234").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("It's time to nominate the channel of the week.").
		WithResponseType("in_channel").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

func TestSlackMenuWithConversationSource(t *testing.T) {
	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("conversations_list").
		Text("Who did you talk to last?").
		ActionType(model.AttachmentActionTypeSelect).
		DataSource(model.AttachmentActionDataSourceConversations).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Upgrade your Slack client to use messages like these.").
		Color("#3AA3E3").
		AttachmentType("default").
		CallbackId("conversations_123").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("Let's get a productive conversation going").
		WithResponseType("in_channel").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}

func TestSlackMenuWithOptionGroups(t *testing.T) {
	optionGroup1, _ := eb.NewMenuOptionGroupBuilder().Text("Doggone bot antics").Options([]model.MenuOption{
		{
			Text:  "Unexpected sentience",
			Value: "AI-2323",
		},
		{
			Text:  "Bot biased toward other bots",
			Value: "SUPPORT-42",
		},
		{
			Text:  "Bot broke my toaster",
			Value: "IOT-75",
		},
	}).Build()
	optionGroup2, _ := eb.NewMenuOptionGroupBuilder().Text("Human error").Options([]model.MenuOption{
		{
			Text:  "Not Penny's boat",
			Value: "LOST-7172",
		},
		{
			Text:  "We built our own CMS",
			Value: "OOPS-1",
		},
	}).Build()

	attach1Action1, _ := eb.NewAttachmentActionBuilder().
		Name("conversations_list").
		Text("Who did you talk to last?").
		ActionType(model.AttachmentActionTypeSelect).
		OptionGroups([]model.MenuOptionGroup{*optionGroup1, *optionGroup2}).
		Build()

	attach1, _ := eb.NewAttachmentBuilder().
		Fallback("Upgrade your Slack client to use messages like these.").
		Color("#3AA3E3").
		AttachmentType("default").
		CallbackId("conversations_123").
		Actions([]model.AttachmentAction{*attach1Action1}).
		Build()

	engagement := eb.NewEngagementBuilder().
		Text("This is an example of attachment groups").
		WithAttachment(attach1).
		Build()
	checkPost(engagement, t)
}