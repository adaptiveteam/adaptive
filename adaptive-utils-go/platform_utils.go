package adaptive_utils_go

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	// "github.com/adaptiveteam/adaptive/engagement-builder/model"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"fmt"
	// "encoding/json"
	"github.com/nlopes/slack"
	// "github.com/nlopes/slack/slackevents"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// Platform encapsulates configuration required to send notifications
type Platform struct {
	Sns awsutils.SnsRequest
	PlatformNotificationTopic string
	Namespace string
	IsInteractiveDebugEnabled bool
}

// Publish sends message to slack user via sns topic
func (p Platform)Publish(msg models.PlatformSimpleNotification) {
	_, err := p.Sns.Publish(msg, p.PlatformNotificationTopic)
	core.ErrorHandler(err, p.Namespace, 
		fmt.Sprintf("Could not publish message to %s topic", p.PlatformNotificationTopic))
}

// PublishAll sends a few messsages
func (p Platform)PublishAll(notes []models.PlatformSimpleNotification) {
	for _, note := range notes {
		p.Publish(note)
	}
}
// RecoverGracefully sends an error message to Slack.
func (p Platform)RecoverGracefully(request slack.InteractionCallback) {
	if err := recover(); err != nil {
		p.Publish(InteractionCallbackSimpleResponse(request,
			fmt.Sprintf("Error: %s", err)))
	}
}

// Debug prints debug message. If configuration allows, the message is directly sent to chat
func (p Platform)Debug(request slack.InteractionCallback, message string) {
	msg := "Debug: " + message
	if p.IsInteractiveDebugEnabled {
		p.Publish(InteractionCallbackSimpleResponse(request, msg))
	}
	fmt.Println(msg)
}
// ErrorHandler handles an error and appends message.
func (p Platform)ErrorHandler(request slack.InteractionCallback, msg string, err error) {
	if err != nil {
		message := fmt.Sprintf("%s while serving request %s \n(%s)", msg, request.CallbackID, err.Error())
		if p.IsInteractiveDebugEnabled {
			p.Publish(InteractionCallbackSimpleResponse(request, "Error: " + message))
		}
		core.ErrorHandler(err, p.Namespace, message)
	}
}

// ActionNameRule is the routing rule based on action name
func ActionNameRule(request slack.InteractionCallback) string {
	return request.ActionCallback.AttachmentActions[0].Name
}

// SelectedOptionRule is the routing rule based on selected menu option
func SelectedOptionRule(request slack.InteractionCallback) string {
	return request.ActionCallback.AttachmentActions[0].SelectedOptions[0].Value
}

// CallbackActionRule is the routing rule based on action name inside MessageCallback
func (p Platform)CallbackActionRule(request slack.InteractionCallback) string {
	mc := MessageCallbackParseUnsafe(request.CallbackID, p.Namespace)
	return mc.Action
}

// DispatchInteractionCallback dispatches request using provided routing table
func (p Platform)DispatchInteractionCallback(r RequestHandlers) func (slack.InteractionCallback, common.DynamoDBConnection) {
	return func (request slack.InteractionCallback, conn common.DynamoDBConnection) {
		defer p.RecoverGracefully(request)
		notes, err2 := r.DispatchByRule(ActionNameRule)(request, conn)
		core.ErrorHandler(err2, "DispatchInteractionCallback", "DispatchByRule")
		p.PublishAll(notes)
	}
}

// DispatchDialogSubmission dispatches request using provided routing table
func (p Platform)DispatchDialogSubmission(r DialogSubmissionHandlers) func (slack.InteractionCallback, slack.DialogSubmissionCallback, common.DynamoDBConnection) {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn common.DynamoDBConnection) {
		defer p.RecoverGracefully(request)
		notes, err2 := r.DispatchByRule(p.CallbackActionRule)(request, dialog, conn)
		core.ErrorHandler(err2, "DispatchDialogSubmission", "DispatchByRule")
		p.PublishAll(notes)
	}
}

// DispatchDialogSubmissionByRule dispatches request using provided routing table
func (p Platform)DispatchDialogSubmissionByRule(r DialogSubmissionHandlers, rule RequestRoutingRule) func (slack.InteractionCallback, slack.DialogSubmissionCallback, common.DynamoDBConnection) {
	return func (request slack.InteractionCallback, dialog slack.DialogSubmissionCallback, conn common.DynamoDBConnection) {
		defer p.RecoverGracefully(request)
		notes, err2 := r.DispatchByRule(rule)(request, dialog, conn)
		core.ErrorHandler(err2, "DispatchDialogSubmissionByRule", "DispatchByRule")
		p.PublishAll(notes)
	}
}
