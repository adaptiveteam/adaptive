package adaptive_utils_go

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
)

func Coachees(table, index string, targetQY models.TargetQY, d *awsutils.DynamoRequest) ([]models.CoachingRelationship, error) {
	var rels []models.CoachingRelationship
	err := d.QueryTableWithIndex(table, awsutils.DynamoIndexExpression{
		IndexName: index,
		Condition: "coach_quarter_year = :cqy",
		Attributes: map[string]interface{}{
			":cqy": fmt.Sprintf("%s:%d:%d", targetQY.Target, targetQY.Quarter, targetQY.Year),
		},
	}, map[string]string{}, true, -1, &rels)
	return rels, err
}
