package collaboration_report

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-nlp"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/logger"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/dialog-fetcher"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// Margins
	marginLeft   = 50
	marginRight  = 50
	marginBottom = 65
	marginTop    = 50
	// Font sizes
	fontSizeTitle                 = 48
	fontSizeSubTitle              = 20
	fontSizeDateLine              = 14
	fontSizeHeadingOne            = 18
	fontSizeHeadingTwo            = 16
	fontSizeHeadingThree          = 14
	fontSizeNormal                = 10
	fontSizeTableOfContentsHeader = 28
	fontSizeTableOfContentsLine   = 14
	// Padding
	padding = 5
	// Tables
	borderWidth = 1
	// Summarization
	summaryBig   = 5
	summarySmall = 3
	// Indentation
	indentOne   = 10
	indentTwo   = 20
	indentThree = 30
	indentFour  = 40
	// Footer
	footerWidth = 40
	// Coaching Topic Tolerance
	coachingTopicTolerance = 80
)

func documentLayout(c *creator.Creator) {
	c.SetPageSize(creator.PageSizeLetter)
	c.SetPageMargins(marginLeft, marginRight, marginTop, marginBottom)
}

func documentFooters(c *creator.Creator, fonts fontMap) {
	c.DrawFooter(func(block *creator.Block, args creator.FooterFunctionArgs) {
		line := c.NewLine(0, footerWidth, c.Width(), footerWidth)
		line.SetLineWidth(footerWidth)
		line.SetColor(creator.ColorBlack)

		p := newParagraph(
			c,
			false,
			fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages),
			fonts["Bold"],
			creator.ColorRGBFromHex("#FFFFFF"),
			creator.TextAlignmentLeft,
			fontSizeNormal,
			0,
			0,
			0,
			0,
		)
		p.SetPos(280, 15)

		_ = block.Draw(line)
		_ = block.Draw(p)
	},
	)
}

func documentHeaders(c *creator.Creator) {
	// Do nothing for right now
}

func documentTableOfContents(c *creator.Creator) {
	c.AddTOC = true
	toc := c.TOC()
	hstyle := c.NewTextStyle()
	hstyle.Color = creator.ColorRGBFromArithmetic(0.2, 0.2, 0.2)
	hstyle.FontSize = fontSizeTableOfContentsHeader
	toc.SetHeading("Table of Contents", hstyle)
	lstyle := c.NewTextStyle()
	lstyle.FontSize = fontSizeTableOfContentsLine
	toc.SetLineStyle(lstyle)
}

func documentFrontPage(
	userName string,
	year int,
	quarter int,
	c *creator.Creator,
	f fontMap,
) {
	c.CreateFrontPage(func(args creator.FrontpageFunctionArgs) {

		name := newParagraph(
			c,
			false,
			userName,
			f["Regular"],
			creator.ColorBlack,
			creator.TextAlignmentCenter,
			fontSizeTitle,
			0,
			0,
			0,
			0,
		)

		reportTitle := newParagraph(
			c,
			false,
			"Collaboration Report for Q"+strconv.Itoa(quarter)+" "+strconv.Itoa(year),
			f["Bold"],
			creator.ColorBlack,
			creator.TextAlignmentCenter,
			fontSizeSubTitle,
			0,
			0,
			padding,
			0,
		)

		t := time.Now().UTC()
		dateStr := "Generated on - " + t.Format("January 02, 2006")

		dateLine := newParagraph(
			c,
			false,
			dateStr,
			f["Bold"],
			creator.ColorBlack,
			creator.TextAlignmentCenter,
			fontSizeDateLine,
			0,
			0,
			padding,
			0,
		)

		combinedHeight := name.Height() + reportTitle.Height() + dateLine.Height()
		pageHeight := c.Height() - 100 // 100 is a magic number.  I have no idea why I need to add this.
		center := pageHeight/2 - combinedHeight/2
		name.SetMargins(0, 0, center, 0)
		_ = c.Draw(name)
		_ = c.Draw(reportTitle)
		_ = c.Draw(dateLine)
	})
}

func getFontMap() (f fontMap, err error) {
	err = nil
	f = make(fontMap)
	f["Regular"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Regular.ttf")
	if err == nil {
		f["Black"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Black.ttf")
		f["BlackItalic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-BlackItalic.ttf")
		f["Bold"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Bold.ttf")
		f["BoldItalic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-BoldItalic.ttf")
		f["Italic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Italic.ttf")
		f["Light"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Light.ttf")
		f["LightItalic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-LightItalic.ttf")
		f["Medium"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Medium.ttf")
		f["MediumItalic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-MediumItalic.ttf")
		f["Regular"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Regular.ttf")
		f["Thin"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-Thin.ttf")
		f["ThinItalic"], err = model.NewPdfFontFromTTFFile("./fonts/Roboto-ThinItalic.ttf")
	}

	return f, err
}

func writePerformanceAnalysis(
	c *creator.Creator,
	f fontMap,
	received coachingList,
	given coachingList,
	topicToValueTypeMapping map[string]string,
	quarter int,
	year int,
	dialogDao fetch_dialog.DAO,
	logger logger.AdaptiveLogger,
) (tags map[string]string) {
	var analyses map[string]string
	analyses, tags = generateSummaryAnalysis(received, given, topicToValueTypeMapping, quarter, year, dialogDao, logger)
	logger.WithField("analyses", &analyses).Infof("Retrieved analyses")
	sortOrder := []string{
		"StrongRed",
		"WeakRed",
		"StrongYellow",
		"WeakYellow",
		"WeakGreen",
		"StrongGreen",
		"Neutral",
	}
	sortedAnalysis := make([]string, 0)
	sortedTags := make([]string, 0)

	ch := newChapter(
		c,
		false,
		"Feedback Analysis",
		f["Regular"],
		creator.ColorBlack,
		fontSizeHeadingOne,
	)

	for _, s := range sortOrder {
		for kind, analyis := range analyses {
			if tags[kind] == s {
				sortedAnalysis = append(sortedAnalysis, analyis)
				sortedTags = append(sortedTags, kind)
			}
		}
	}
	options := loadDialogUnsafe(dialogDao, CoachingIntro, "summary-explanation")
	logger.WithField("options", &options).Infof("Dialog entries for summary-explanation")

	_ = ch.Add(newParagraph(
		c,
		false,
		core_utils_go.RandomString(options.Dialog),
		f["Regular"],
		creator.ColorBlack,
		creator.TextAlignmentLeft,
		fontSizeNormal,
		0,
		0,
		padding,
		0,
	))

	ideaTable := c.NewTable(3)
	ideaTable.SetMargins(0, 0, padding, 0)
	_ = ideaTable.SetColumnWidths(0.10, 0.20, 0.70)

	analyticColumns := []string{"", "Topic", "Analysis"}
	for i, analyticColumn := range analyticColumns {
		cellBackground := creator.ColorWhite
		if i == 0 {
			cellBackground = creator.ColorBlack
		}
		newCell(
			ideaTable,
			creator.ColorBlack,
			cellBackground,
			creator.CellBorderSideAll,
			creator.CellBorderStyleSingle,
			borderWidth,
			creator.CellHorizontalAlignmentCenter,
			creator.CellVerticalAlignmentMiddle,
			newParagraph(
				c,
				false,
				analyticColumn,
				f["Bold"],
				creator.ColorBlack,
				creator.TextAlignmentLeft,
				fontSizeNormal,
				0,
				0,
				0,
				0,
			),
		)
	}

	for i := range sortedAnalysis {
		analytics := []string{"", sortedTags[i], sortedAnalysis[i]}
		for j, analytic := range analytics {
			cellBackground := creator.ColorWhite
			if j == 0 {
				cellBackground = getColor(tags[sortedTags[i]])
			}
			newCell(
				ideaTable,
				creator.ColorBlack,
				cellBackground,
				creator.CellBorderSideAll,
				creator.CellBorderStyleSingle,
				borderWidth,
				creator.CellHorizontalAlignmentLeft,
				creator.CellVerticalAlignmentMiddle,
				newParagraph(
					c,
					false,
					analytic,
					f["Regular"],
					creator.ColorBlack,
					creator.TextAlignmentLeft,
					fontSizeNormal,
					0,
					0,
					0,
					0,
				),
			)
		}
	}

	_ = ch.Add(ideaTable)
	_ = c.Draw(ch)

	return
}

func writePerformanceSummary(c *creator.Creator, f fontMap, received coachingList) {
	summary, err := nlp.GetSummaryText(summaryBig, received.createTextBlob(), nlp.English)

	if err != nil {
		log.Println("Error retrieving performance summary ", err)
	}

	ch := newChapter(
		c,
		false,
		"Feedback Summary",
		f["Regular"],
		creator.ColorBlack,
		fontSizeHeadingOne,
	)

	_ = ch.Add(newParagraph(
		c,
		false,
		summary,
		f["Regular"],
		creator.ColorBlack,
		creator.TextAlignmentLeft,
		fontSizeNormal,
		0,
		0,
		padding,
		0,
	))
	_ = c.Draw(ch)
}

func writeCoachingIdeas(
	c *creator.Creator,
	f fontMap,
	received coachingList,
	dialogDao fetch_dialog.DAO,
) {
	unfilteredCoachingSuggestions, err := nlp.GetTextCategoriesText(received.createTextBlob(), nlp.English)

	// First sort by Relevance
	sort.Slice(unfilteredCoachingSuggestions, func(i, j int) bool {
		return unfilteredCoachingSuggestions[i].GetRelevance() > unfilteredCoachingSuggestions[j].GetRelevance()
	})

	coachingDialogContext := map[string]string{
		"Adaptability":                  "adaptability",
		"Autonomy":                      "autonomy",
		"Collaboration and cooperation": "collaboration",
		"Commitment":                    "commitment",
		"Communication":                 "communication",
		"Customer relations":            "customer-relations",
		"Efficiency":                    "efficiency",
		"Expertise and skills":          "expertise-and-skills",
		"Learning capacity":             "innovation",
		"Training needs":                "leadership",
		"Innovation":                    "learning-capacity",
		"Leadership":                    "motivation",
		"Motivation":                    "personality",
		"Personality":                   "reliability",
		"Work environment":              "stress-management",
		"Reliability":                   "training-needs",
		"Stress management":             "work-environment",
	}

	coachingDialogSubject := map[string]string{
		"very positive": "very-positive",
		"positive":      "positive",
		"negative":      "negative",
		"very negative": "very-negative",
	}

	// Now filter out the stuff with a low relevance score
	// That has a sentiment we can work with
	// and is a topic  we've written for.
	coachingSuggestions := unfilteredCoachingSuggestions[:0]
	for _, x := range unfilteredCoachingSuggestions {
		if x.GetRelevance() > coachingTopicTolerance &&
			coachingDialogContext[x.GetLabel()] != "" &&
			coachingDialogSubject[x.GetSentiment()] != "" {
			coachingSuggestions = append(coachingSuggestions, x)
		}
	}

	if err != nil {
		log.Println("Error retrieving coaching ideas ", err)
	}

	ch := newChapter(
		c,
		false,
		"Coaching Ideas",
		f["Regular"],
		creator.ColorBlack,
		fontSizeHeadingOne,
	)

	var coachingIntro = coachingIdeaAnalysis(dialogDao)
	if len(coachingSuggestions) > 0 {
		coachingIntro = coachingIntro + "\n\nThese topics are sorted by how relevant they are to your feedback. The colors to the far left hand side are an indication of the sentiment found in your feedback from your colleagues about the given topic. Dark green means very positive, light green means positive, grey means neutral, light red means negative, and dark read means very negative."
		_ = ch.Add(newParagraph(
			c,
			false,
			coachingIntro,
			f["Regular"],
			creator.ColorBlack,
			creator.TextAlignmentLeft,
			fontSizeNormal,
			0,
			0,
			padding,
			0,
		))

		ideaTable := c.NewTable(3)
		ideaTable.SetMargins(0, 0, padding, 0)
		_ = ideaTable.SetColumnWidths(0.10, 0.20, 0.70)

		ideaCols := []string{"", "Topic", "Idea"}
		for i, idea := range ideaCols {
			cellBackground := creator.ColorWhite
			if i == 0 {
				cellBackground = creator.ColorBlack
			}
			newCell(
				ideaTable,
				creator.ColorBlack,
				cellBackground,
				creator.CellBorderSideAll,
				creator.CellBorderStyleSingle,
				borderWidth,
				creator.CellHorizontalAlignmentCenter,
				creator.CellVerticalAlignmentMiddle,
				newParagraph(
					c,
					false,
					idea,
					f["Bold"],
					creator.ColorBlack,
					creator.TextAlignmentLeft,
					fontSizeNormal,
					0,
					0,
					0,
					0,
				),
			)
		}

		sentimentColor := map[string]string{
			"very positive": StrongGreen,
			"positive":      WeakGreen,
			"neutral":       Neutral,
			"none":          Neutral,
			"negative":      WeakRed,
			"very negative": StrongRed,
		}
		for _, each := range coachingSuggestions {
			options := loadDialogUnsafe(dialogDao, coachingDialogContext[each.GetLabel()], coachingDialogSubject[each.GetSentiment()])

			dialog := core_utils_go.RandomString(options.Dialog)
			ideas := []string{"", each.GetLabel(), dialog}
			for i, idea := range ideas {
				backgroundColor := creator.ColorWhite
				alignment := creator.TextAlignmentLeft
				if i == 0 {
					backgroundColor = getColor(sentimentColor[each.GetSentiment()])
					alignment = creator.TextAlignmentCenter
				}
				newCell(
					ideaTable,
					creator.ColorBlack,
					backgroundColor,
					creator.CellBorderSideAll,
					creator.CellBorderStyleSingle,
					borderWidth,
					creator.CellHorizontalAlignmentLeft,
					creator.CellVerticalAlignmentMiddle,
					newParagraph(
						c,
						false,
						idea,
						f["Regular"],
						creator.ColorBlack,
						alignment,
						fontSizeNormal,
						0,
						0,
						0,
						0,
					),
				)
			}
		}

		_ = ch.Add(ideaTable)
		_ = c.Draw(ch)
	}
}

func writeFeedbackSummary(c *creator.Creator, f fontMap, received coachingList, kind string, topicToValueTypeMapping map[string]string) *creator.Chapter {
	if received.length() > 0 {
		feedbackSummary, err := nlp.GetSummaryText(summarySmall,
			received.kindCoaching(kind, topicToValueTypeMapping).createTextBlob(), nlp.English)

		if err != nil {
			log.Println("Error retrieving feedback summary ", err)
		}

		c.NewPage()

		ch := newChapter(
			c,
			false,
			strings.Title(kind)+" Feedback Summary",
			f["Regular"],
			creator.ColorBlack,
			fontSizeHeadingOne,
		)

		_ = ch.Add(
			newParagraph(
				c,
				false,
				feedbackSummary,
				f["Regular"],
				creator.ColorBlack,
				creator.TextAlignmentLeft,
				fontSizeNormal,
				0,
				0,
				padding,
				0,
			))

		_ = c.Draw(ch)
		return ch
	} else {
		return nil
	}
}

func writeTopic(c *creator.Creator, sc *creator.Chapter, f fontMap, topic string, received coachingList) {
	Topiccoaching := received.topicCoaching(topic)
	sc.SetMargins(indentOne, 0, padding, 0)
	writeTopicSummary(c, sc, f, topic, Topiccoaching)
	writeTopicDetails(c, sc, f, Topiccoaching)
}

func writeTopicSummary(c *creator.Creator, sc *creator.Chapter, f fontMap, topic string, tc coachingList) {
	if tc.length() > 0 {
		topicSummary, err := nlp.GetSummaryText(summarySmall, tc.topicCoaching(topic).createTextBlob(), nlp.English)

		if err != nil {
			log.Println("Error retrieving topic summary ", err)
		}

		sc.GetHeading().SetFont(f["Regular"])
		sc.GetHeading().SetFontSize(fontSizeHeadingTwo)
		sc.GetHeading().SetColor(creator.ColorBlack)

		_ = sc.Add(newParagraph(
			c,
			false,
			topicSummary,
			f["Regular"],
			creator.ColorBlack,
			creator.TextAlignmentLeft,
			fontSizeNormal,
			indentTwo,
			0,
			padding,
			0,
		))
	}
}

func writeTopicDetails(c *creator.Creator, sc *creator.Chapter, f fontMap, tc coachingList) {

	_ = sc.Add(newParagraph(
		c,
		false,
		"Details",
		f["Bold"],
		creator.ColorBlack,
		creator.TextAlignmentLeft,
		fontSizeHeadingThree,
		indentThree,
		0,
		padding,
		0,
	))

	details := tc.justFeedback()

	if details.length() == 0 {
		_ = sc.Add(newParagraph(
			c,
			false,
			"No coaching feedback for this topic.",
			f["Regular"],
			creator.ColorBlack,
			creator.TextAlignmentLeft,
			fontSizeNormal,
			indentFour,
			0,
			padding,
			0,
		))
	} else if details.length() == 1 {
		_ = sc.Add(newParagraph(
			c,
			false,
			details.index(0).GetComments(),
			f["Regular"],
			creator.ColorBlack,
			creator.TextAlignmentLeft,
			fontSizeNormal,
			indentFour,
			0,
			padding,
			0,
		))
	} else {
		detailsTable := c.NewTable(1)
		detailsTable.SetMargins(indentFour, indentFour, padding, 0)
		lineColor := creator.ColorRGBFromHex("#D5D8DC")
		for i := 0; i < details.length(); i++ {
			newCell(
				detailsTable,
				lineColor,
				creator.ColorWhite,
				creator.CellBorderSideTop,
				creator.CellBorderStyleSingle,
				borderWidth,
				creator.CellHorizontalAlignmentLeft,
				creator.CellVerticalAlignmentMiddle,
				newParagraph(
					c,
					false,
					details.index(i).GetComments(),
					f["Regular"],
					creator.ColorBlack,
					creator.TextAlignmentLeft,
					fontSizeNormal,
					0,
					0,
					0,
					padding,
				),
			)
		}

		cell := detailsTable.NewCell()
		cell.SetBackgroundColor(creator.ColorWhite)
		cell.SetBorder(creator.CellBorderSideTop, creator.CellBorderStyleSingle, borderWidth)
		cell.SetBorderColor(lineColor)
		_ = sc.Add(detailsTable)
	}
}

func newParagraph(
	c *creator.Creator,
	draw bool,
	text string,
	font *model.PdfFont,
	color creator.Color,
	alignment creator.TextAlignment,
	fontSize float64,
	leftMargin float64,
	rightMargin float64,
	topMargin float64,
	bottomMargin float64,
) *creator.Paragraph {
	p := c.NewParagraph(text)
	p.SetFont(font)
	p.SetFontSize(fontSize)
	p.SetColor(color)
	p.SetMargins(
		leftMargin,
		rightMargin,
		topMargin,
		bottomMargin,
	)
	p.SetTextAlignment(alignment)
	if draw {
		_ = c.Draw(p)
	}
	return p
}

func newChapter(
	c *creator.Creator,
	draw bool,
	text string,
	font *model.PdfFont,
	color creator.Color,
	fontSize float64,
) *creator.Chapter {
	ch := c.NewChapter(text)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(fontSize)
	ch.GetHeading().SetColor(color)
	if draw {
		_ = c.Draw(ch)
	}
	return ch
}

func newCell(
	table *creator.Table,
	borderColor creator.Color,
	backgroundColor creator.Color,
	borderSides creator.CellBorderSide,
	borderStyle creator.CellBorderStyle,
	borderWidth float64,
	horizontalAlignment creator.CellHorizontalAlignment,
	verticalAlignment creator.CellVerticalAlignment,
	drawable creator.VectorDrawable,
) {
	cell := table.NewCell()
	cell.SetBorderColor(borderColor)
	cell.SetBackgroundColor(backgroundColor)
	cell.SetBorder(borderSides, borderStyle, borderWidth)
	cell.SetHorizontalAlignment(horizontalAlignment)
	cell.SetVerticalAlignment(verticalAlignment)
	_ = cell.SetContent(drawable)
}
