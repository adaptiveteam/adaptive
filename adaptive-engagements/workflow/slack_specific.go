package workflow

import (
	// "github.com/adaptiveteam/adaptive/core-utils-go"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/daos/common"
	"encoding/json"

	utils "github.com/adaptiveteam/adaptive/adaptive-utils-go"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/slack-go/slack"
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

func GetPlatformAPIImpl(conn common.DynamoDBConnection) PlatformAPIForTeamID {
	return func (teamID models.TeamID) mapper.PlatformAPI {

		// token, err2 := platform.GetToken(teamID)(conn)
		// core_utils_go.ErrorHandler(err2, "GetPlatformAPIImpl", "GetPlatformAPIImpl")

		return mapper.SlackAdapterForTeamID(conn)
	}
}
func ConstructEnvironmentWithoutPrefix(conn common.DynamoDBConnection, postpone PostponeEvent, log alog.AdaptiveLogger,
	resolveCommunity ResolveCommunity) Environment {
	return Environment{
		GetPlatformAPI: GetPlatformAPIImpl(conn),
		LogInfof:       log.Infof,
		PostponeEvent:  postpone,
		ResolveCommunity: resolveCommunity,
	}
}
