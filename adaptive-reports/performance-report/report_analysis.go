package collaboration_report

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-nlp"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/business-time"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/gonum/stat"
	"github.com/unidoc/unidoc/pdf/creator"
)

const (
	StrongRed    = "StrongRed"
	WeakRed      = "WeakRed"
	StrongYellow = "StrongYellow"
	WeakYellow   = "WeakYellow"
	WeakGreen    = "WeakGreen"
	StrongGreen  = "StrongGreen"
	Neutral      = "Neutral"
)

const (
	packageName   = "collaboration-report"
	Quantity      = "quantity"
	Overall       = "overall"
	Network       = "network"
	Sentiment     = "sentiment"
	Energy        = "energy"
	Consistency   = "consistency"
	CoachingIntro = "coachingIntro"
)

func getColor(severity string) (rv creator.Color) {
	severityMap := map[string]creator.Color{
		StrongRed:    creator.ColorRGBFromHex("#C0392B"),
		WeakRed:      creator.ColorRGBFromHex("#F5B7B1"),
		StrongYellow: creator.ColorRGBFromHex("#F1C40F"),
		WeakYellow:   creator.ColorRGBFromHex("#F9E79F"),
		WeakGreen:    creator.ColorRGBFromHex("#ABEBC6"),
		StrongGreen:  creator.ColorRGBFromHex("#27AE60"),
		Neutral:      creator.ColorRGBFromHex("#AEB6BF"),
	}

	rv = severityMap[severity]
	return rv
}

func generateSummaryAnalysis(
	received CoachingList,
	topicToValueTypeMapping map[string]string,
	quarter int,
	year int,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
) (rv map[string]string, tags map[string]string) {

	rv = make(map[string]string, 0)
	tags = make(map[string]string, 0)
	quarterAnalysisFunctions := getSingleQuarterAnalysisFunctions(topicToValueTypeMapping, dialogDao, logger)
	for key, value := range quarterAnalysisFunctions {
		analysis, advice := value(received)
		tags[key] = analysis
		rv[key] = advice
	}
	for key, value := range getMultiQuarterAnalysisFunctions(dialogDao) {
		analysis, advice := value(quarter, year, received)
		tags[key] = analysis
		rv[key] = advice
	}

	return rv, tags
}

type analysisFunction = func(list CoachingList) (string, string)

func getSingleQuarterAnalysisFunctions(
	topicToValueTypeMapping map[string]string,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
) (rv map[string]func(CoachingList) (analysis string, advice string)) {
	rv = make(map[string]analysisFunction, 0)
	rv["Feedback Quantity"] = feedbackQuantity(dialogDao)
	rv["Overall"] = quarterPerformanceOverall(topicToValueTypeMapping, dialogDao)
	rv["Network Strength"] = quarterNetworkOverall(dialogDao)
	rv["Sentiment"] = quarterSentimentOverall(dialogDao, logger)
	rv["Relationship"] = relationshipOverview(topicToValueTypeMapping, dialogDao)

	return rv
}

func getMultiQuarterAnalysisFunctions(dialogDao fetch_dialog.DAO) (rv map[string]func(quarter int, year int, list CoachingList) (analysis string, advice string)) {
	rv = make(map[string]func(quarter int, year int, list CoachingList) (string, string), 0)
	rv["Consistency"] = consistencyOverview(dialogDao)
	return rv
}

func feedbackQuantity(dialogDao fetch_dialog.DAO) analysisFunction {
	return func(c CoachingList) (analysis string, advice string) {
		topics := c.topics()
		score := float32(c.justFeedback().length()) / float32(len(topics))

		var subject string
		if score < 4 {
			subject = "below"
			analysis = StrongRed
		} else if score >= 4 && score < 7 {
			subject = "meets"
			analysis = WeakYellow
		} else if score >= 7 && score < 10 {
			subject = "above"
			analysis = WeakGreen
		} else {
			subject = "exceeds"
			analysis = StrongGreen
		}

		options := loadDialogUnsafe(dialogDao, Quantity, subject)

		advice = fmt.Sprintf(core_utils_go.RandomString(options.Dialog), score)
		return analysis, advice
	}
}

func quarterPerformanceOverall(topicToValueTypeMapping map[string]string, dialogDao fetch_dialog.DAO) analysisFunction {
	return func(c CoachingList) (analysis string, advice string) {
		performanceFeedback := c.kindCoaching("performance", topicToValueTypeMapping)
		var subject, language string
		if len(performanceFeedback) > 0 {
			score := c.kindCoaching("performance", topicToValueTypeMapping).calculateScore()
			language = getRatingLanguage(score)

			if score < 2 {
				subject = "below"
				analysis = StrongRed
			} else if score >= 2 && score < 3 {
				subject = "meets"
				analysis = WeakYellow
			} else if score >= 4 && score < 4 {
				subject = "above"
				analysis = WeakGreen
			} else {
				subject = "exceeds"
				analysis = StrongGreen
			}

			options := loadDialogUnsafe(dialogDao, Overall, subject)
			advice = fmt.Sprintf(core_utils_go.RandomString(options.Dialog), language)
		} else {
			analysis = Neutral
			options := loadDialogUnsafe(dialogDao, Overall, "not-enough-data")
			advice = core_utils_go.RandomString(options.Dialog)
		}

		return analysis, advice
	}
}

func quarterNetworkOverall(dialogDao fetch_dialog.DAO) analysisFunction {
	return func(c CoachingList) (analysis string, advice string) {
		var sources []string
		for i := 0; i < c.length(); i++ {
			sources = append(sources, c.index(i).GetSource())
		}
		sources = unique(sources)
		score := len(sources)

		var subject string
		if score < 4 {
			subject = "small"
			analysis = StrongYellow
		} else if score >= 4 && score < 7 {
			subject = "medium"
			analysis = WeakGreen
		} else if score >= 7 && score < 10 {
			subject = "large"
			analysis = StrongGreen
		} else {
			subject = "jumbo"
			analysis = StrongYellow
		}
		options := loadDialogUnsafe(dialogDao, Network, subject)
		advice = fmt.Sprintf(core_utils_go.RandomString(options.Dialog), score)
		return analysis, advice
	}
}

func quarterSentimentOverall(dialogDao fetch_dialog.DAO, logger logger.AdaptiveLogger) analysisFunction {
	return func(c CoachingList) (analysis string, advice string) {
		score, err := nlp.GetTextSentimentText(c.createTextBlob(), nlp.English)
		if err == nil {
			var subject string
			if score.GetSentiment() == "very positive" {
				subject = "very-positive"
				analysis = StrongGreen
			} else if score.GetSentiment() == "positive" {
				subject = "positive"
				analysis = WeakGreen
			} else if score.GetSentiment() == "none" || score.GetSentiment() == "neutral" {
				subject = "neutral"
				analysis = WeakYellow
			} else if score.GetSentiment() == "negative" {
				subject = "negative"
				analysis = WeakRed
			} else if score.GetSentiment() == "very negative" {
				subject = "very-negative"
				analysis = StrongRed
			}
			options := loadDialogUnsafe(dialogDao, Sentiment, subject)
			advice = core_utils_go.RandomString(options.Dialog)
		} else {
			logger.WithError(err).Errorf("Could not get sentiment text")
		}
		return analysis, advice
	}
}

func relationshipOverview(topicToValueTypeMapping map[string]string, dialogDao fetch_dialog.DAO) analysisFunction {
	return func(c CoachingList) (analysis string, advice string) {
		relationshipFeedback := c.typeCoaching("relationship", topicToValueTypeMapping)
		var subject, language string
		if len(relationshipFeedback) > 0 {
			score := c.typeCoaching("relationship", topicToValueTypeMapping).calculateScore()
			language = getRatingLanguage(score)

			if score < 2 {
				subject = "below"
				analysis = StrongRed
			} else if score >= 2 && score < 3 {
				subject = "meets"
				analysis = WeakYellow
			} else if score >= 3 && score < 4 {
				subject = "above"
				analysis = WeakGreen
			} else {
				subject = "exceeds"
				analysis = StrongGreen
			}
			options := loadDialogUnsafe(dialogDao, Energy, subject)
			advice = fmt.Sprintf(core_utils_go.RandomString(options.Dialog), language)
		} else {
			analysis = Neutral
			options := loadDialogUnsafe(dialogDao, Energy, "not-enough-data")
			advice = core_utils_go.RandomString(options.Dialog)
		}

		return analysis, advice
	}
}

func consistencyOverview(dialogDao fetch_dialog.DAO) func(
	quarter int,
	year int,
	c CoachingList,
) (analysis string, advice string) {
	return func(
		quarter int,
		year int,
		c CoachingList,
	) (analysis string, advice string) {

		quarters := []business_time.Date{
			business_time.NewDateFromQuarter(quarter, year),
		}

		for i := 1; i < 4; i++ {
			previousQuarter := business_time.NewDateFromQuarter(
				quarters[i-1].GetPreviousQuarter(),
				quarters[i-1].GetPreviousQuarterYear(),
			)
			quarters = append(quarters, previousQuarter)
		}

		var scores []float64

		for _, each := range quarters {
			qf := c.feedbackForQuarter(each.GetQuarter(), each.GetYear()).justScores()
			if len(qf) > 0 {
				scores = append(scores, stat.Mean(qf, nil))
			}
		}

		var subject string
		if len(scores) > 2 {
			score := stat.Variance(scores, nil)

			if score > 80 {
				subject = "below"
				analysis = StrongRed
			} else if score > 60 && score <= 80 {
				subject = "meets"
				analysis = WeakYellow
			} else if score > 40 && score <= 60 {
				subject = "above"
				analysis = WeakGreen
			} else {
				subject = "exceeds"
				analysis = StrongGreen
			}

			options := loadDialogUnsafe(dialogDao, Consistency, subject)

			advice = fmt.Sprintf(core_utils_go.RandomString(options.Dialog), score)
		} else {
			analysis = Neutral
			options := loadDialogUnsafe(dialogDao, Consistency, "not-enough-data")

			advice = core_utils_go.RandomString(options.Dialog)
		}

		return analysis, advice
	}
}

func coachingIdeaAnalysis(
	dialogDao fetch_dialog.DAO,
) (coachingIntro string) {
	options := loadDialogUnsafe(dialogDao, CoachingIntro, "summary-explanation")
	coachingIntro = core_utils_go.RandomString(options.Dialog)
	return coachingIntro
}

func loadDialogUnsafe(dialogDao fetch_dialog.DAO, contextAlias, subject string) fetch_dialog.DialogEntry {
	options, err := dialogDao.FetchByAlias(packageName, contextAlias, subject)
	core_utils_go.ErrorHandler(err, "reporting", "failed to load dialog for `"+packageName+"#"+contextAlias+":"+subject+"`")
	return options
}
