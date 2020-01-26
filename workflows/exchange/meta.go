package exchange

import (
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
)

// WorkflowInfo identifies a workflow.
type WorkflowInfo struct {
	Name string
	Init wf.State
}
