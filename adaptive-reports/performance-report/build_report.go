package collaboration_report

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/unidoc/unipdf/v3/creator"
	"math"
	"strings"
)

func createPdfReport(
	// The last year of feedback received
	received CoachingList,
	// The users name (e.g., Chris Creel)
	UserName string,
	// The quarter for which this report was produced
	Quarter int,
	// The year for which this report was produced
	Year int,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
) (pdf *creator.Creator, tags map[string]string, err error) {
	SetUniDocGlobalLicenseIfAvailable()
	pdf = creator.New()

	var fm fontMap
	fm, err = getFontMap()
	if err == nil {
	
		receivedForQuarter := received.feedbackForQuarter(Quarter, Year)
		// givenForQuarter := given.FeedbackForQuarter(Quarter,Year)

		sortedTopics := received.getSortedAttribute(func(c Coaching) string {
			return c.Topic
		})
		sortedTypes := received.getSortedAttribute(func(c Coaching) string {
			return c.Type
		})
		topicToValueTypeMapping := received.getTopicToValueTypeMapping()
		documentLayout(pdf)
		documentFooters(pdf, fm)
		documentHeaders(pdf)
		documentFrontPage(UserName, Year, Quarter, pdf, fm)
		writePerformanceSummary(pdf, fm, receivedForQuarter)
		tags = writePerformanceAnalysis(pdf, fm, received, topicToValueTypeMapping, Quarter, Year, dialogDao, logger)
		writeCoachingIdeas(pdf, fm, receivedForQuarter, dialogDao)
		for _, each := range sortedTypes {
			s := writeFeedbackSummary(pdf, fm, receivedForQuarter, each, topicToValueTypeMapping)
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
						sc = pdf.NewChapter(feedback)
					}
					writeTopic(pdf, sc, fm, topic, receivedForQuarter.topicCoaching(topic))
					_ = pdf.Draw(sc)
				}
			}
		}
		documentTableOfContents(pdf)
	}
	return
}
