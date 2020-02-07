package fetch_dialog

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	//"os"
	"reflect"
	"testing"
)

func CheckFetchBySubject(t *testing.T, dialogTableName string) {

	dao := localStackDao()
	resultSubject,err := dao.FetchByContextSubject(
		"dialog/collaboration-report/report-contents/analysis/energy",
		"above")
	if err != nil {t.Errorf("Error in NewFetchDialogByContext.FetchDialog (%v)", err)}

	resultDialogID, _, err := dao.FetchByDialogID("e53bf79b-dff2-4a34-9c92-7f3bf07cdea5")
	
	if err != nil {t.Errorf("Error in NewFetchDialogByID.FetchDialog (%v)", err)}

	resultAlias,err := dao.FetchByAlias(
		"collaboration_report",
		"energy",
		"above")

	if err != nil {t.Errorf("Error in NewFetchDialogByAlias.FetchDialog (%v)", err)}

	if !reflect.DeepEqual(resultSubject, resultDialogID) {
		t.Errorf("Expected resultSubject, %v\n resultDialogID, %v, to be the same.", resultSubject, resultDialogID)
	}

	if !reflect.DeepEqual(resultSubject, resultAlias) {
		t.Errorf("Expected resultSubject, %v\n resultAlias, %v, to be the same.", resultSubject, resultAlias)
	}
}

func CheckDialogMocking(t *testing.T) {
	context := "/test/test/test"
	subject := "test-subject"
	mockResult := DialogEntry{
		Context:context,
		Subject:subject,
		Updated:"2019-5-5",
		Dialog: []string{"dialog one", "dialog two", "dialog three"},
		Comments: []string{"comment one", "comment two", "comment three"},
		DialogID:"test-uuid-for-dialog-id",
		LearnMoreLink:"https://learn-more-link.com",
		LearnMoreContent:"test learn more content",
		BuildBranch:"build",
		CultivationBranch:"cultivate",
		MasterBranch:"master",
		BuildID:"test-uuid-for-build-id",
	}
	dao := NewInMemoryDAO()

	err := dao.Create(mockResult)
	if err != nil {t.Errorf("Error using In-memory DAO: %v", err)}

	resultMockDialog,err := dao.FetchByContextSubject(context, subject)

	if err != nil {t.Errorf("Couldn't fetch from In-memory DAO: %v", err)}

	if !reflect.DeepEqual(mockResult, resultMockDialog) {
		t.Errorf("Expected mockResult, %v\n resultMock, %v, to be the same.", mockResult, resultMockDialog)
	}
}
func CheckDialogEntryEquality(t *testing.T){
	nde1 := NewDialogEntry(
		"testing/context",
		"test_subject",
		"2019-12-12",
		[]string{"test1", "test2", "test3"},
		[]string{"comment1", "comment2", "comment3"},
		"test_dialog_id",
		"test_learn_more_link",
		"test_learn_more_content",
		"test_build_branch",
		"test_cultivation_branch",
		"test_master_branch",
		"test_build_id",
	)

	nde2 := DialogEntry{
		Context: "testing/context",
		Subject: "test_subject",
		Updated: "2019-12-12",
		Dialog: []string{"test1", "test2", "test3"},
		Comments: []string{"comment1", "comment2", "comment3"},
		DialogID: "test_dialog_id",
		LearnMoreLink: "test_learn_more_link",
		LearnMoreContent: "test_learn_more_content",
		BuildBranch: "test_build_branch",
		CultivationBranch: "test_cultivation_branch",
		MasterBranch: "test_master_branch",
		BuildID: "test_build_id",
	}

	if !reflect.DeepEqual(
		nde1,
		nde2,
	) {
		t.Errorf("Expected nde1, %v\n nde2, %v, to be the same.", nde1, nde2)
	}
}

// TestPanicsInNewDialogEntry can be run without localstack environment
func TestPanicsInNewDialogEntry(t *testing.T) {
	assert.Panics(t, func() {
		NewDialogEntry(
			"",
			"",
			"",
			[]string{"test1", "test2", "test3"},
			[]string{"comment1", "comment2", "comment3"},
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		)
	}, "The code did not panic")
	fmt.Println("Test_FetchBySubject: Completed")
}
