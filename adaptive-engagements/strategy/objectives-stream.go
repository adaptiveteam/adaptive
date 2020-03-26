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

func readByPlatformIDPager(conn daosCommon.DynamoDBConnection, queryInput dynamodb.QueryInput) func()  (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
	return func() (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
		if queryInput.Limit == nil || *queryInput.Limit > 0 {
			var instances []models.StrategyObjective
			err = conn.Dynamo.QueryInternal(queryInput, &instances)
			if err == nil {
				sl = AsInterfaceSlice(instances)
				if len(sl) > 0 {
					last := instances[len(instances) - 1]
					queryInput.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{ 
						string(strategyObjective.PlatformID): daosCommon.DynS(string(last.PlatformID)),
						string(strategyObjective.ID): daosCommon.DynS(last.ID),
					}
				}
				if queryInput.Limit != nil {
					if int64(len(instances)) > *queryInput.Limit {
						sl = sl[0:*queryInput.Limit]
						*queryInput.Limit = 0
					} else {
						*queryInput.Limit = *queryInput.Limit - int64(len(instances))
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
			// var uos []userObjective.UserObjective
			// uos, err = userObjective.ReadOrEmpty(so.ID)(conn)
			// res = len(uos) > 0 && uos[0].Completed == 0 
			// return
		}
	})
}
