package adaptive_utils_go

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMessageCallbackParsing(t *testing.T) {
	mc := &models.MessageCallback{
		Module: "coaching",
		Source: "source",
		Topic:  "feedback",
		Action: "ask",
		Target: "target",
		Month:  "1",
		Year:   "2019",
	}
	assert.True(t, len(mc.ToCallbackID()) > 0)
	assert.True(t, strings.Contains(mc.ToCallbackID(), "ask"))

	// set action to a different value
	mc.Set("Action", "collect")
	assert.True(t, strings.Contains(mc.ToCallbackID(), "collect"))

	_, err := ParseToCallback(mc.ToCallbackID())
	assert.True(t, err == nil)
}

func TestMessageInvalidCallbackParsing(t *testing.T) {
	_, err := ParseToCallback("1:2:3:4:5:6")
	assert.True(t, err != nil)
}
