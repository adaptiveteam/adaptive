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
	issues "github.com/adaptiveteam/adaptive/workflows/issues"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// WorkflowInfo identifies a workflow.
type WorkflowInfo struct {
	Name string
	Init wf.State
}

// IssuesWorkflow is a description of an issues workflow
var IssuesWorkflow = WorkflowInfo{Name: issues.IssuesNamespace, Init: issues.InitState}

// var IssuesWorkflowImpl = issues.IssueWorkflow(d, clientID, logger)
// var IssuesWorkflow = IssuesWorkflowImpl.GetNamedTemplate()

const communityNamespace = "community"

var CommunityPath models.Path = models.ParsePath("/" + communityNamespace)

// var CreateIDOPath models.Path = CommunityPath.Append(CreateIDOWorkflow.Name)
var CreateIDOIssuePath models.Path = CommunityPath.Append(issues.IssuesNamespace)

var logger = alog.LambdaLogger(logrus.InfoLevel)

// EnterWorkflow sends the given event to the provided workflow.
// It modifies NamespacePayload4.CallbackID according to 
// the workflow namespace and event.
func EnterWorkflow(workflow WorkflowInfo, np models.NamespacePayload4, conn common.DynamoDBConnection, event wf.Event) error {
	initEvent := wf.ExternalActionPathWithData(CommunityPath.Append(workflow.Name), workflow.Init, event, map[string]string{}, false)
	logger.Infof("Starting workflow %s with path %s", workflow.Name, initEvent.Encode())
	np.InteractionCallback.CallbackID = initEvent.Encode()
	return InvokeWorkflow(np, conn)
}

// InvokeWorkflow sends `np` to a respective workflow.
// passes connection to workflow implementation.
func InvokeWorkflow(np models.NamespacePayload4, conn common.DynamoDBConnection) (err error) {
	if strings.HasPrefix(np.InteractionCallback.CallbackID, "/") {
		err = invokeWorkflowInner(np, conn,
			wf.TriggerImmediateEventForAnotherUser{
				UserID:     np.SlackRequest.InteractionCallback.User.ID,
				ActionPath: wf.ActionPathFromCallbackID(np),
			})
	} else {
		logger.Warnf("Unknown CallbackID %s", np.InteractionCallback.CallbackID)
		err = errors.New(fmt.Sprintf("Unknown CallbackID %s", np.InteractionCallback.CallbackID))
	}
	return
}

func invokeWorkflowInner(np models.NamespacePayload4, conn common.DynamoDBConnection, action wf.TriggerImmediateEventForAnotherUser) (err error) {
	np.SlackRequest.InteractionCallback.User.ID = action.UserID
	var furtherActions []wf.TriggerImmediateEventForAnotherUser

	furtherActions, err = communityRoutes(np, conn).Handler()(action.ActionPath.ToRelActionPath(), np)
	for _, a := range furtherActions {
		err = invokeWorkflowInner(np, conn, a)
		if err != nil {
			return
		}
	}
	return
}
// NB: It'll fail if there are issues in templates constructors.
// Though, it is unprobable
func workflows(conn common.DynamoDBConnection) []wf.NamedTemplate {
	IssuesWorkflowImpl := issues.IssueWorkflow(
		conn.Dynamo, conn.ClientID, conn.PlatformID, logger)
	return []wf.NamedTemplate{
		// CreateIDOWorkflow,
		IssuesWorkflowImpl.GetNamedTemplate(),
	}
}

// var allRoutes = communityRoutes()

func prepareEnvironmentWithoutPrefix(conn common.DynamoDBConnection) (env wf.Environment) {
	schema := models.SchemaForClientID(conn.ClientID)
	platformDAO := utilsPlatform.NewDAOFromSchema(conn.Dynamo, "workflows", schema)
	platformAdapter := mapper.SlackAdapter2(platformDAO)

	env = wf.ConstructEnvironmentWithoutPrefix(platformAdapter, wf.PostponeEventHandler(conn), logger)
	return
}


func communityRoutes(np models.NamespacePayload4, conn common.DynamoDBConnection) (routes wf.Routes) {
	workflowRoutes := wf.ToRoutingTable(CommunityPath, prepareEnvironmentWithoutPrefix(conn),
		workflows(conn))
	routes = map[string]wf.RequestHandler{
		communityNamespace: workflowRoutes.Handler(),
	}
	return
}

// func constructActionPath(prefix models.Path, state wf.State, event wf.Event) models.ActionPath {
// 	return wf.ExternalActionPath(prefix, state, event)
// }

// // SelectedIDOWorkflow allows to switch to a different workflow implementation
// var SelectedIDOWorkflow = IssuesWorkflow
// // var SelectedIDOWorkflow = CreateIDOWorkflow

// func onCreateIDONow1(np models.NamespacePayload4) error {
// 	return enterWorkflow(IssuesWorkflow, np, "")
// }

// func onViewIDOs(np models.NamespacePayload4) error {
// 	return enterWorkflow(SelectedIDOWorkflow, np, "view-idos")
// }
