package workflows

import (
	"fmt"
	"strings"

	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	utilsPlatform "github.com/adaptiveteam/adaptive/adaptive-utils-go/platform"
	common "github.com/adaptiveteam/adaptive/daos/common"
	mapper "github.com/adaptiveteam/adaptive/engagement-slack-mapper"
	"github.com/adaptiveteam/adaptive/workflows/issues"
	"github.com/adaptiveteam/adaptive/workflows/request_coach"
	"github.com/adaptiveteam/adaptive/workflows/coachees"
	"github.com/adaptiveteam/adaptive/workflows/closeout"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/adaptiveteam/adaptive/workflows/exchange"
)

// WorkflowInfo -
type WorkflowInfo = exchange.WorkflowInfo
// IssuesWorkflow is a description of an issues workflow
var IssuesWorkflow = issues.IssuesWorkflow
// RequestCoachWorkflow -
var RequestCoachWorkflow = request_coach.RequestCoachWorkflow
// RequestCloseoutWorkflow - 
var RequestCloseoutWorkflow = closeout.RequestCloseoutWorkflow
// ViewCoacheesWorkflow -
var ViewCoacheesWorkflow = coachees.ViewCoacheesWorkflow
// var IssuesWorkflowImpl = issues.IssueWorkflow(d, clientID, logger)
// var IssuesWorkflow = IssuesWorkflowImpl.GetNamedTemplate()

const communityNamespace = exchange.CommunityNamespace

var CommunityPath models.Path = exchange.CommunityPath

const FeedbackNamespace = exchange.FeedbackNamespace

var CoachingPath models.Path = exchange.CoachingPath
// var ViewCoacheeIDOsPath models.Path = CoachingPath.Append(ViewCoacheeIDOs.Name)

var logger = alog.LambdaLogger(logrus.InfoLevel)

// EnterWorkflow sends the given event to the provided workflow.
// It modifies NamespacePayload4.CallbackID according to 
// the workflow namespace and event.
func EnterWorkflow(workflow WorkflowInfo, np models.NamespacePayload4, conn common.DynamoDBConnection, event wf.Event) error {
	if conn.ClientID == "" {
		return errors.New("EnterWorkflow: clientID == ''")
	}
	initEvent := wf.ExternalActionPathWithData(workflow.Prefix.Append(workflow.Name), workflow.Init, event, map[string]string{}, false)
	logger.Infof("Starting workflow %s with path %s", workflow.Name, initEvent.Encode())
	np.InteractionCallback.CallbackID = initEvent.Encode()
	return InvokeWorkflow(np, conn)
}

// InvokeWorkflow sends `np` to a respective workflow.
// passes connection to workflow implementation.
func InvokeWorkflow(np models.NamespacePayload4, conn common.DynamoDBConnection) (err error) {
	if strings.HasPrefix(np.InteractionCallback.CallbackID, "/") {
		err = invokeWorkflowInner(np,
			wf.TriggerImmediateEventForAnotherUser{
				UserID:     np.SlackRequest.InteractionCallback.User.ID,
				ActionPath: wf.ActionPathFromCallbackID(np),
			})(conn)
	} else {
		logger.Warnf("Unknown CallbackID %s", np.InteractionCallback.CallbackID)
		err = errors.New(fmt.Sprintf("Unknown CallbackID %s", np.InteractionCallback.CallbackID))
	}
	return
}

func invokeWorkflowInner(np models.NamespacePayload4, action wf.TriggerImmediateEventForAnotherUser) func (conn common.DynamoDBConnection) (err error) {
	return func (conn common.DynamoDBConnection) (err error) {
		np.SlackRequest.InteractionCallback.User.ID = action.UserID
		np.PlatformID = conn.PlatformID
		np.InteractionCallback.CallbackID = action.ActionPath.Encode()
		var furtherActions []wf.TriggerImmediateEventForAnotherUser
		logger.
			WithField("userID", action.UserID).
			WithField("action.ActionPath", np.InteractionCallback.CallbackID).
			Info("invokeWorkflowInner")
		furtherActions, err = communityRoutes(np, conn).Handler()(action.ActionPath.ToRelActionPath(), np)
		for _, a := range furtherActions {
			err = invokeWorkflowInner(np, a)(conn)
			if err != nil {
				return
			}
		}
		return
	}
}
// NB: It'll fail if there are issues in templates constructors.
// Though, it is unprobable
func communityWorkflows(conn common.DynamoDBConnection) (templates []wf.NamedTemplate) {
	IssuesWorkflowImpl := issues.CreateIssueWorkflow(conn, logger)
	RequestCoachWorkflowImpl := request_coach.CreateRequestCoachWorkflow(conn, logger)
	RequestCloseoutWorkflowImpl := closeout.CreateRequestCloseoutWorkflow(conn, logger)
	templates = []wf.NamedTemplate{
		IssuesWorkflowImpl.GetNamedTemplate(),
		RequestCoachWorkflowImpl.GetNamedTemplate(),
		RequestCloseoutWorkflowImpl.GetNamedTemplate(),
	}
	for _, t := range templates {
		logger.Infof("Community Workflow template: %s", t.Name)
	}
	return
}

func feedbackWorkflows(conn common.DynamoDBConnection) (templates []wf.NamedTemplate) {
	CoacheesWorkflowImpl := coachees.CreateCoacheesWorkflow(conn, logger)
	templates = []wf.NamedTemplate{
		CoacheesWorkflowImpl.GetNamedTemplate(),
	}
	for _, t := range templates {
		logger.Infof("Feedback Workflow template: %s", t.Name)
	}
	return
}

func prepareEnvironmentWithoutPrefix(conn common.DynamoDBConnection) (env wf.Environment) {
	schema := models.SchemaForClientID(conn.ClientID)
	platformDAO := utilsPlatform.NewDAOFromSchema(conn.Dynamo, "workflows", schema)
	platformAdapter := mapper.SlackAdapter2(platformDAO)

	env = wf.ConstructEnvironmentWithoutPrefix(platformAdapter, 
		wf.PostponeEventHandler(conn), logger)
	return
}


func communityRoutes(np models.NamespacePayload4, conn common.DynamoDBConnection) (routes wf.Routes) {
	communityRoutes := wf.ToRoutingTable(CommunityPath, 
		prepareEnvironmentWithoutPrefix(conn),
		communityWorkflows(conn))
	feedbackRoutes := wf.ToRoutingTable(CoachingPath,
		prepareEnvironmentWithoutPrefix(conn),
		feedbackWorkflows(conn))
	routes = map[string]wf.RequestHandler{
		communityNamespace: communityRoutes.Handler(),
		FeedbackNamespace: feedbackRoutes.Handler(),
	}
	return
}

// InvokeWorkflowByPath is 
func InvokeWorkflowByPath(immediateAction wf.TriggerImmediateEventForAnotherUser) func (conn common.DynamoDBConnection) (err error) {
	np := models.NamespacePayload4{}
	return invokeWorkflowInner(np, immediateAction)
}
