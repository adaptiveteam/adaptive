package common

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)

type DynamoNamespace struct {
	Dynamo    *awsutils.DynamoRequest `json:"dynamo"`
	Namespace string                  `json:"namespace"`
}

type AttachmentEntity struct {
	MC       models.MessageCallback
	Title    ui.RichText
	Text     ui.RichText
	Fallback ui.PlainText
	Actions  []ebm.AttachmentAction
	Fields   []models.KvPair
}

// UserIDsToDisplayNames is a function that fetches user information for ids.
// For each user their display name is put to `Name` and user id to `Value`.
type UserIDsToDisplayNames func([]string) []models.KvPair

func KvPairFilterInplace(kvPairs []models.KvPair, predicate func(models.KvPair) bool) []models.KvPair {
	n := 0
	for _, x := range kvPairs {
		if predicate(x) {
			kvPairs[n] = x
			n++
		}
	}
	kvPairs = kvPairs[:n]
	return kvPairs
}

// NotKey is a predicate that returns true if the key is not equal to the given one
func NotKey(unwantedKey string) func(models.KvPair) bool {
	return func(kvPair models.KvPair) bool { return kvPair.Key != unwantedKey }
}

func TaggedUser(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}
