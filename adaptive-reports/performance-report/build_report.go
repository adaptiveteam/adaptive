package collaboration_report

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/unidoc/unipdf/v3/creator"
	"log"
	"math"
	"strings"
)

func buildReport(
// The last year of feedback received
	ReceivedBytes []byte,
// The last year of feedback given
	GivenBytes []byte,
// The users name (e.g., Chris Creel)
	UserName string,
// The quarter for which this report was produced
	Quarter int,
// The year for which this report was produced
	Year int,
// Name and location for where to store the file.
	FileName string,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
	getCompetencyUnsafe GetCompetencyUnsafe,
) (tags map[string]string, err error) {
	SetUniDocGlobalLicenseIfAvailable()
	c := creator.New()

	fm, err := getFontMap()
	if err == nil {
		received, err := newCoachingListFromStream(ReceivedBytes, getCompetencyUnsafe)
		if err == nil {
			logger.WithField("received", &received).Infof("Retrieved received feedback")
			given, err := newCoachingListFromStream(GivenBytes, getCompetencyUnsafe)
			if err == nil {
				logger.WithField("given", &given).Infof("Retrieved given feedback")

				receivedForQuarter := received.feedbackForQuarter(Quarter, Year)
				// givenForQuarter := given.FeedbackForQuarter(Quarter,Year)

				sortedTopics := received.getSortedAttribute(func(c coaching) string {
					return c.Topic
				})
				sortedTypes := received.getSortedAttribute(func(c coaching) string {
					return c.Type
				})
				topicToValueTypeMapping := received.getTopicToValueTypeMapping()
				documentLayout(c)
				documentFooters(c, fm)
				documentHeaders(c)
				documentFrontPage(UserName, Year, Quarter, c, fm)
				writePerformanceSummary(c, fm, receivedForQuarter)
				tags = writePerformanceAnalysis(c, fm, received, given, topicToValueTypeMapping, Quarter, Year, dialogDao, logger)
				writeCoachingIdeas(c, fm, receivedForQuarter, dialogDao)
				for _, each := range sortedTypes {
					s := writeFeedbackSummary(c, fm, receivedForQuarter, each, topicToValueTypeMapping)
					for _, topic := range sortedTopics {
						tpe := topicToValueTypeMapping[topic]
						if tpe == each {
							var language = "n/a"
							if !math.IsNaN(receivedForQuarter.topicCoaching(topic).calculateScore()) {
								language = getRatingLanguage(receivedForQuarter.topicCoaching(topic).calculateScore())
							}
							feedback := fmt.Sprintf(
								"%s (%s)",
								strings.Title(topic),
								language,
							)
							var sc *creator.Chapter
							if s != nil {
								sc = c.NewChapter(feedback)
							}
							writeTopic(c, sc, fm, topic, receivedForQuarter.topicCoaching(topic))
							_ = c.Draw(sc)
						}
					}
				}
				documentTableOfContents(c)
				err = c.WriteToFile(FileName)
				if err != nil {
					log.Println("Error writing file", err)
				}
			}
		}

	}
	return tags, err
}
