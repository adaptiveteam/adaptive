package issues

import (
	// "github.com/adaptiveteam/adaptive/workflows/exchange"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/pkg/errors"
	// "strings"
	wf "github.com/adaptiveteam/adaptive/adaptive-engagements/workflow"
	userObjective "github.com/adaptiveteam/adaptive/daos/userObjective"
	// issuesUtils "github.com/adaptiveteam/adaptive/adaptive-utils-go/issues"
)

// getNewAndOldIssues loads issue if necessary
func (w workflowImpl)getNewAndOldIssues(ctx wf.EventHandlingContext) (newAndOldIssues NewAndOldIssues, err error) {
	return w.WorkflowContext.GetNewAndOldIssues(ctx)
}

// prefetch reads joined tables and puts related data into issue
func (w workflowImpl)prefetch(ctx wf.EventHandlingContext, 
	issue *Issue,
	) (err error ) {
	isShowingProgress := ctx.GetFlag(isShowingProgressKey)
	if isShowingProgress {
		// 	objectiveProgress := LatestProgressUpdateByObjectiveID(issue.UserObjective.ID)
		issue.PrefetchedData.Progress, err = IssueProgressReadAll(issue.UserObjective.ID, 0)(w.DynamoDBConnection)
		if err != nil { return }
	}
	return w.prefetchIssueWithoutProgress(ctx.PlatformID, issue)
}

// prefetchIssueWithoutProgress loads issue information ignoring context
func (w workflowImpl)prefetchIssueWithoutProgress(
	platformID models.PlatformID,
	issue *Issue,
	) (err error ) {
	
	w.AdaptiveLogger.
		WithField("issue.UserObjective.ID",issue.UserObjective.ID).
		WithField("issue.StrategyObjective.ID",issue.StrategyObjective.ID).
		Info("prefetchIssueWithoutProgress")
	if issue.UserObjective.AccountabilityPartner != "none" && 
		issue.UserObjective.AccountabilityPartner != "requested" && 
		issue.UserObjective.AccountabilityPartner != "" {
		var partners [] models.User
		partners, err = UserRead(issue.UserObjective.AccountabilityPartner)(w.DynamoDBConnection)
		for _, p := range partners {
			issue.PrefetchedData.AccountabilityPartner = p
		}
		if err != nil { return }
	}

	switch issue.StrategyAlignmentEntityType {
	case userObjective.ObjectiveStrategyObjectiveAlignment:
		issue.PrefetchedData.AlignedCapabilityObjective, err = StrategyObjectiveRead(issue.StrategyAlignmentEntityID)(w.DynamoDBConnection)
	case userObjective.ObjectiveStrategyInitiativeAlignment:
		issue.PrefetchedData.AlignedCapabilityInitiative, err = StrategyInitiativeRead(issue.StrategyAlignmentEntityID)(w.DynamoDBConnection)
	case userObjective.ObjectiveCompetencyAlignment:
		issue.PrefetchedData.AlignedCompetency, err = CompetencyRead(issue.StrategyAlignmentEntityID)(w.DynamoDBConnection)
	}
	if err != nil {
		w.AdaptiveLogger.
			WithError(err).
			WithField("issue.StrategyAlignmentEntityType", issue.StrategyAlignmentEntityType).
			WithField("issue.StrategyAlignmentEntityID", issue.StrategyAlignmentEntityID).
			Infof("prefetchIssueWithoutProgress, couldn't load issue.PrefetchedData.Aligned*")
		err = nil
	}

	itype := issue.GetIssueType()
	switch itype {
	case IDO:
		// see above - prefetched data
	case SObjective:
		// already prefetched?
		if len(issue.StrategyObjective.CapabilityCommunityIDs) > 0 {
			capCommID := issue.StrategyObjective.CapabilityCommunityIDs[0]
			w.AdaptiveLogger.
				WithField("capCommID", capCommID).
				Infof("prefetchIssueWithoutProgress, prefetching AlignedCapabilityCommunity by CapabilityCommunityIDs[0]")
			issue.PrefetchedData.AlignedCapabilityCommunity, err = CapabilityCommunityRead(capCommID)(w.DynamoDBConnection)
		}
		// splits := strings.Split(issue.UserObjective.ID, "_")
		// if len(splits) == 2 {
		// 	soID := splits[0]
		// 	capCommID := splits[1]
		// 	issue.PrefetchedData.AlignedCapabilityObjective, err = w.StrategyObjectiveDAO.Read(platformID, soID)
		// 	if err != nil { return }
		// 	issue.PrefetchedData.AlignedCapabilityCommunity, err = w.CapabilityCommunityDAO.Read(platformID, capCommID)
		// } else {
		// 	issue.PrefetchedData.AlignedCapabilityObjective, err = w.StrategyObjectiveDAO.Read(platformID, issue.UserObjective.ID)
		// }
	case Initiative:
		initCommID := issue.StrategyInitiative.InitiativeCommunityID
		if initCommID != "" {
			issue.PrefetchedData.AlignedInitiativeCommunity, err =
				StrategyInitiativeCommunityRead(initCommID)(w.DynamoDBConnection)
			if err != nil { return }
		}
		capObjID := issue.StrategyInitiative.CapabilityObjective
		if capObjID != "" {
			issue.PrefetchedData.AlignedCapabilityObjective, err = StrategyObjectiveRead(capObjID)(w.DynamoDBConnection)
		}
	default:
		w.AdaptiveLogger.WithField("issue.ID", issue.UserObjective.ID).Info("Not aligned with any strategy")
	}
	err = errors.Wrap(err, "{prefetch}")
	return
}

func (w workflowImpl)prefetchManyIssuesWithoutProgress(
	platformID models.PlatformID,
	issues []Issue,
)(prefetchedIssues []Issue, err error ) {
	for _, issue := range issues {
		err = w.prefetchIssueWithoutProgress(platformID, &issue)
		if err != nil {
			return
		} 
		prefetchedIssues = append(prefetchedIssues, issue)
	}
	return
}
