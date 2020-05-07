package engagement_slack_mapper

import (
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
)
// PlatformAdapter encapsulates mapper and provides a couple of helper functions
type PlatformAdapter interface {
	PostSync(psn models.PlatformSimpleNotification) (id MessageID, err error)
	PostAsync(psn models.PlatformSimpleNotification)
	PostSyncUnsafe(psn models.PlatformSimpleNotification) (id MessageID)
	PostManyAsync(psns ...models.PlatformSimpleNotification)
	PostManySync(psns ...models.PlatformSimpleNotification) (errs []error)
	PostManySyncUnsafe(psns ...models.PlatformSimpleNotification)
}

// SlackAPIAdapterImpl encapsulates information necessary to connect to Slack
// It implements PlatformAdapter interface
type SlackAPIAdapterImpl struct {
	API *slack.Client
}

// SlackAdapter constructs PlatformAdapter for Slack
func SlackAdapter(platformToken string) PlatformAdapter {
	return SlackAPIAdapterImpl{
		API: slack.New(platformToken),
	}
}

// PostAsync posts message and forgets about it.
func (s SlackAPIAdapterImpl)PostAsync(psn models.PlatformSimpleNotification) {
	core_utils_go.Go("PostSyncUnsafe", func(){ s.PostSyncUnsafe(psn)})
}

// MessageID identifies message that has been posted.
type MessageID struct {
	ConversationID platform.ConversationID
	Ts string
}

// PostSync sends message to Slack
func (s SlackAPIAdapterImpl)PostSync(psn models.PlatformSimpleNotification) (id MessageID, err error) {
	msgParams := ConvertPSN(psn)
	conversationID := GetConversationID(psn)
	channel, timestamp, err := s.API.PostMessage(conversationID, msgParams...)
	id = MessageID{
		ConversationID: platform.ConversationID(channel),
		Ts: timestamp,
	}
	return
}

// PostSyncUnsafe sends message to Slack. Panics in case of errors
func (s SlackAPIAdapterImpl)PostSyncUnsafe(psn models.PlatformSimpleNotification) (id MessageID) {
	id, err := s.PostSync(psn)

	if err != nil {
		fmt.Printf("Error posting message %s to Slack: %v", psn.Message, err)
		panic(err)
	}
	return id
}

// ConvertPSN converts PlatformSimpleNotification to a collection of MsgOptions
func ConvertPSN(psn models.PlatformSimpleNotification)(msgParams []slack.MsgOption){
	// Converting generic attachments to slack attachments
	var msgOption slack.MsgOption
	if len(psn.Attachments) > 0 {
		msgOption = slack.MsgOptionAttachments(Attachments(psn.Attachments)...)
	} else {
		msgOption = slack.MsgOptionAttachments(slack.Attachment{})
	}
	// base message param
	msgParams = []slack.MsgOption{
		slack.MsgOptionText(psn.Message, false),
		slack.MsgOptionAsUser(psn.AsUser), 
		msgOption,
	}

	if psn.ThreadTs != "" {
		msgParams = append(msgParams, slack.MsgOptionTS(psn.ThreadTs))
	}
	return 	
}

// GetConversationID retrieves conversation id from PlatformSimpleNotification
func GetConversationID(psn models.PlatformSimpleNotification) (postTo string) {
	if psn.Channel != "" {
		// When channel is set, post to that. Else, post to the user.
		return psn.Channel
	}
	return psn.UserId
}

// PostManyAsync posts a few messages in async mode
func (s SlackAPIAdapterImpl)PostManyAsync(psns ...models.PlatformSimpleNotification) {
	core_utils_go.Go("PostManySyncUnsafe", func(){s.PostManySyncUnsafe(psns...)})
}

// PostManySync posts a few messages synchronously
func (s SlackAPIAdapterImpl)PostManySync(psns ...models.PlatformSimpleNotification) (errs []error) {
	for _, psn := range psns {
		_, err := s.PostSync(psn)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

// PostManySyncUnsafe posts a few messages synchronously. Panics in case of any errors
func (s SlackAPIAdapterImpl)PostManySyncUnsafe(psns ...models.PlatformSimpleNotification) {
	errs := s.PostManySync(psns...)
	for _, err := range errs {
		fmt.Printf("Error posting multi-message to Slack: %v", err)
	}
	if len(errs) > 0 {
		panic(errs[0])
	}
}
