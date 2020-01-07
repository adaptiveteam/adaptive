package adaptive_utils_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func readFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func checkParsing(t *testing.T, filename, eventType string) {
	bytes, err := readFile(filename)
	assert.Nil(t, err, "Could not read the file")
	eav, err := ParseApiRequest(string(bytes))
	assert.Nil(t, err, "Could not parse to Callback")
	fmt.Println(eav.Type)
	assert.True(t, eav.Type == eventType)
}

func TestApiMsgParsing(t *testing.T) {
	// test callback message
	checkParsing(t, "testdata/callback.json", "event_callback")
	// test interaction message
	checkParsing(t, "testdata/interactive_message.json", "interactive_message")
	// test dialog submission
	checkParsing(t, "testdata/dialog_submission.json", "dialog_submission")
}

func TestSlackCallbackParsing(t *testing.T) {
	bytes, err := readFile("testdata/callback.json")
	assert.Nil(t, err, "Could not read the file")
	ar, _ := ParseApiRequest(string(bytes))
	eav := ParseAsCallbackMsg(ar)
	assert.True(t, eav.Type == "message")
}

func TestInteractionMsgParsing(t *testing.T) {
	bytes, err := readFile("testdata/interactive_message.json")
	assert.Nil(t, err, "Could not read the file")
	eav, err := ParseAsInteractionMsg(string(bytes))
	assert.Nil(t, err, "Could not parse to Callback")
	assert.True(t, eav.Type == "interactive_message")
}

func TestDialogSubmissionParsing(t *testing.T) {
	bytes, err := readFile("testdata/dialog_submission.json")
	assert.Nil(t, err, "Could not read the file")
	eav, err := ParseAsInteractionMsg(string(bytes))
	assert.Nil(t, err, "Could not parse to Callback")
	assert.True(t, eav.Type == "dialog_submission")
}
