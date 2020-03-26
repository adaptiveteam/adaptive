package common

import (
	"github.com/adaptiveteam/adaptive/workflows/exchange"
	// "github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/pkg/errors"
	// "strings"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	// userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	alog "github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
)
// WorkflowContext contains common state of workflows
type WorkflowContext struct {
	alog.AdaptiveLogger
	daosCommon.DynamoDBConnection
}

const NewAndOldIssuesKey = "NewAndOldIssues"
// GetNewAndOldIssues loads issue if necessary
func (w WorkflowContext)GetNewAndOldIssues(ctx wf.EventHandlingContext) (newAndOldIssues exchange.NewAndOldIssues, err error) {
	issueID := exchange.GetIssueID(ctx)
	itype := exchange.GetIssueType(ctx)
	if itype == "" {
		err = errors.Errorf("issueType is not defined in the context: CallbackID=%s, data is %v", ctx.Request.CallbackID, wf.ShowData(ctx.Data))
		return
	}
	log := w.AdaptiveLogger.
		WithField("issueID", issueID).
		WithField("IssueTypeFromContext", itype)
	log.Info("getNewAndOldIssues")
	newAndOldIssuesI, ok := ctx.RuntimeData[NewAndOldIssuesKey]
	if !ok {
		isShowingProgress := ctx.GetFlag(exchange.IsShowingProgressKey)
		log.Infof("GetNewAndOldIssues: runtime data is empty. Reading from database. isShowingProgress=%v", isShowingProgress)
		newAndOldIssues, err = issuesUtils.ReadNewAndOldIssuesAndPrefetch(itype, issueID, isShowingProgress)(w.DynamoDBConnection)
		if err != nil { 
			err = errors.Wrapf(err, "getNewAndOldIssues/w.IssueDAO.ReadNewAndOldIssuesAndPrefetch")
			return 
		}
		// NB: CANNOT modify input context! ctx.RuntimeData = runtimeData(newAndOldIssues)
	} else {
		newAndOldIssues = newAndOldIssuesI.(exchange.NewAndOldIssues)
	}
	err = errors.Wrap(err, "{getNewAndOldIssues}")
	return
}
