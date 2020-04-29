package holidays

import (
	"github.com/adaptiveteam/adaptive/daos/adHocHoliday"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"testing"
)

func TestDao(t *testing.T) {
	if false { // TODO: enable holidays test
		namespace := "test"
		d := awsutils.NewDynamo("us-east-1", "localhost:4570", namespace)

		adHocHoliday.CreateUnsafe(models.AdHocHoliday{
			Name:       "holiday1",
			PlatformID: "ivan",
		})
		conn := daosCommon.DynamoDBConnection{
			Dynamo: d,
			ClientID: "test",
			PlatformID: "test",
		}
		holidays := AllUnsafe(conn)
		for _, h := range holidays {
			if h.Name == "holiday1" {
				return
			}
		}
		t.Errorf("Holiday not found")
	}
}
