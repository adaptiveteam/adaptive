package feedback

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/aws/aws-sdk-go/aws"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
)

// DAO is a CRUD wrapper around the _adaptive_user_feedback Dynamo DB table
// Requires access to the table and the index.
type DAO interface {
	Read(userID string, quarterYear string) (feedbacks []models.UserFeedback, err error)
	ReadUnsafe(userID string, quarterYear string) (feedbacks []models.UserFeedback)
	Create(user models.UserFeedback) error
	CreateUnsafe(user models.UserFeedback)

	IsThereFeedbackFromTo(from string, to string, quarterYear string) bool
}

// DAOImpl - a container for all information needed to access a DynamoDB table
type DAOImpl struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
	models.UserFeedbackTableSchema
}

// NewDAO creates an instance of DAO that will provide access to UserFeedback table
func NewDAO(dynamo *awsutils.DynamoRequest, namespace, table string, feedbackSourceQuarterYearIndex string) DAO {
	if table == "" { panic("Cannot create User DAO without table") }
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		UserFeedbackTableSchema: models.UserFeedbackTableSchema{Name: table, 
			FeedbackSourceQuarterYearIndex: feedbackSourceQuarterYearIndex},
	}
}

// NewDAOFromSchema creates an instance of DAO that will provide access to adaptiveValues table
func NewDAOFromSchema(dynamo *awsutils.DynamoRequest, namespace string, schema models.Schema) DAO {	
	return DAOImpl{Dynamo: dynamo, Namespace: namespace, 
		UserFeedbackTableSchema: schema.UserFeedback}
}

// Read reads User's feedback records for the quarter
func (d DAOImpl) Read(userID string, quarterYear string) (feedbacks []models.UserFeedback, err error) {
	err = d.Dynamo.QueryTableWithIndex(d.Name, awsutils.DynamoIndexExpression{
		IndexName: d.FeedbackSourceQuarterYearIndex,
		// there is no != operator for ConditionExpression
		Condition: "quarter_year = :qy AND #source = :s",
		Attributes: map[string]interface{}{
			":s":  userID,
			":qy": quarterYear,
		},
	}, map[string]string{"#source": "source"}, true, -1, &feedbacks)
	return
}
	
// ReadUnsafe reads the User. Panics in case of any errors
func (d DAOImpl) ReadUnsafe(userID string, quarterYear string) (feedbacks []models.UserFeedback) {
	out, err := d.Read(userID, quarterYear)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not query %s index", d.FeedbackSourceQuarterYearIndex))
	return out
}

// Create saves the User's feedback.
func (d DAOImpl) Create(userFeedback models.UserFeedback) error {
	return d.Dynamo.PutTableEntryWithCondition(userFeedback, d.Name, 
		"attribute_not_exists(id)")
}

// CreateUnsafe saves the User.
func (d DAOImpl) CreateUnsafe(user models.UserFeedback) {
	err := d.Create(user)
	core.ErrorHandler(err, d.Namespace, fmt.Sprintf("Could not create %s in %s", user.Id, d.Name))
}

func dynString(str string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{S: aws.String(str)}
}


// IsThereFeedbackFromTo searches database for feedback `from` `to`
func (d DAOImpl) IsThereFeedbackFromTo(from string, to string, quarterYear string) bool {
	feedbacks := d.ReadUnsafe(from, quarterYear)
	for _, f := range feedbacks {
		if f.Target == to {
			return true
		}
	}
	return false
}
