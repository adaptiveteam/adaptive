package workflow_test

import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	"github.com/nlopes/slack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/sirupsen/logrus"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"

)

const FirstEvent wf.Event =  "first"

func FirstWorkflow_Workflow() wf.Template {
	return wf.Template{
		Init: "init", // initial state is "init". This is used when the user first triggers the workflow
		FSA: map[struct{wf.State; wf.Event}] wf.Handler {
			{State: "init", Event: ""}: wf.SimpleHandler(FirstWorkflow_OnInit, "next"),
			{State: "next", Event: "first"}: wf.SimpleHandler(FirstWorkflow_OnUserSelected(1), "done"),
			{State: "next", Event: "second"}: wf.SimpleHandler(FirstWorkflow_OnShowDialog, "show-dialog"),
			{State: "show-dialog", Event: "submit"}: wf.SimpleHandler(FirstWorkflow_OnDialogSubmitted, "done"),
			{State: "show-dialog", Event: "cancel"}: wf.SimpleHandler(FirstWorkflow_OnDialogCancelled, "done"),
		},
		Parser: wf.Parser,
	}
}

func FirstWorkflow_OnInit(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	out.Interaction = wf.MenuMessage("Select the next step", 
		wf.MenuOption("first", "first"),
		wf.MenuOption("second", "second"),
	)
	return
}

func FirstWorkflow_OnUserSelected(i int) wf.Handler {
	return func (ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		out.Interaction = wf.SimpleResponses(
			platform.Post(platform.ConversationID(ctx.Request.User.ID), 
				platform.MessageContent{Message: ui.Sprintf("You've selected %d", i)},
			),
		)
		return
	}
}

// FirstWorkflow_OnShowDialog returns a dialog to user
func FirstWorkflow_OnShowDialog(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	out.Interaction = wf.OpenSurvey(ebm.AttachmentActionSurvey{
		Title: "My dialog",
	})
	return
}

func FirstWorkflow_OnDialogSubmitted(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	out.Interaction = wf.SimpleResponses(
		platform.Post(platform.ConversationID(ctx.Request.User.ID), 
			platform.MessageContent{Message: "Dialog submitted"},
		),
	)
	return
}

func FirstWorkflow_OnDialogCancelled(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
	out.Interaction = wf.SimpleResponses(
		platform.Post(platform.ConversationID(ctx.Request.User.ID), 
			platform.MessageContent{Message: "Dialog cancelled"},
		),
	)
	return
}

func CreateInput(callback string) models.NamespacePayload4 {
	return  models.NamespacePayload4{
		PlatformRequest: models.PlatformRequest{
			SlackRequest: models.SlackRequest{
				InteractionCallback: slack.InteractionCallback{
					CallbackID: callback,
					User: slack.User{ID: "U1"},
				},
			},
		},
	}
}
type EnvMock struct {
	response platform.Response
	dialog ebm.AttachmentActionSurvey
}

func (e *EnvMock)PostAsync(response platform.Response) chan mapper.MessageID {
	return make(chan mapper.MessageID)
}
func (e *EnvMock)PostSync(response platform.Response) (id mapper.MessageID, err error) {
	if response.Type != platform.DeleteMessageByURLType && 
		response.Type != platform.DeleteTargetMessageType { 
		e.response = response
	}
	return mapper.MessageID {ConversationID: "convID", Ts: "ts"}, nil
}
func (e *EnvMock)PostSyncUnsafe(response platform.Response) (id mapper.MessageID) {
	id, _ = e.PostSync(response)
	return
}

func (e *EnvMock)ShowDialog(survey ebm.AttachmentActionSurvey2) error {
	e.dialog = survey.AttachmentActionSurvey
	return nil
}

var _ = Describe("Workflow", func() {
	
	Context("FirstWorkflow_Workflow", func(){
		template := FirstWorkflow_Workflow()
		var mock EnvMock
		prefix := models.ParsePath("/test/first-workflow")
		log := logger.LambdaLogger(logrus.InfoLevel)
		env := wf.Environment{
			Prefix: prefix,
			GetPlatformAPI: func (pid models.TeamID) mapper.PlatformAPI {
				return &mock
			},
			LogInfof: func(format string, args ...interface{}) { log.Infof(format, args ...)},
		}
		handler := template.GetRequestHandler(env)
		It("should handle an initial message", func(){
			np := CreateInput("")
			_, err := handler(wf.ActionPathFromCallbackID(np).ToRelActionPath(), np)
			Ω(err).Should(BeNil())
			Ω(mock.response.Type).Should(Equal(platform.PostToConversationType))
		})
		It("should handle next/first message", func(){
			np := CreateInput("/test/first-workflow?Event=first&State=next")
			_, err := handler(wf.ActionPathFromCallbackID(np).ToRelActionPath(), np)
			Ω(err).Should(BeNil())
			Ω(mock.response.Type).Should(Equal(platform.PostToConversationType))
		})
		It("should show dialog on second message", func(){
			np := CreateInput("/test/first-workflow?Event=second&State=next")
			_, err := handler(wf.ActionPathFromCallbackID(np).ToRelActionPath(), np)
			Ω(err).Should(BeNil())
			Ω(mock.dialog.Title).Should(Equal("My dialog"))
		})
	})
})
