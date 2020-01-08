package adaptive_utils_go

import (
	"encoding/json"
	"errors"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
)

// UserProfileLambda encapsulates lambda that extracts user tokens
type UserProfileLambda struct {
	Region            string
	Namespace         string
	ProfileLambdaName string
}

// UserTokenSync gets token associated with a user by invoking profile lambda
// It performs synchronous request to the remote lambda.
func (u UserProfileLambda) UserTokenSync(userID string) (ut models.UserToken, err error) {
	return u.UserTokenAdvSync(userID, "")
}

// UserTokenAdvSync gets token associated with a user by invoking profile lambda
// It performs synchronous request to the remote lambda.
func (u UserProfileLambda) UserTokenAdvSync(userID string, platformID models.PlatformID) (ut models.UserToken, err error) {
	ut = models.UserToken{}
	if u.ProfileLambdaName == "" {
		err = errors.New("UserProfileLambda#ProfileLambdaName: Invalid argument, it should not be empty")
	} else {
		l := awsutils.NewLambda(u.Region, "", u.Namespace)
		userEngBytes, err := json.Marshal(models.UserEngage{UserId: userID, PlatformID: platformID})
		if err == nil {
			output, err := l.InvokeFunction(u.ProfileLambdaName, userEngBytes, false)
			if err == nil {
				err = json.Unmarshal(output.Payload, &ut)
			}
		}
	}
	return
}

// UserTokenSyncUnsafe gets token associated with a user by invoking profile lambda
// It performs synchronous request to the remote lambda. Panics in case of errors.
func (u UserProfileLambda) UserTokenSyncUnsafe(userID string) (ut models.UserToken) {
	return u.UserTokenAdvSyncUnsafe(userID, "")
}

// UserTokenAdvSyncUnsafe gets token associated with a user by invoking profile lambda
// It performs synchronous request to the remote lambda. Panics in case of errors.
func (u UserProfileLambda) UserTokenAdvSyncUnsafe(userID string, platformID models.PlatformID) (ut models.UserToken) {
	ut, err := u.UserTokenAdvSync(userID, platformID)
	core.ErrorHandler(err, u.Namespace, "Could not retrieve user token for "+userID)
	return
}

// UserToken gets token associated with a user by invoking profile lambda
// deprecated. Use UserProfileLambda.UserTokenSyncUnsafe
func UserToken(userID, profileLambda, region, namespace string) (models.UserToken, error) {
	l := UserProfileLambda{
		Region:            region,
		Namespace:         namespace,
		ProfileLambdaName: profileLambda,
	}
	return l.UserTokenSync(userID)
}

func ChatAttachment(title, text, fallback string, callbackId string,
	actions []ebm.AttachmentAction,
	fields []ebm.AttachmentField,
	timestamp int64) *ebm.Attachment {
	builder := eb.NewAttachmentBuilder().
		Title(title).
		Text(text).
		Fallback(fallback).
		Color(models.BlueColorHex).
		AttachmentType(models.DefaultAttachmentType).
		MarkDownIn([]ebm.MarkdownField{ebm.MarkdownFieldText, ebm.MarkdownFieldPretext}).
		CallbackId(callbackId).
		Footer(ebm.AttachmentFooter{Timestamp: timestamp})

	if len(actions) > 0 {
		builder.Actions(actions)
	}

	if len(fields) > 0 {
		builder.Fields(fields)
	}

	attach, _ := builder.Build()
	return attach
}

func AddChatEngagement(mc models.MessageCallback, title, text, fallback, userId string, actions []ebm.AttachmentAction,
	fields []ebm.AttachmentField, platformID models.PlatformID, urgent bool, table string,
	d *awsutils.DynamoRequest, namespace string, timestamp int64,
	check models.UserEngagementCheckWithValue) {
	eng := MakeUserEngagement(mc, title, text, fallback, userId,
		actions, fields, urgent, namespace, timestamp, check, platformID)
	AddEng(eng, table, d, namespace)
}

func urgency(urgent bool) models.PriorityValue {
	return core.IfThenElse(urgent, models.UrgentPriority, models.HighPriority).(models.PriorityValue)
}

// MakeUserEngagement creates a model of UserEngagement that will be posted
// to Dynamo table eventually
func MakeUserEngagement(mc models.MessageCallback,
	title, text, fallback, userID string,
	actions []ebm.AttachmentAction,
	fields []ebm.AttachmentField, urgent bool,
	namespace string, timestamp int64,
	check models.UserEngagementCheckWithValue,
	platformID models.PlatformID) models.UserEngagement {
	callbackID := mc.ToCallbackID()

	attach := ChatAttachment(title, text, fallback, callbackID, actions, fields, timestamp)
	engagement := eb.NewEngagementBuilder().
		Id(callbackID).
		WithResponseType(models.SlackInChannel).
		WithAttachment(attach).
		Build()
	bytes, err := engagement.ToJson()
	core.ErrorHandler(err, namespace, "Could not convert engagement to JSON")
	return models.UserEngagement{UserID: userID, TargetID: mc.Target, ID: callbackID,
		Script: string(bytes), Priority: urgency(urgent), CreatedAt: core.CurrentRFCTimestamp(),
		UserEngagementCheckWithValue: check, PlatformID: platformID}
}
