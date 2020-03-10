package reporting_transformed_model_streaming_lambda

import (
	"github.com/adaptiveteam/adaptive/daos/common"
	"testing"

	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	"github.com/google/uuid"
	"gotest.tools/assert"
)

func TestUserTypeInference(t *testing.T) {
	teamID := models.ParseTeamID(common.PlatformID(uuid.New().String()))
	user := models.User{
		// UserProfile: models.UserProfile{},
		PlatformID:  teamID.ToPlatformID(),
		// PlatformOrg: "",
		// IsAdmin:     false,
		// // Deleted:     false,
		// CreatedAt:   "",
		// ModifiedAt:  "",
		// IsShared:    false,
	}
	var u interface{}
	u = user
	if w, ok := u.(models.User); ok {
		assert.Assert(t, w.PlatformID == teamID.ToPlatformID())
	} else {
		t.Fail()
	}
}
