package collaboration_report

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	fetch_dialog "github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/sirupsen/logrus"
)

func TestBuildReportFullData_Test(t *testing.T) {
	received, err := readCoachingsFromJSON("test-data/test-data.json", MockGetCompetency)
	l := logger.LambdaLogger(logrus.InfoLevel)
	if err == nil {
		_, err = BuildReportWithCustomValuesTyped(
			received,
			"Christopher Creel",
			1,
			2019,
			"test-reports/adaptive.pdf",
			createTestDao(),
			l,
		)
	}
}

func TestBuildReportIncompleteData(t *testing.T) {
	received, err := readCoachingsFromJSON("test-data/test-incomplete-data.json", MockGetCompetency)
	l := logger.LambdaLogger(logrus.InfoLevel)
	if err == nil {
		_, err = BuildReportWithCustomValuesTyped(
			received,
			"Christopher Creel",
			1,
			2019,
			"test-reports/adaptive-incomplete.pdf",
			createTestDao(),
			l,
		)
	}
}

func MockGetCompetency(adaptiveValueID string)[]adaptiveValue.AdaptiveValue{
	return []models.AdaptiveValue{{ID: adaptiveValueID, Name: "V1", Description: "Desc1", ValueType: "performance"}}
}

func createTestDao() fetch_dialog.DAO {
	dialogFetcherDAO := fetch_dialog.NewInMemoryDAO()
	alias := func(a string) {
		if dialogFetcherDAO.CreateAlias(fetch_dialog.ContextAliasEntry{
			Context: "test_context",
			ApplicationAlias:   a,
		}) != nil {
			log.Panicf("unable to fetch dialog for alias %v", a)
		}
	}
	subject := func(s string) {
		if dialogFetcherDAO.Create(fetch_dialog.DialogEntry{
			Context: "test_context", Subject: s,
			Dialog: []string{"text for " + s + " - %v"},
		}) != nil {
			log.Panicf("unable to fetch dialog for subject %v", s)
		}
	}
	subject("summary-explanation")
	subject("not-enough-data")

	subject("exceeds")
	subject("meets")

	subject("positive")
	subject("negative")

	subject("large")
	subject("small")

	subject("below")

	alias("collaboration-report#quantity")
	alias("collaboration-report#overall")
	alias("collaboration-report#network")
	alias("collaboration-report#sentiment")
	alias("collaboration-report#energy")
	alias("collaboration-report#consistency")
	alias("collaboration-report#coachingIntro")
	alias("collaboration-report#efficiency")
	alias("collaboration-report#collaboration")
	alias("collaboration-report#communication")

	return dialogFetcherDAO
}

func readCoachingsFromJSON(file string, getCompetencyUnsafe GetCompetencyUnsafe) (received []Coaching, err error) {
	var receivedBytes []byte
	jsonFile, err := os.Open(file)
	if err == nil {
		receivedBytes, err = ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Panic("unable to read JSON file")
		}
		err = jsonFile.Close()
		if err != nil {
			log.Panic("unable to close JSON file")
		}
	} else {
		receivedBytes = nil
	}
	return NewCoachingListFromStream(receivedBytes, getCompetencyUnsafe)
}

func Test_getRating(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name       string
		args       args
		wantRating string
	}{
		{
			"does not meet 0",
			args{
				value: 0.0,
			},
			"did not meet expectations",
		},
		{
			"does not meet 1.0",
			args{
				value: 1.0,
			},
			"did not meet expectations",
		},
		{
			"does not meet 1.4",
			args{
				value: 1.4,
			},
			"did not meet expectations",
		},
		{
			"almost approaching 1.5",
			args{
				value: 1.5,
			},
			"is almost approaching expectations",
		},
		{
			"almost approaching 1.52",
			args{
				value: 1.52,
			},
			"is almost approaching expectations",
		},
		{
			"almost approaching 1.99",
			args{
				value: 1.99,
			},
			"is almost approaching expectations",
		},
		{
			"almost approaching 1.9999999999999",
			args{
				value: 1.9999999999999,
			},
			"is almost approaching expectations",
		},
		{
			"does not meet 1.9",
			args{
				value: 1.9,
			},
			"is almost approaching expectations",
		},
		{
			"approaching 2.0",
			args{
				value: 2.0,
			},
			"is approaching expectations",
		},
		{
			"approaching 2.3",
			args{
				value: 2.3,
			},
			"is approaching expectations",
		},
		{
			"better then approaching 2.5",
			args{
				value: 2.5,
			},
			"almost meets expectations",
		},
		{
			"better then approaching",
			args{
				value: 2.6,
			},
			"almost meets expectations",
		},
		{
			"better then approaching",
			args{
				value: 2.99,
			},
			"almost meets expectations",
		},
		{
			"meets 3.0",
			args{
				value: 3.0,
			},
			"meets expectations",
		},
		{
			"meets 3.9999999",
			args{
				value: 3.9999999,
			},
			"is almost above expectations",
		},
		{
			"meets 3.00000001",
			args{
				value: 3.00000001,
			},
			"meets expectations",
		},
		{
			"exceeds 5.0",
			args{
				value: 5.0,
			},
			"exceeded expectations",
		},
		{
			"exceeds 5.1",
			args{
				value: 5.1,
			},
			"exceeded expectations",
		},
		{
			"exceeds 4.50000000000001",
			args{
				value: 4.50000000000001,
			},
			"almost exceeded expectations",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRating := getRatingLanguage(tt.args.value); gotRating != tt.wantRating {
				t.Errorf("getRating() = %v, want %v", gotRating, tt.wantRating)
			}
		})
	}
}
