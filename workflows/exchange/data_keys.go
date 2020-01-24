package exchange

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
)

const IssueIDKey = "iid"
const IssueTypeKey = "itype"


const CommunityNamespace = "community"

var CommunityPath models.Path = models.ParsePath("/" + CommunityNamespace)
