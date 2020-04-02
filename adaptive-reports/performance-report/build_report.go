package collaboration_report

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/values"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/unidoc/unipdf/v3/common/license"
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
	competencyDao values.DAO,
	logger logger.AdaptiveLogger,
) (tags map[string]string, err error) {

	uniDocLicense := `
-----BEGIN UNIDOC LICENSE KEY-----
eyJsaWNlbnNlX2lkIjoiMzY2Y2YyMDEtYmFiNi00ZTZkLTQyMDAtNWM5ODU5NmExNDY3IiwiY3VzdG9tZXJfaWQiOiIwYmUwMTk0My1iM2RlLTQwNTQtNTY3Yy04NjkyYzhjZTkyMzQiLCJjdXN0b21lcl9uYW1lIjoiQWRhcHRpdmUuVGVhbSIsImN1c3RvbWVyX2VtYWlsIjoiY2hyaXMuY3JlZWxAYWRhcHRpdmUudGVhbSIsInRpZXIiOiJidXNpbmVzcyIsImNyZWF0ZWRfYXQiOjE1ODQ5MDI0NjAsImV4cGlyZXNfYXQiOjE2MTY0NTc1OTksImNyZWF0b3JfbmFtZSI6IlVuaURvYyBTdXBwb3J0IiwiY3JlYXRvcl9lbWFpbCI6InN1cHBvcnRAdW5pZG9jLmlvIiwidW5pcGRmIjp0cnVlLCJ1bmlvZmZpY2UiOmZhbHNlLCJ0cmlhbCI6ZmFsc2V9
+
DwkMzUJbA6A7nJR7yHFn8C+aTec4EMRqjElgOB7doUsyRl5oFbvcO/KxKdbLSlmmVBQM01iq6FqThIcYlPqMGpvlFfXbPxGKK2Cf31CmhkdX4X0yV9fmkDcFgKTHllg4oUqzgKMglccXm1bqLgWNFnw6DXw1LNboYQ8Iebv6+CxZzN2viesEekWv4qCLJK+DgMnV2rldjreh3dmhcGT2A33vhqwSliSJeOp81MoVW1QMqDDYn2R8YjZ7mBKqs17m2/s3mek7zEpwoeQzGlspC/m1vPgP0yqAQwAdWslix/fK/9xruFbkAZCijwjNLOs/tuoUefb3DKHWH3RZ1bjrbg==
-----END UNIDOC LICENSE KEY-----
`
	licenseError := license.SetLicenseKey(
		uniDocLicense,
		"adaptive.team",
	)
	err = licenseError

	if err == nil {

		c := creator.New()

		var fm fontMap
		fm, err = getFontMap()
		if err == nil {
			received, err := newCoachingListFromStream(ReceivedBytes, competencyDao)
			if err == nil {
				logger.WithField("received", &received).Infof("Retrieved received feedback")

				var given coachingList
				given, err = newCoachingListFromStream(GivenBytes, competencyDao)
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
	}

	return
}
