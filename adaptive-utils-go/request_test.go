package adaptive_utils_go

import (
	"encoding/json"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"gotest.tools/assert"
	"strings"
	"testing"
)

func TestRequestParsing(t *testing.T) {
	apiJSONBytes, err := readFile("testdata/member_left_channel.json")
	assert.NilError(t, err, "Could not read the file")
	var request events.APIGatewayProxyRequest
	err = json.Unmarshal(apiJSONBytes, &request)
	assert.NilError(t, err, "Could not parse to APIGatewayProxyRequest")

	requestPayload := request.Body//strings.Replace(request.Body, "payload=", "", -1)
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(requestPayload),
		slackevents.OptionNoVerifyToken(),
	)
	assert.NilError(t, err, "Could not parse to EventsAPIEvent")

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		eventType := eventsAPIEvent.InnerEvent.Type
		if eventType == "member_left_channel" {
			id := core.Uuid()
			callbackEvent := eventsAPIEvent.Data.(*slackevents.EventsAPICallbackEvent)
			apiAppID := callbackEvent.APIAppID

			np := models.NamespacePayload4{
				ID:        id,
				Namespace: "adaptive-channel",
				PlatformRequest: models.PlatformRequest{
					TeamID:       models.ParseTeamID(daosCommon.PlatformID(apiAppID)),
					SlackRequest: models.EventsAPIEvent(requestPayload),
				},
			}
			// Emulating publishing to SNS
			npSerBytes, err := json.Marshal(np)
			assert.NilError(t, err, "Could not marshal NamespacePayload")

			payload := string(npSerBytes)
			// Emulating reading from SNS
			var parsedNp models.NamespacePayload4
			parsedNp, err = models.UnmarshalNamespacePayload4JSON(payload)
			assert.NilError(t, err, "Payload=%s", payload)
			assert.Equal(t, np.ID, parsedNp.ID)
			if parsedNp.Namespace == "adaptive-channel" {
				evt := parsedNp.ToEventsAPIEventUnsafe()
				eventType := evt.InnerEvent.Type
				assert.Equal(t, "member_left_channel", eventType)
				slackMsg := *eventsAPIEvent.InnerEvent.Data.(*slack.MemberLeftChannelEvent)
				assert.Assert(t, slackMsg.User != "")
			} else {
				t.Fatal("Was expecting adaptive-channel namespace")
			}

		} else {
			t.Fatal("Was expecting member_joined_channel event")
		}

	} else {
		t.Fatal("Request is not event_callback")
	}
}

func TestApiRequestParsingX12LiteralBug(t *testing.T) {
	apiJsonBytes, err := readFile("testdata/api_request_x12_literal_bug.json")
	assert.NilError(t, err, "Could not read the file")
	var request events.APIGatewayProxyRequest
	err = json.Unmarshal(apiJsonBytes, &request)
	assert.NilError(t, err, "Could not parse to APIGatewayProxyRequest")

	requestPayload := strings.Replace(request.Body, "payload=", "", -1)
	_, err = slackevents.ParseEvent(
		json.RawMessage(requestPayload),
		slackevents.OptionNoVerifyToken(),
	)
	assert.NilError(t, err, "Could not parse to eventsAPIEvent")
}

func TestPayloadParsing(t *testing.T) {
	apiJSONBytes, err := readFile("testdata/api_request_x12_literal_bug.json")
	assert.NilError(t, err, "Could not read the file")
	var request events.APIGatewayProxyRequest
	err = json.Unmarshal(apiJSONBytes, &request)
	assert.NilError(t, err, "Could not parse to APIGatewayProxyRequest")

	slackRequest := models.ParseBodyAsSlackRequestUnsafe(request.Body)
	assert.Equal(t, slackRequest.Type, models.EventsAPIEventSlackRequestType)
	evt := slackRequest.ToEventsAPIEventUnsafe()
	assert.Equal(t, evt.TeamID, "TCUFJTB45")
}

func TestInteractionCallback(t *testing.T) {
	interactionCallback := "{\"type\":\"\",\"token\":\"\",\"callback_id\":\"\",\"response_url\":\"\",\"trigger_id\":\"\",\"action_ts\":\"\",\"team\":{\"id\":\"\",\"name\":\"\",\"domain\":\"\"},\"channel\":{\"id\":\"\",\"created\":0,\"is_open\":false,\"is_group\":false,\"is_shared\":false,\"is_im\":false,\"is_ext_shared\":false,\"is_org_shared\":false,\"is_pending_ext_shared\":false,\"is_private\":false,\"is_mpim\":false,\"unlinked\":0,\"name_normalized\":\"\",\"num_members\":0,\"priority\":0,\"user\":\"\",\"name\":\"\",\"creator\":\"\",\"is_archived\":false,\"members\":null,\"topic\":{\"value\":\"\",\"creator\":\"\",\"last_set\":0},\"purpose\":{\"value\":\"\",\"creator\":\"\",\"last_set\":0},\"is_channel\":false,\"is_general\":false,\"is_member\":false,\"locale\":\"\"},\"user\":{\"id\":\"\",\"team_id\":\"\",\"name\":\"\",\"deleted\":false,\"color\":\"\",\"real_name\":\"\",\"tz_label\":\"\",\"tz_offset\":0,\"profile\":{\"first_name\":\"\",\"last_name\":\"\",\"real_name\":\"\",\"real_name_normalized\":\"\",\"display_name\":\"\",\"display_name_normalized\":\"\",\"email\":\"\",\"skype\":\"\",\"phone\":\"\",\"image_24\":\"\",\"image_32\":\"\",\"image_48\":\"\",\"image_72\":\"\",\"image_192\":\"\",\"image_original\":\"\",\"title\":\"\",\"status_expiration\":0,\"team\":\"\",\"fields\":[]},\"is_bot\":false,\"is_admin\":false,\"is_owner\":false,\"is_primary_owner\":false,\"is_restricted\":false,\"is_ultra_restricted\":false,\"is_stranger\":false,\"is_app_user\":false,\"is_invited_user\":false,\"has_2fa\":false,\"has_files\":false,\"presence\":\"\",\"locale\":\"\",\"updated\":0,\"enterprise_user\":{\"id\":\"\",\"enterprise_id\":\"\",\"enterprise_name\":\"\",\"is_admin\":false,\"is_owner\":false,\"teams\":null}},\"original_message\":{\"replace_original\":false,\"delete_original\":false,\"blocks\":null},\"message\":{\"replace_original\":false,\"delete_original\":false,\"blocks\":null},\"name\":\"\",\"value\":\"\",\"message_ts\":\"\",\"attachment_id\":\"\",\"actions\":[],\"submission\":null}"
	var ic slack.InteractionCallback
	err := json.Unmarshal([]byte(interactionCallback), &ic)
	assert.NilError(t, err, "Could not parse to InteractionCallback")
	assert.Equal(t, slack.InteractionType(""), ic.Type)
}

func TestActionCallback(t *testing.T) {
	actionCallback := "[{\"block_id\":\"1\"}]"
	//interactionCallback := "{\"type\":\"\",\"token\":\"\",\"callback_id\":\"\",\"response_url\":\"\",\"trigger_id\":\"\",\"action_ts\":\"\",\"team\":{\"id\":\"\",\"name\":\"\",\"domain\":\"\"},\"channel\":{\"id\":\"\",\"created\":0,\"is_open\":false,\"is_group\":false,\"is_shared\":false,\"is_im\":false,\"is_ext_shared\":false,\"is_org_shared\":false,\"is_pending_ext_shared\":false,\"is_private\":false,\"is_mpim\":false,\"unlinked\":0,\"name_normalized\":\"\",\"num_members\":0,\"priority\":0,\"user\":\"\",\"name\":\"\",\"creator\":\"\",\"is_archived\":false,\"members\":null,\"topic\":{\"value\":\"\",\"creator\":\"\",\"last_set\":0},\"purpose\":{\"value\":\"\",\"creator\":\"\",\"last_set\":0},\"is_channel\":false,\"is_general\":false,\"is_member\":false,\"locale\":\"\"},\"user\":{\"id\":\"\",\"team_id\":\"\",\"name\":\"\",\"deleted\":false,\"color\":\"\",\"real_name\":\"\",\"tz_label\":\"\",\"tz_offset\":0,\"profile\":{\"first_name\":\"\",\"last_name\":\"\",\"real_name\":\"\",\"real_name_normalized\":\"\",\"display_name\":\"\",\"display_name_normalized\":\"\",\"email\":\"\",\"skype\":\"\",\"phone\":\"\",\"image_24\":\"\",\"image_32\":\"\",\"image_48\":\"\",\"image_72\":\"\",\"image_192\":\"\",\"image_original\":\"\",\"title\":\"\",\"status_expiration\":0,\"team\":\"\",\"fields\":[]},\"is_bot\":false,\"is_admin\":false,\"is_owner\":false,\"is_primary_owner\":false,\"is_restricted\":false,\"is_ultra_restricted\":false,\"is_stranger\":false,\"is_app_user\":false,\"is_invited_user\":false,\"has_2fa\":false,\"has_files\":false,\"presence\":\"\",\"locale\":\"\",\"updated\":0,\"enterprise_user\":{\"id\":\"\",\"enterprise_id\":\"\",\"enterprise_name\":\"\",\"is_admin\":false,\"is_owner\":false,\"teams\":null}},\"original_message\":{\"replace_original\":false,\"delete_original\":false,\"blocks\":null},\"message\":{\"replace_original\":false,\"delete_original\":false,\"blocks\":null},\"name\":\"\",\"value\":\"\",\"message_ts\":\"\",\"attachment_id\":\"\",\"actions\":{\"AttachmentActions\":null,\"BlockActions\":null},\"submission\":null}"
	var ac slack.ActionCallbacks
	err := ac.UnmarshalJSON([]byte(actionCallback))
	//err := json.Unmarshal([]byte(actionCallback), &ac)
	assert.NilError(t, err, "Could not parse to ActionCallbacks")
	assert.Equal(t, 1, len(ac.BlockActions))
	//ac.AttachmentActions[0].
}

func TestActionCallback3(t *testing.T) {
	action := slack.AttachmentAction{
		Text: "text",
	}
	ac := slack.ActionCallbacks{AttachmentActions: []*slack.AttachmentAction{&action}}
	bytes, err := ac.MarshalJSON()
	if err != nil {
		t.Errorf("Could not marshal ActionCallbacks: %v\n", err)
	}
	err = ac.UnmarshalJSON(bytes)
	actionCallback := "[{\"name\":\"\",\"text\":\"text\",\"type\":\"\"}]"
	assert.Equal(t, actionCallback, string(bytes))
	var ac2 slack.ActionCallbacks
	err = json.Unmarshal([]byte(actionCallback), &ac2)
	if err != nil {
		t.Errorf("Could not parse to ActionCallbacks: %v\n", err)
	}
	assert.Equal(t, 1, len(ac2.AttachmentActions))
	assert.Equal(t, "text", ac2.AttachmentActions[0].Text)
}

func TestActionCallback4(t *testing.T) {
	ac := slack.ActionCallbacks{AttachmentActions: []*slack.AttachmentAction{}}
	bytes, err := json.Marshal(ac)
//	ac.MarshalJSON()
	if err != nil {
		t.Errorf("Could not marshal ActionCallbacks: %v\n", err)
	}
	err = ac.UnmarshalJSON(bytes)
	actionCallback := "[]"
	assert.Equal(t, actionCallback, string(bytes))
	var ac2 slack.ActionCallbacks
	err = json.Unmarshal([]byte(actionCallback), &ac2)
	if err != nil {
		t.Errorf("Could not parse to ActionCallbacks: %v\n", err)
	}
	assert.Equal(t, 0, len(ac2.AttachmentActions))
}