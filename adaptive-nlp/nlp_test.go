package nlp

import (
	"fmt"
	fetch_dialog "github.com/adaptiveteam/dialog-fetcher"
	"net/url"
	"os"
	"testing"
)

func testFetcherWithReturn(
	tableName string,
	subject string,
	context string,
) (dialogOptions []string, err error) {
	dialogOptions = []string{"Hello","World"}
	return  dialogOptions,err
}

func testFetcherWithoutReturn(
	tableName string,
	subject string,
	context string,
) (dialogOptions []string, err error) {
	return  dialogOptions,nil
}

const posBreakingText = "This is a test sentence to demonstrate breaking out parts of speech. " +
	"We also want to derive sentiment." +
	"I hate that." +
	"I wonder how I ended up with ths job.  Oh yeah.  I chose it."

func improvementsContains(improvements []ImprovementID, improvement ImprovementID) bool {
	for _, i := range improvements {
		if i == improvement {
			return true
		}
	}
	return false
}

// ConfigureGlobalEnvironmentVariable sets AWS_REGION for running real
// services during test.
// An alternative is to use mock-services.
func ConfigureGlobalEnvironmentVariable() {
	_ = os.Setenv("AWS_REGION", "us-west-2")
}
func Test_all(t *testing.T) {
	ConfigureGlobalEnvironmentVariable()
	ensureGlobalConnectionsAreOpen()

	improvements, _ := globalConnections.GetImprovements("but", English)
	if len(improvements) < 3 {
		t.Errorf("Expected 4 or more recommndations for 'but.' Got %v: %v", len(improvements), improvements)
	}
	if  !improvementsContains(improvements, TooShort) ||
		!improvementsContains(improvements, WordsToAvoid) ||
		!improvementsContains(improvements, Actionable) ||
		!improvementsContains(improvements, Subjectivity)  {
		t.Errorf("Expected known recommndations for 'but.' Got: %v", improvements)
	}

	objectiveText := "The sky is blue. So is the ocean."
	u, err := url.Parse("https://raw.githubusercontent.com/adaptiveteam/chris.creel/master/README.md")
	if err == nil {
		textCategories, err := GetTextCategoriesURL(u, English)
		if err == nil {
			for i := 0; i < len(textCategories); i++ {
				fmt.Println(
					textCategories[i].GetLabel(),
					textCategories[i].GetAbsRelevance(),
					textCategories[i].GetRelevance(),
					textCategories[i].GetSentiment(),
				)
			}
		}
	}

	s, err := GetTextSentimentURL(u, English)
	if err == nil {
		fmt.Println(
			"Sentiment", s.GetSentiment(),
			"Subjectivity", s.GetSubjectivity(),
			"Irony", s.GetIrony(),
			"Agreement", s.GetAgreement(),
			"Confidence", s.GetConfidence(),
		)
	} else {
		t.Errorf("Error in GetTextSentimentURL %v", err)
	}

	summary, err := GetSummaryURL(5, u, English)
	if err == nil {
		fmt.Println(summary)
	} else {
		t.Errorf("Error in GetSummaryURL %v", err)
	}

	pos, err := globalConnections.GetPartOfSpeech(posBreakingText, English)
	if err == nil {
		fmt.Println(pos)
	} else {
		t.Errorf("Error in GetPartOfSpeech %v", err)
	}

	ac := IsActionable(pos)	
	fmt.Println("Action oriented phrase", ac)

	ps := IsPersonal(pos)
	fmt.Println("Feedback is personal", ps)

	tt, err := globalConnections.GetTranslation(posBreakingText, English, Spanish)
	if err == nil {
		fmt.Println(tt)
	} else {
		t.Errorf("Error in GetTranslation %v", err)
	}

	tt, err = globalConnections.GetTranslation(posBreakingText, English, Spanish)
	if err == nil {
		fmt.Println(tt)
	} else {
		t.Errorf("Error in GetTranslation %v", err)
	}

	tt, err = globalConnections.GetTranslation(posBreakingText, Korean, Hebrew)
	if err == nil {
		fmt.Printf("Expected an error on translation from Korean to Hebrew. Got a real translation due to language autodetection: %s", tt)
	}

	tc, err := GetTextCategoriesText(posBreakingText, English)
	if err == nil {
		fmt.Println(tc)
	} else {
		t.Errorf("Error in GetTextCategoriesText %v", err)
	}

	// Trying to invoke `Missing required parameter(s)` error
	_, err = GetTextCategoriesText("", English)
	if err != nil {
		t.Errorf("Got an error on an empty string")
	}

	//ironicText := "I love waking up with migraines."
	ironicText := "Get AWS Professional Solutions Architect Certification. I am going to identify potential courses in Udemy that could help me with this. I will review these with my colleagues to see how helpful they would be."
	ts, err := GetTextSentimentText(ironicText, English)
	if err == nil {
		if ts.GetIrony() {
			t.Errorf("Error in detecting irony %v", err)
		}
	} else {
		t.Errorf("Error in GetTextSentimentText %v", err)
	}

	ts, err = GetTextSentimentText(objectiveText, English)
	if err == nil {
		if ts.GetSubjectivity() {
			t.Error("Expected to detect objective statement but got subjective")
		}
	} else {
		t.Errorf("Error in GetTextSentimentText %v", err)
	}

	st, err := GetSummaryText(1, posBreakingText, English)
	if err == nil {
		fmt.Println(st)
	} else {
		t.Errorf("Error in GetSummaryText %v", err)
	}

	// Trying to invoke `Missing required parameter(s)` error
	_, err = GetSummaryText(1, "", English)
	if err != nil {
		t.Errorf("Got an error on an empty string")
	}

	_, err = globalConnections.MeaningCloud.GetTextSentiment("wrong", "still wrong", English)
	if err == nil {
		t.Error("Expected an error on bad parameters for getTextSentiment but got none")
	}

	dialogFetcherDAO := fetch_dialog.NewInMemoryDAO()

	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "too-short",
		Dialog: []string{"Your text is too short."},
	})
	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "subjective",
		Dialog: []string{"Your text is not subjective enough."},
	})
	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "avoid-words",
		Dialog: []string{"There are some words that you might prefer to avoid."},
	})
	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "actionable",
		Dialog: []string{"If you could do your language more actionable..."},
	})
	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "not-positive",
		Dialog: []string{"Try making your phrase more positive. It'll have much stronger impact."},
	})
	dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
		Context: "test_context", Subject: "good",
		Dialog: []string{"Your text is good."},
	})

	_, err = dialogFetcherDAO.FetchByContextSubject("test_context", "too-short")
	if err != nil {
		t.Errorf("err %v", err)
	}
	recommendations, errList := GetDialog(
		ironicText,
		English,
		dialogFetcherDAO, "test_context",
	)

	recommendations, errList = GetDialog(
		"The ocean is blue.",
		English,
		dialogFetcherDAO, "test_context",
	)
	if len(errList) == 0 {
		if len(recommendations) == 0 {
			t.Error("Expected recommendations but got none")
		} else {
			fmt.Println(recommendations)
		}
	} else {
		for _, err := range errList {
			t.Errorf("Error in GetDialog: %v", err)
		}
	}

	recommendations, errList = GetDialog(
		"You are a terrible person. But, I love you so much",
		English,
		dialogFetcherDAO, "test_context",
	)
	if len(errList) == 0 {
		if len(recommendations) == 0 {
			t.Error("Expected recommendations but got none")
		} else {
			fmt.Println(recommendations)
		}
	} else {
		for _, err := range errList {
			t.Errorf("Error in GetDialog 2: %v", err)
		}
	}

	recommendations, errList = GetDialog(
		"I hate you so very much.",
		English,
		dialogFetcherDAO, "test_context",
	)
	if len(errList) == 0 {
		if len(recommendations) == 0 {
			t.Error("Expected recommendations but got none")
		} else {
			fmt.Println(recommendations)
		}
	} else {
		for _, err := range errList {
			t.Errorf("Error in GetDialog 2: %v", err)
		}
	}

	wordsToAvoid := GetWordsToAvoid("The ocean is blue. But so are you")
	if wordsToAvoid == true {
		fmt.Println("Identified words to avoid.")
	} else {
		t.Error("Did not detect words to avoid.")
	}

	tooShort := GetTooShort("The ocean is blue, but so are you")
	if tooShort == true {
		fmt.Println(tooShort, "is too short")
	} else {
		t.Error("Did not detect too short text")
	}

	notTooShort := GetTooShort("I know you were trying your best to do a good job.  Hopefully, we can learn from the experience and you will come out a better person in the end.  Next time, let's sit down ahead of the sprint to better understand what the expectations are.")
	if notTooShort != true {
		fmt.Println(tooShort, "is correctly identified as not too short")
	} else {
		t.Error("Incorrectly identified as too short.")
	}

	tooLong := GetTooLong("The ocean is blue, but so are you. The ocean is blue, but so are you. The ocean is blue, but so are you. The ocean is blue, but so are you. The ocean is blue, but so are you")
	if wordsToAvoid == true {
		fmt.Println(tooLong, "is too long")
	} else {
		t.Error("Did not detect too long text")
	}

	// Testing empty strings
	recommendations, errList = GetDialog(
		"",
		English,
		dialogFetcherDAO, "test_context",
	)
	if len(recommendations) != 1 || len(errList) > 0 {
		t.Errorf("Expected graceful failure with empty text. %v\n%v", recommendations, errList)
	}
	summary, err = globalConnections.MeaningCloud.GetSummary(5, "txt", "", English)
	if len(summary) > 0 || err != nil {
		t.Error("Expected graceful failure with empty text.")
	}
	sentiment, err := globalConnections.MeaningCloud.GetTextSentiment("txt", "", English)
	if sentiment.GetSentiment() != "none" || err != nil {
		t.Error("Expected graceful failure with empty text.")
	}

	// testing bad parameters
	_, err = globalConnections.MeaningCloud.GetTextSentiment("BAD!!!", "", English)
	if err == nil {
		t.Error("Expected graceful failure with bad parameters.")
	}

	_, err = globalConnections.MeaningCloud.GetTextCategories("BAD!!!", "", English)
	if err == nil {
		t.Error("Expected graceful failure with empty parameters.")
	}

	_, err = globalConnections.MeaningCloud.GetSummary(5, "BAD!!!", "", English)
	if err == nil {
		t.Error("Expected graceful failure with empty parameters.")
	}

	summary, err = globalConnections.MeaningCloud.GetSummary(5, "url", "https://en.wikipedia.org/wiki/Organizational_chart", English)
	if err != nil {
		t.Error("Expected graceful failure with empty parameters.")
	}
}
