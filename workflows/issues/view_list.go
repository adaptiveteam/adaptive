package issues

import (
	"log"
	issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	ui "github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type IssuePredicateFactory = func(wf.EventHandlingContext) IssuePredicate

func filterUserAdvocate(ctx wf.EventHandlingContext) IssuePredicate {
	return func(issue Issue) bool { return issue.StrategyObjective.Advocate == ctx.Request.User.ID }
}
func unfiltered(wf.EventHandlingContext) IssuePredicate { return func(Issue) bool { return true } }

// IssueQuery loads issues from DB
type IssueQuery = func(conn DynamoDBConnection) (issues []Issue, err error)
type IssueQueryFactory = func(wf.EventHandlingContext) IssueQuery

// configures context as if it represented a different issue.
func setIssue(ctx *wf.EventHandlingContext, issue Issue) {
	newAndOldIssues := NewAndOldIssues{
		NewIssue: issue,
		OldIssue: issue,
		Updated:  false,
	}
	ctx.Data[issueIDKey] = issue.UserObjective.ID
	if issue.UserObjective.ID == "" {
		log.Printf("INVALID(1): issueID is empty %v\n", issue)
	}
	ctx.Data[issueTypeKey] = string(issue.GetIssueType())
	ctx.RuntimeData = runtimeData(newAndOldIssues)
}

func (w workflowImpl) queryFromPredicate(issueType IssueType, issueFilterFactory IssuePredicateFactory) IssueQueryFactory {
	return func(ctx wf.EventHandlingContext) IssueQuery {
		issueFilter := issueFilterFactory(ctx)
		userID := ctx.Request.User.ID
		// issueType := getIssueTypeFromContext(ctx)
		isCompleted := 0 // only not finished

		return func(conn DynamoDBConnection) (issues []Issue, err error) {
			issues, err = issuesUtils.SelectFromIssuesWhereTypeAndUserID(userID, issueType, isCompleted)(conn)
			for _, issue := range issues {
				w.AdaptiveLogger.
					WithField("issue.UserObjective.ID", issue.UserObjective.ID).
					WithField("issue.StrategyObjective.ID", issue.StrategyObjective.ID).
					Infof("queryFromPredicate/for issues")
				if issue.UserObjective.ID == "" || (
					issueType == SObjective && issue.StrategyObjective.ID == "") || (
					issueType == Initiative && issue.StrategyInitiative.ID == "") {
					w.AdaptiveLogger.Warnf("INVALID(4): Issue ID is incorrect: %v", issue)
				}
			}
			if err != nil {
				return
			}

			issues = filterIssues(issues, issueFilter)
			return
		}
	}
}
func (w workflowImpl) OnViewListOfIssues(issueType IssueType, issueFilterFactory IssuePredicateFactory) wf.Handler {
	qf := w.queryFromPredicate(issueType, issueFilterFactory)
	tc := getTypeClass(issueType)
	return w.OnViewListOfQueryIssues(issueType, qf, tc.IssueTypeName()+"s")
}

// OnViewListOfQueryIssues shows the list of elements returned by the query.
// 
// queryItemPluralTitle - the user-facing name of the list element in plural form.
// It's used in the title of the list.
func (w workflowImpl) OnViewListOfQueryIssues(issueType IssueType, 
	issueQueryFactory IssueQueryFactory, 
	queryItemPluralTitle ui.PlainText,
) wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		ctx.Data[issueTypeKey] = string(issueType)
		return w.OnViewListOfQueryIssuesWithTypeInContext(issueQueryFactory, 
			func(_ IssueType) ui.PlainText {return queryItemPluralTitle})(ctx)
	}
}

func (w workflowImpl) OnViewListOfQueryIssuesWithTypeInContext(issueQueryFactory IssueQueryFactory, 
	queryItemPluralTitle func (IssueType) ui.PlainText) wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		issueType := getIssueTypeFromContext(ctx)
		var issues []Issue
		issues, err = issueQueryFactory(ctx)(w.DynamoDBConnection)
		if err != nil {
			return
		}

		var prefetchedIssues []Issue
		prefetchedIssues, err = w.prefetchManyIssuesWithoutProgress(ctx.PlatformID, issues)
		if err != nil {
			return
		}

		threadMessages := wf.InteractiveMessages()
		for _, issue := range prefetchedIssues {
			setIssue(&ctx, issue) // We reuse `standardView` that shows the current issue
			out, err = w.standardView(ctx)
			if err != nil {
				return
			}
			threadMessages = append(threadMessages, out.Messages...)
		}
		var msg ui.RichText
		if len(threadMessages) == 0 {
			msg = ui.Sprintf("There are no %s yet.", queryItemPluralTitle(issueType))
		} else {
			msg = ui.Sprintf("You can find the list of %s in the below thread. :point_down:", queryItemPluralTitle(issueType))
		}
		out.Interaction.Messages = wf.InteractiveMessages(wf.InteractiveMessage{
			PassiveMessage: wf.PassiveMessage{
				Text: msg,
				OverrideOriginal: true,
				IsPermanentMessage: true, // we don't want the thread title to disappear ever
			},
			Thread: threadMessages,
		})
		// out.Interaction.KeepOriginal = false
		return
	}
}

// OnPromptStaleIssues shows a prompt if there are stale issues
func (w workflowImpl) OnPromptStaleIssues(issueType IssueType) wf.Handler {
	return func(ctx wf.EventHandlingContext) (out wf.EventOutput, err error) {
		log := w.AdaptiveLogger.
			WithField("issueType", issueType).
			WithField("Handler", "OnPromptStaleIssues")

		log.Infof("start")
		// TODO: check in database if we indeed have stale issues
		out = ctx.Prompt(
			ui.Sprintf("You have %s(s) that haven't been updated in last 7 days. Would you like to update them?", issueType.Template()),
			wf.Button(ConfirmEvent, "Yes"),
			wf.Button(DismissEvent, "Skip this, please"),
		)
		out.DataOverride = map[string]string{issueTypeKey: string(issueType)}
		log.Infof("out=%v", out)
		return
	}
}

