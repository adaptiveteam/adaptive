package issues

import (
	"github.com/adaptiveteam/adaptive/workflows/common"
	"time"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

// IssueTypeClass contains various helper methods for converting issue to/from forms and dialogs
type IssueTypeClass interface {
	IsCapabilityCommunityNeeded() bool
	IsInitiativeCommunityNeeded() bool
	// Empty returns a new issue of this type
	Empty() Issue
	CreateDialog(w workflowImpl, ctx wf.EventHandlingContext, issue Issue) (survey ebm.AttachmentActionSurvey, err error)
	// ExtractFromContext reads id from data. If it's not empty then we are updating
	// an existing issue, otherwise creating a new one.
	ExtractFromContext(ctx wf.EventHandlingContext, id string, updated bool, oldIssue Issue) (newIssue Issue)
	// View that represents issue information. If the issue is changed,
	// the view will contain the difference.
	// Progress list might be empty.
	// View(w workflowImpl, isShowingDetails, isShowingProgress bool,
	// 	newAndOldIssues NewAndOldIssues,
	// ) (fields []ebm.AttachmentField)
	// ObjectiveToFields - TODO rename to ShortViewForStrategy
	// ObjectiveToFields(w workflowImpl, newAndOldIssues NewAndOldIssues) (fields []ebm.AttachmentField)
	// TODO: check other contexts.
	CloseoutDisagreementContext() string
	IssueTypeName() ui.PlainText
}

type IDOImpl struct{}
type SObjectiveImpl struct{}
type InitiativeImpl struct{}

func getTypeClass(itype IssueType) (tc IssueTypeClass) {
	switch itype {
	case IDO:
		tc = IDOImpl{}
	case SObjective:
		tc = SObjectiveImpl{}
	case Initiative:
		tc = InitiativeImpl{}
	}
	return
}

func (IDOImpl) IsCapabilityCommunityNeeded() bool        { return false }
func (SObjectiveImpl) IsCapabilityCommunityNeeded() bool { return true }
func (InitiativeImpl) IsCapabilityCommunityNeeded() bool { return false }

func (IDOImpl) IsInitiativeCommunityNeeded() bool        { return false }
func (SObjectiveImpl) IsInitiativeCommunityNeeded() bool { return false }
func (InitiativeImpl) IsInitiativeCommunityNeeded() bool { return true }

func (IDOImpl) Empty() (issue Issue) {
	issue.UserObjective.ObjectiveType = userObjective.IndividualDevelopmentObjective
	issue.UserObjective.CreatedAt = core.ISODateLayout.Format(time.Now())
	issue.UserObjective.CreatedDate = core.USDateLayout.Format(time.Now())
	return
}
func (SObjectiveImpl) Empty() (issue Issue) {
	issue.UserObjective.ObjectiveType = userObjective.StrategyDevelopmentObjective
	issue.UserObjective.CreatedAt = core.ISODateLayout.Format(time.Now())
	issue.UserObjective.CreatedDate = core.USDateLayout.Format(time.Now())
	issue.StrategyAlignmentEntityType = userObjective.ObjectiveStrategyObjectiveAlignment
	return
}
func (InitiativeImpl) Empty() (issue Issue) {
	issue.UserObjective.ObjectiveType = userObjective.StrategyDevelopmentObjective
	issue.UserObjective.CreatedAt = core.ISODateLayout.Format(time.Now())
	issue.UserObjective.CreatedDate = core.USDateLayout.Format(time.Now())
	issue.StrategyAlignmentEntityType = userObjective.ObjectiveStrategyInitiativeAlignment
	return
}
func (IDOImpl) CloseoutDisagreementContext() string {
	return common.GetDialogContext(CloseoutDisagreementContext, IDO)
}
func (SObjectiveImpl) CloseoutDisagreementContext() string {
	return common.GetDialogContext(CloseoutDisagreementContext, SObjective)
}
func (InitiativeImpl) CloseoutDisagreementContext() string {
	return common.GetDialogContext(CloseoutDisagreementContext, Initiative)
}

func (IDOImpl) IssueTypeName() ui.PlainText        { return "Individual Development Objective" }
func (SObjectiveImpl) IssueTypeName() ui.PlainText { return "Strategy Objective" }
func (InitiativeImpl) IssueTypeName() ui.PlainText { return "Initiative" }
