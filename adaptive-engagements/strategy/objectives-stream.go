package strategy

import (
	"github.com/adaptiveteam/adaptive/pagination"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)


// IsStrategyObjectiveNotCompleted a predicate
func IsStrategyObjectiveNotCompleted(so models.StrategyObjective) func(conn daosCommon.DynamoDBConnection) (res bool, err error) {
	return func(conn daosCommon.DynamoDBConnection) (res bool, err error) {
		var uos []userObjective.UserObjective
		uos, err = userObjective.ReadOrEmpty(so.ID)(conn)
		res = len(uos) > 0 && uos[0].Completed == 0
		return
	}
}

// AsStrategyObjectiveSlice casts each element to StrategyObjective
func AsStrategyObjectiveSlice(is pagination.InterfaceSlice) (res []strategyObjective.StrategyObjective) {
	for _, i := range is {
		res = append(res, i.(strategyObjective.StrategyObjective))
	}
	return
}
// AsInterfaceSlice casts each element to interface {}
func AsInterfaceSlice(sos []strategyObjective.StrategyObjective) (res pagination.InterfaceSlice) {
	for _, so := range sos {
		res = append(res, so)
	}
	return

}
// DynamoDBKey - key
type DynamoDBKey = map[string]*dynamodb.AttributeValue
// KeyExtractor obtains key value from the given instance
type KeyExtractor = func(interface{}) DynamoDBKey

// StrategyObjectiveKey retrieves key of the instance
func StrategyObjectiveKey(i interface{}) DynamoDBKey {
	so := i.(strategyObjective.StrategyObjective)
	return map[string]*dynamodb.AttributeValue{ 
		string(strategyObjective.PlatformID): daosCommon.DynS(string(so.PlatformID)),
		string(strategyObjective.ID):         daosCommon.DynS(so.ID),
	}
}
// InterfaceQueryPager constructs a pager that will be yielding pages using 
// DynamoDB Query.
// sliceOfEntities = []models.StrategyObjective
// asInterfaceSlice = AsInterfaceSlice
// keyExtractor = StrategyObjectiveKey
// TODO: combine these into a type class
func InterfaceQueryPager(
	conn daosCommon.DynamoDBConnection, 
	queryInput dynamodb.QueryInput,
	sliceOfEntities interface{},
	asInterfaceSlice func (interface{}) pagination.InterfaceSlice,
	keyExtractor KeyExtractor,
	) pagination.InterfacePager {
	return func() (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
		if queryInput.Limit == nil || *queryInput.Limit > 0 {
			// var instances []models.StrategyObjective
			err = conn.Dynamo.QueryInternal(queryInput, sliceOfEntities)
			if err == nil {
				sl = asInterfaceSlice(sliceOfEntities)
				if len(sl) > 0 {
					last := sl[len(sl) - 1]
					queryInput.ExclusiveStartKey = keyExtractor(last)
				}
				if queryInput.Limit != nil {
					if int64(len(sl)) > *queryInput.Limit {
						sl = sl[0:*queryInput.Limit]
						*queryInput.Limit = 0
					} else {
						*queryInput.Limit = *queryInput.Limit - int64(len(sl))
					}
				}
				
			}
		} 
		if (queryInput.Limit == nil || *queryInput.Limit > 0) && err == nil {
			ip = readByPlatformIDPager(conn, queryInput)
		} else {
			ip = pagination.InterfacePagerPure()
		}
		return
	}
}

func readByPlatformIDPager(conn daosCommon.DynamoDBConnection, queryInput dynamodb.QueryInput) func()  (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
	var instances []models.StrategyObjective
	return InterfaceQueryPager(conn, queryInput,
		instances,
		func (i interface{}) pagination.InterfaceSlice {return  AsInterfaceSlice(i.([]models.StrategyObjective))},
		StrategyObjectiveKey,
	)
}
// StrategyObjective_ReadByPlatformIDStream create stream that will 
// read all strategy objectives for platform id
func StrategyObjective_ReadByPlatformIDStream() (stm daosCommon.InterfaceStream) {
	platformIDIndex := string(strategyObjective.PlatformIDIndex)
	stm = func (conn daosCommon.DynamoDBConnection) pagination.InterfacePager {
		queryInput:= dynamodb.QueryInput{
			TableName:                 aws.String(strategyObjective.TableName(conn.ClientID)),
			KeyConditionExpression:    aws.String(string(strategyObjective.PlatformID + " = :a0")),
			IndexName:                 &platformIDIndex,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":a0": daosCommon.DynS(string(conn.PlatformID)),
			},
			// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html
			ScanIndexForward:          aws.Bool(true),
		}
		return readByPlatformIDPager(conn, queryInput)
	}
	return	
}

// SelectFromStrategyObjectiveJoinUserObjectiveWhereNotCompletedStream - stream-based SQL:
// SELECT * FROM _strategy_objective
// JOIN _user_objective ON _strategy_objective.id = _user_objective.id
// WHERE _user_objective.completed=0
func SelectFromStrategyObjectiveJoinUserObjectiveWhereNotCompletedStream() (stm daosCommon.InterfaceStream) {
	return StrategyObjective_ReadByPlatformIDStream().
	FilterF(
		func (conn daosCommon.DynamoDBConnection) func (i interface{}) (bool, error) {
		return func (i interface{}) (res bool, err error) {
			so := i.(strategyObjective.StrategyObjective)
			return IsStrategyObjectiveNotCompleted(so)(conn)
		}
	})
}
