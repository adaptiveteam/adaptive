package fetch_dialog

import (
	"fmt"
)

const (
	knownDialogID = "e53bf79b-dff2-4a34-9c92-7f3bf07cdea5"
	knownDialogContext = "dialog/collaboration-report/report-contents/analysis/energy"
	mockPackageName = "collaboration_report"
	mockAlias = "energy"
	mockSubject = "above"
)

var (
	expectedDialogEntry = DialogEntry{
		Context:ConvertContextPathToHash(knownDialogContext),
		Subject:mockSubject,
		Updated:"2019-5-5",
		Dialog: []string{"dialog one", "dialog two", "dialog three"},
		Comments: []string{"comment one", "comment two", "comment three"},
		DialogID: knownDialogID,
		LearnMoreLink:"https://learn-more-link.com",
		LearnMoreContent:"test learn more content",
		BuildBranch:"build",
		CultivationBranch:"cultivate",
		MasterBranch:"master",
		BuildID:"test-uuid-for-build-id",
	}
	expectedAlias = ContextAliasEntry{
		Context: knownDialogContext,
		Alias: mockPackageName + "#" + mockAlias,
		BuildID: knownDialogID,
	}
)

func addMockData(dao DAO) error {
	err := dao.Create(expectedDialogEntry)
	if err != nil {
		fmt.Printf("Couldn't create dialog entry %v\n", err)
	} else {
		err = dao.CreateAlias(expectedAlias)
		if err != nil {
			fmt.Printf("Couldn't create dialog entry %v\n", err)
		}
	}
	return err
}
