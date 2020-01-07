package holidays

import (
	"github.com/adaptiveteam/adaptive/adaptive-engagements/common"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	"testing"
)

func TestDao(t *testing.T) {
	namespace := "test"
	d := awsutils.NewDynamo("us-east-1", "localhost:4570", namespace)
	dns := common.DynamoNamespace{Dynamo: d, Namespace: namespace}
	daoHolidaysTable := NewDAO(&dns, "adHocHolidaysTable", "platformIndex")

	daoHolidaysTable.AddAdHocHoliday(models.AdHocHoliday{
		Name:       "holiday1",
		PlatformID: "ivan",
	})
	holidays := daoHolidaysTable.ForPlatformID("test").AllUnsafe()
	for _, h := range holidays {
		if h.Name == "holiday1" {
			return
		}
	}
	t.Errorf("Holiday not found")
}
