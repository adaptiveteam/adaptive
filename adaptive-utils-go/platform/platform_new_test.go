package platform_test

import (
	"encoding/json"
	"testing"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/stretchr/testify/assert"
)

func TestResponseParsing(t *testing.T) {
	r := platform.Post(
		//"userID", 
		platform.ConversationID("convID"),
		platform.MessageContent{
			Message: "hello",
		},
	)

	bytes, err := json.Marshal(r)
	str := string(bytes)
	assert.Equal(t, nil, err)
	assert.Equal(t, "{\"type\":\"post-to-channel\",\"post_to_conversation\":{\"conversation_id\":\"convID\",\"body\":{\"message\":\"hello\"}}}", str)
	var deserialized platform.Response
	json.Unmarshal(bytes, &deserialized)
	assert.Equal(t, r.PostToConversation.ConversationID, deserialized.PostToConversation.ConversationID)
}
