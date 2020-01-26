package exchange

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

// WorkflowInfo identifies a workflow.
type WorkflowInfo struct {
	Prefix models.Path
	Name   string
	Init   wf.State
}
