package common

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/adaptiveteam/adaptive/pagination"
)
// DynamoDBKey - key
type DynamoDBKey = map[string]*dynamodb.AttributeValue
// KeyExtractor obtains key value from the given instance
type KeyExtractor = func(interface{}) DynamoDBKey

// Queryable is a type class for an entity T 
type Queryable interface {
	// PointerToSliceOfEntities returns *[]T:
	// 	   var instances []T
	//     return &instances
	PointerToSliceOfEntities() interface{}
	// AsInterfaceSlice casts input to []T and then converts each entity to []interface{}
	// func _AsInterfaceSlice(i interface{}) (res pagination.InterfaceSlice) {
	//  ts := i.([]T)
	// 	for _, t := range ts {
	// 		res = append(res, t)
	// 	}
	// 	return
	// }
	AsInterfaceSlice(interface{}) pagination.InterfaceSlice
	// KeyExtractor extracts Key from T
	// func StrategyObjectiveKey(i interface{}) DynamoDBKey {
	// 	so := i.(strategyObjective.StrategyObjective)
	// 	return map[string]*dynamodb.AttributeValue{ 
	// 		string(strategyObjective.PlatformID): daosCommon.DynS(string(so.PlatformID)),
	// 		string(strategyObjective.ID):         daosCommon.DynS(so.ID),
	// 	}
	// }
	KeyExtractor(interface{}) DynamoDBKey
}
