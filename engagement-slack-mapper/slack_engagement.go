package engagement_slack_mapper

import (
	"github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/slack-go/slack"
)

type SlackEngagement struct {
	Message model.Message
}

func (s *SlackEngagement) MsgOptions() []slack.MsgOption {
	return slackMapper(s.Message)
}
