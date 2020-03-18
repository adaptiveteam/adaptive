package adaptive_utils_go

import (
	"log"
	"github.com/adaptiveteam/adaptive/daos/dialogEntry"
	"github.com/adaptiveteam/adaptive/daos/common"
	"errors"
	"fmt"
	"sort"

	nlp "github.com/adaptiveteam/adaptive/adaptive-nlp"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	dialogFetcher "github.com/adaptiveteam/adaptive/dialog-fetcher"
	eb "github.com/adaptiveteam/adaptive/engagement-builder"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

const (
	AnalysisFeedbackText = "I have some ideas to improve what you wrote below."
	// NoRecommendationsTemplate is the template when there are no recommendations.
	NoRecommendationsTemplate = "Your comments look good :thumbsup:"
)

var (
	defaultColor = models.BlueColorHex
	goodColor    = "#2ECC71" // Strong green, no issues found
	badColor     = "#E74C3C" // strong red
	colorMap     = map[int]string{
		0: goodColor,
		1: "#F9E79F", // Weak yellow, one issue found
		2: "#F1C40F", // Strong yellow, 2 issues found
		3: "#F5B7B1", // Weak red, 3 issues found
		4: badColor,
		5: badColor,
		6: badColor,
	}
)

func mapImprovementIDStringKeys(m map[nlp.ImprovementID]string) (keys []string) {
	for k := range m {
		keys = append(keys, string(k))
	}
	return
}

func sortedMapImprovementIDString(m map[nlp.ImprovementID]string) (op []models.KvPair) {
	keys := mapImprovementIDStringKeys(m)
	sort.Strings(keys)
	for _, k := range keys {
		op = append(op, models.KvPair{Key: k, Value: m[nlp.ImprovementID(k)]})
	}
	return op
}

// // Deprecated: It seems that it's not used outside.
// func TextRecommendations(text, context, dialogTableName, namespace string) (int, string) {
// 	d := awsutils.NewDynamo(NonEmptyEnv("AWS_REGION"), "", namespace)
// 	dialogFetcherDao := dialogFetcher.NewDAO(d, dialogTableName)
// 	//cnt, summary := textRecommendations(text, context, dialogFetcherDao, namespace)

// 	recs, errs := nlp.GetDialog(text, nlp.English, dialogFetcherDao, context)
// 	panicErrorList(errs, namespace)
// 	cnt := len(recs)
// 	_, isGood := recs[GoodDescriptionSubject]
// 	if len(recs) == 1 && isGood {
// 		cnt = 0
// 	}

// 	summary := RenderRecommendations(recs)

// 	return cnt, string(summary)
// }

func panicErrorList(errs []error, namespace string) {
	// Combining multiple errors into one
	var errsMkStr string
	var err error
	for _, each := range errs {
		errsMkStr = errsMkStr + " " + each.Error()
	}
	if errsMkStr != core.EmptyString {
		err = errors.New(errsMkStr)
	}

	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not get nlp recommendations for text"))
}

// RenderRecommendations joins recommendations into summary
func RenderRecommendations(recs map[nlp.ImprovementID]string) (res ui.RichText) {
	goodText, isGood := recs[nlp.GoodDescriptionSubject]
	res = ui.RichText("")
	if isGood {
		res = ui.RichText(goodText)
	} else {
		sortedRecMap := sortedMapImprovementIDString(recs)
		for _, rec := range sortedRecMap {
			res = res + ui.RichText(rec.Key).Bold() + "\n" + ui.RichText(rec.Value) + "\n\n"
		}
	}
	return
}

// ColorStatusByRecommendationsLength returns color based on the number of recommendations
// Deprecated: Use ColorStatusByRecommendations
func ColorStatusByRecommendationsLength(len int) (color string) {
	color = defaultColor
	if col, ok := colorMap[len]; ok {
		// Setting color according to the number of issues in Meaning Cloud analysis
		color = col
	}
	return
}

func ColorStatusByTextAnalysisResults(textAnalysisResults TextAnalysisResults) (color string) {
	return ColorStatusByIsGoodAndLength(textAnalysisResults.IsGood, textAnalysisResults.RecommendationsCount)
}

func ColorStatusByIsGoodAndLength(isGood bool, recommendationsCount int) (color string) {
	if isGood {
		color = goodColor
	} else {
		if recommendationsCount >= len(colorMap) {
			color = badColor
		} else {
			col, hasColor := colorMap[recommendationsCount]
			if hasColor {
				// Setting color according to the number of issues in Meaning Cloud analysis
				color = col
			} else {
				color = defaultColor
			}
		}
	}
	return
}

// ColorStatusByRecommendations returns bar color for the given recommendations.
func ColorStatusByRecommendations(recs map[nlp.ImprovementID]string) (color string) {
	_, isGood := recs["good"]
	return ColorStatusByIsGoodAndLength(isGood, len(recs))
}

// RepaintAttachmentsWithColor repaints the given list of attachments with the provided color
func RepaintAttachmentsWithColor(attachs []ebm.Attachment, color string) (colorCodedAttachs []ebm.Attachment) {
	for _, each := range attachs {
		newAttach, _ := eb.LoadAttachmentBuilder(&each).Color(color).Build()
		colorCodedAttachs = append(colorCodedAttachs, *newAttach)
	}
	return
}

// RecommendationsTemplate renders recommendation in the form of attachments
func RecommendationsTemplate(originalText string, analysisSummary ui.RichText, color string) []ebm.Attachment {
	pretext := core.IfThenElse(color == goodColor, NoRecommendationsTemplate, AnalysisFeedbackText).(string)
	attach1, _ := eb.NewAttachmentBuilder().
		Color(color).
		Pretext(core.TextWrap(pretext, core.Underscore)).
		MarkDownIn([]ebm.MarkdownField{ebm.MarkdownFieldText, ebm.MarkdownFieldPretext}).
		Text(string(analysisSummary)).
		Title(core.EmptyString).
		Build()
	attach2, _ := eb.NewAttachmentBuilder().
		Title("Original Text").
		Text(fmt.Sprintf("%s", originalText)).
		Build()
	return []ebm.Attachment{*attach1, *attach2}
}

// RecommendationsMessage constructs PlatformSimpleNotification that will contain recommendations.
func RecommendationsMessage(originalText string, analysisSummary ui.RichText, color string) (note models.PlatformSimpleNotification) {
	if analysisSummary != "" {
		attachments := RecommendationsTemplate(originalText, ui.RichText(analysisSummary), color)
		note = models.PlatformSimpleNotification{Attachments: attachments}
	} else {
		// When no feedback, post the feedback is good
		note = models.PlatformSimpleNotification{Message: string(NoRecommendationsTemplate)}
	}
	return
}

type ThreadID struct {
	UserID   string
	ThreadTs string
}

type OriginalMessageID struct {
	UserID    string
	ChannelID string
	Ts        string
}

// SendNoteToUserThread updates the message to send it to specific thread and user.
func (conversationContext ConversationContext) SendNoteToUserThread(note models.PlatformSimpleNotification) models.PlatformSimpleNotification {
	note.AsUser = true
	note.UserId = conversationContext.UserID
	note.Channel = conversationContext.ConversationID
	note.ThreadTs = conversationContext.ThreadTs
	return note
}

// UpdateOriginalMessageInUserChannel updates the message to send it to specific channel and user to override message with given ts.
func (conversationContext ConversationContext) UpdateOriginalMessageInUserChannel(note models.PlatformSimpleNotification) models.PlatformSimpleNotification {
	note.AsUser = true
	note.UserId = conversationContext.UserID
	note.Channel = conversationContext.ConversationID
	note.Ts = conversationContext.OriginalMessageTs
	return note
}

// ECAnalysis performs analysis and posts recommendations to the same thread
// ECAnalysis breaks SRP. Uses environment variables internally.
// Deprecated: Use AnalyzeText and PresentTextAnalysisResults instead.
func ECAnalysis(originalText, context, label, dialogTableName string, callbackId, userId, channelId, ts, threadTs string, updatedAttachs []ebm.Attachment, s *awsutils.SnsRequest, notificationTopic, namespace string) {
	// Do analysis on the text
	// Once we receive the analysis from Meaning Cloud on the user's feedback, we post that result to the original message's thread
	d := awsutils.NewDynamo(NonEmptyEnv("AWS_REGION"), "", namespace)
	dialogFetcherDao := dialogFetcher.NewDAO(d, dialogTableName)
	conversationID := channelId
	if conversationID == "" {
		conversationID = userId
	}
	conversationContext := ConversationContext{
		ThreadTs:          threadTs,
		UserID:            userId,
		ConversationID:    conversationID,
		OriginalMessageTs: ts,
	}
	platform := Platform{
		Sns:                       *s,
		PlatformNotificationTopic: notificationTopic,
		Namespace:                 namespace,
	}

	input := TextAnalysisInput{
		Text:                       originalText,
		OriginalMessageAttachments: updatedAttachs,
		Context:                    context,
		Namespace:                  namespace,
	}
	recommendations, errList := nlp.GetImprovements(originalText, nlp.English)
	recs, errList2 := nlp.FetchDialogForImprovements(dialogFetcherDao, context, recommendations)
	errs := append(errList, errList2...)

	panicErrorList(errs, namespace)
	_, isGood := recs["good"]
	recommendationsCount := len(recs)
	if isGood {
		recommendationsCount = 0
	}
	analysisSummary := RenderRecommendations(recs)
	analysisResults := TextAnalysisResults{
		RecommendationsCount: recommendationsCount,
		IsGood:               isGood,
		Summary:              analysisSummary,
		TextAnalysisInput:    input,
	}
	notes := conversationContext.PresentTextAnalysisResults(analysisResults)
	// 	[]models.PlatformSimpleNotification{
	// 	conversationContext.SendNoteToUserThread(note),
	// 	conversationContext.UpdateOriginalMessageInUserChannel(colorCodedOriginalMessageOverrideNote),
	// }
	platform.PublishAll(notes)
}

// TextAnalysisInput captures the input arguments for text analysis
type TextAnalysisInput struct {
	Text                       string
	OriginalMessageAttachments []ebm.Attachment
	Context                    string
	Namespace                  string
}

// TextAnalysisResults encapsulates information obtained from text analysis.
type TextAnalysisResults struct {
	TextAnalysisInput
	RecommendationsCount int
	IsGood               bool
	Summary              ui.RichText
}

// GetConversationID returns either channel id or user id whatever is not empty
func GetConversationID(userID, channelID string) string {
	if channelID != "" {
		return channelID
	}
	return userID
}

// ConversationContext is a data structure that represents current conversation with the user.
type ConversationContext struct {
	// UserID is used to retrieve platform information - token, id...
	UserID string
	// ConversationID is either channel id or user id. It's used to
	ConversationID    string
	OriginalMessageTs string
	ThreadTs          string
}

// AnalyzeTextC performs a few checks of the input text and produces some recommendations.
// This is the same as AnalyzeText apart from using DynamoDBConnection instead of dialogFetcher.DAO
func AnalyzeTextC(input TextAnalysisInput)func (common.DynamoDBConnection)(result TextAnalysisResults, err error) {
	return func (conn common.DynamoDBConnection)(result TextAnalysisResults, err error) {
		dialogFetcherDao := dialogFetcher.NewDAO(conn.Dynamo, dialogEntry.TableName(conn.ClientID))
		var errors []error
		result, errors = AnalyzeText(dialogFetcherDao, input)
		for _, e := range errors {
			log.Printf("AnalyzeTextC ERROR: %v\n", e)
			err = errors[0]
		}
		return
	}
}
// AnalyzeText performs a few checks of the input text and produces some recommendations.
func AnalyzeText(dialogFetcherDao dialogFetcher.DAO, input TextAnalysisInput) (result TextAnalysisResults, errors []error) {
	recommendations, errList := nlp.GetImprovements(input.Text, nlp.English)
	recs, errList2 := nlp.FetchDialogForImprovementsOrGood(dialogFetcherDao, input.Context, recommendations)
	errs := append(errList, errList2...)

	l := len(recs)
	_, isGood := recs[nlp.GoodDescriptionSubject]
	if isGood {
		l = 0
	}
	return TextAnalysisResults{
		TextAnalysisInput:    input,
		RecommendationsCount: l,
		IsGood:               isGood,
		Summary:              RenderRecommendations(recs),
	}, errs
}

// AnalyzeTextUnsafe performs a few checks of the input text and produces some recommendations.
// panics in case of any errors.
func AnalyzeTextUnsafe(dialogFetcherDao dialogFetcher.DAO, input TextAnalysisInput) TextAnalysisResults {
	result, errors := AnalyzeText(dialogFetcherDao, input)
	panicErrorList(errors, input.Namespace)
	return result
}

// PresentTextAnalysisResults represents text analysis results to user in the given conversation context.
func (conversationContext ConversationContext) PresentTextAnalysisResults(analysisResults TextAnalysisResults) []models.PlatformSimpleNotification {
	color := ColorStatusByIsGoodAndLength(analysisResults.IsGood, analysisResults.RecommendationsCount)

	note := RecommendationsMessage(analysisResults.TextAnalysisInput.Text, analysisResults.Summary, color)
	// Update the original attachments with the new color
	colorCodedOriginalMessageOverrideNote := models.PlatformSimpleNotification{
		Attachments: RepaintAttachmentsWithColor(analysisResults.TextAnalysisInput.OriginalMessageAttachments, color)}
	return []models.PlatformSimpleNotification{
		conversationContext.SendNoteToUserThread(note),
		conversationContext.UpdateOriginalMessageInUserChannel(colorCodedOriginalMessageOverrideNote),
	}
}
