package workflow

import (
	"encoding/json"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// APIShowDialog wraps Slack API and returns a function that could be used as showDialog in HandleRequest
func APIShowDialog(api *slack.Client) func(slack.InteractionCallback, ebm.AttachmentActionSurvey, Instance, string) (err error) {
	return func(request slack.InteractionCallback, survey ebm.AttachmentActionSurvey, instance Instance, callbackID string) (err error) {
		survBytes, err := json.Marshal(survey)
		if err == nil {
			logrus.Infof("APIShowDialog callbackID=%s: %v", callbackID, string(survBytes))
			instanceBytes, err := json.Marshal(instance)
			if err == nil {
				logrus.Infof("Dialog %v", string(instanceBytes))
				survState := func() string { return string(instanceBytes) }
				err = utils.SlackSurvey(api, request, survey, callbackID, survState)
			}
		}
		return
	}
}

// SelectedValue extracts the selected value when user clicks on a line in drop down.
func SelectedValue(request slack.InteractionCallback) (value string, err error) {
	actions := request.ActionCallback.AttachmentActions
	if len(actions) > 0 && actions[0].SelectedOptions != nil && len(actions[0].SelectedOptions) > 0 {
		value = actions[0].SelectedOptions[0].Value
	} else {
		err = errors.New("No selection in request")
	}
	return
}

func ConstructEnvironmentWithoutPrefix(adapter mapper.PlatformAdapter2, postpone PostponeEvent, log alog.AdaptiveLogger) Environment {
	return Environment{
		GetPlatformAPI: adapter.ForPlatformID,
		LogInfof:       log.Infof,
		PostponeEvent:  postpone,
	}
}
