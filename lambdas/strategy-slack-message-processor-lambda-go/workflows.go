package lambda

import (
	// "github.com/sirupsen/logrus"
	// "github.com/pkg/errors"
	// "fmt"
	// alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	// "strings"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)
var Namespace = "strategy"
var StrategyPath models.Path = models.ParsePath("/" + Namespace)
var CreateObjectivePath models.Path = StrategyPath.Append(CreateObjectiveWorkflow.Name)

// NB: It'll fail if there are issues in templates constructors.
// Though, it is unprobable
var workflows = []wf.NamedTemplate{CreateObjectiveWorkflow}
// var allRoutes = strategyRoutes()

// func prepareEnvironmentWithoutPrefix() (env wf.Environment) {
// 	log := alog.LambdaLogger(logrus.InfoLevel)
// 	env = wf.ConstructEnvironmentWithoutPrefix(platformAdapter, log)
// 	return
// }

// func invokeWorkflow(np models.NamespacePayload4) (err error) {
// 	if strings.HasPrefix(np.InteractionCallback.CallbackID, "/") {
// 		err = allRoutes.Handler()(wf.ActionPathFromCallbackID(np), np)
// 	} else {
// 		logger.Warnf("Unknown CallbackID %s", np.InteractionCallback.CallbackID)
// 		err = errors.New(fmt.Sprintf("Unknown CallbackID %s", np.InteractionCallback.CallbackID))
// 	}
// 	return
// }

// func strategyRoutes() (routes wf.Routes) {
// 	workflowRoutes := wf.ToRoutingTable(StrategyPath, prepareEnvironmentWithoutPrefix(), workflows)
// 	routes = map[string]wf.RequestHandler {
// 		Namespace: workflowRoutes.Handler(),
// 	}
// 	return 
// }

// func enterWorkflow(workflow wf.NamedTemplate, np models.NamespacePayload4, event wf.Event) error {
// 	initEvent := wf.ExternalActionPath(StrategyPath.Append(workflow.Name), workflow.Template.Init, event)
// 	logger.Infof("Starting workflow %s with path %s", workflow.Name, initEvent.Encode())
// 	np.InteractionCallback.CallbackID = initEvent.Encode()
// 	return invokeWorkflow(np)
// }

// func initWorkflow(workflow wf.NamedTemplate, np models.NamespacePayload4) error {
// 	return enterWorkflow(workflow, np, "")
// }

// func onCreateObjective(np models.NamespacePayload4) error {
// 	return initWorkflow(CreateObjectiveWorkflow, np)
// }

// func onViewStrategyObjectives(np models.NamespacePayload4) error {
// 	return enterWorkflow(CreateObjectiveWorkflow, np, ViewObjectivesEvent)
// }
