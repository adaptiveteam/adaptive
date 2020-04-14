package collaboration_report

import (
	"github.com/adaptiveteam/adaptive/daos/adaptiveValue"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"encoding/json"
	utils "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/unidoc/unipdf/v3/model"
	"math"
	"sort"
)

func getRatingLanguage(value float64) (rating string) {
	const (
		upper = "upper"
		lower = "lower"
	)
	values := map[string]map[string]float64{
		"exceeded expectations": {
			upper: math.MaxFloat64,
			lower: 5.0,
		},
		"almost exceeded expectations": {
			upper: 5.0,
			lower: 4.5,
		},
		"is above expectations": {
			upper: 4.5,
			lower: 4.0,
		},
		"is almost above expectations": {
			upper: 4.0,
			lower: 3.5,
		},
		"meets expectations": {
			upper: 3.5,
			lower: 3.0,
		},
		"almost meets expectations": {
			upper: 3.0,
			lower: 2.5,
		},
		"is approaching expectations": {
			upper: 2.5,
			lower: 2.0,
		},
		"is almost approaching expectations": {
			upper: 2.0,
			lower: 1.5,
		},
		"did not meet expectations": {
			upper: 1.5,
			lower: -math.MaxFloat64,
		},
	}

	for k := range values {
		if value >= values[k][lower] && value < values[k][upper] {
			rating = k
		}
	}

	return rating
}

func (c CoachingList) getSortedAttribute(attribute func(Coaching) string) (sortedTopics []string) {

	// Get all topics
	topics := make([]string, 0)
	for i := 0; i < len(c); i++ {
		topics = append(topics, attribute(c[i]))
	}

	// Remove duplicates
	sortedTopics = utils.Distinct(topics)

	// Sort it all
	sort.Strings(sortedTopics)

	return sortedTopics
}

func (c CoachingList) getTopicToValueTypeMapping() (topicToValueTypeMapping map[string]string) {
	topicToValueTypeMapping = make(map[string]string, 0)
	for _, each := range c {
		topicToValueTypeMapping[each.Topic] = each.Type
	}

	return topicToValueTypeMapping
}

func (c CoachingList) length() int {
	return len(c)
}

func (c CoachingList) index(i int) Coaching {
	return c[i]
}

type GetCompetencyUnsafe = func (competencyID string) []adaptiveValue.AdaptiveValue

func GetCompetencyImpl(conn daosCommon.DynamoDBConnection) GetCompetencyUnsafe {
	return func (competencyID string) (competencies []adaptiveValue.AdaptiveValue) {
		competencies = adaptiveValue.ReadOrEmptyUnsafe(competencyID)(conn)
		competencies = adaptiveValue.AdaptiveValueFilterActive(competencies)
		return
	}
}

// ResolveCompetencies replaces value id in `Topic` with value name and value type in
// `Topic`, `Type`.
func ResolveCompetencies(coachingList CoachingList, getCompetencyUnsafe GetCompetencyUnsafe) (feedbackValueMappedList CoachingList) {
	for _, each := range coachingList {
		competencies := getCompetencyUnsafe(each.Topic)
		for _, competency := range competencies {
			feedbackMapped := Coaching{
				Source:   each.Source,
				Target:   each.Target,
				Topic:    competency.Name,
				Type:     competency.ValueType,
				Rating:   each.Rating,
				Comments: each.Comments,
				Quarter:  each.Quarter,
				Year:     each.Year,
			}
			feedbackValueMappedList = append(feedbackValueMappedList, feedbackMapped)
		}
	}
	return
}
// NewCoachingListFromStream - 
func NewCoachingListFromStream(stream []byte, getCompetencyUnsafe GetCompetencyUnsafe) (feedbackValueMappedList CoachingList, err error) {
	rv := make(CoachingList, 0)
	err = nil
	if len(stream) > 0 {
		err = json.Unmarshal(stream, &rv)
	}
	feedbackValueMappedList = ResolveCompetencies(rv, getCompetencyUnsafe)
	return feedbackValueMappedList, err
}

func (c CoachingList) justFeedback() CoachingList {
	rv := make(CoachingList, 0)
	for _, each := range c {
		if len(each.GetComments()) > 0 {
			rv = append(rv, each)
		}
	}
	return rv
}

func (c CoachingList) feedbackForQuarter(quarter int, year int) CoachingList {
	rv := make(CoachingList, 0)
	for _, each := range c {
		if each.GetQuarter() == quarter && each.GetYear() == year {
			rv = append(rv, each)
		}
	}
	return rv
}

func (c CoachingList) topicCoaching(topic string) CoachingList {
	rv := make(CoachingList, 0)
	for _, each := range c {
		if each.GetTopic() == topic {
			rv = append(rv, each)
		}
	}
	return rv
}

func (c CoachingList) typeCoaching(topicType string, topicToValueTypeMapping map[string]string) CoachingList {
	rv := make(CoachingList, 0)
	for _, each := range c {
		if topicToValueTypeMapping[each.GetTopic()] == topicType {
			rv = append(rv, each)
		}
	}
	return rv
}

func (c CoachingList) topics() (rv []string) {
	for _, each := range c {
		rv = append(rv, each.GetTopic())
	}

	return unique(rv)
}

func (c CoachingList) kindCoaching(kind string, topicToValueTypeMapping map[string]string) CoachingList {
	rv := make(CoachingList, 0)
	for _, each := range c {
		if topicToValueTypeMapping[each.GetTopic()] == kind {
			rv = append(rv, each)
		}
	}
	return rv
}

func (c CoachingList) justScores() (rv []float64) {
	rv = make([]float64, 0)
	for _, each := range c {
		rv = append(rv, each.GetRating())
	}
	return rv
}

func (c CoachingList) createTextBlob() string {
	var textBlob string
	for _, each := range c {
		textBlob = each.GetComments() + "\n" + textBlob
	}
	return textBlob
}

func (c CoachingList) calculateScore() float64 {
	scores := float64(0.0)
	length := float64(len(c))
	if length > 0 {
		for i := 0; i < len(c); i++ {
			scores += float64(c[i].GetRating())
		}
		return scores / length
	}
	return scores
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type fontMap map[string]*model.PdfFont

type CoachingList []Coaching

type Coaching struct {
	Source   string  `json:"source"`
	Target   string  `json:"target"`
	Topic    string  `json:"topic"`
	Type     string  `json:"type,omitempty"`
	Rating   float64 `json:"rating"`
	Comments string  `json:"comments"`
	Quarter  int     `json:"quarter"`
	Year     int     `json:"year"`
}

func (c Coaching) GetSource() string {
	return c.Source
}

func (c Coaching) GetTopic() string {
	return c.Topic
}

func (c Coaching) GetType() string {
	return c.Type
}

func (c Coaching) GetComments() string {
	return c.Comments
}

func (c Coaching) GetRating() float64 {
	return c.Rating
}

func (c Coaching) GetQuarter() int {
	return c.Quarter
}

func (c Coaching) GetYear() int {
	return c.Year
}

func (c Coaching) Set(
	source string,
	target string,
	topic string,
	comments string,
	rating float64,
	quarter int,
	year int,
) {
	c.Source = source
	c.Target = target
	c.Topic = topic
	c.Comments = comments
	c.Quarter = quarter
	c.Year = year
}
