package strategy

import (
	// "log"
	"github.com/adaptiveteam/adaptive/core-utils-go"
	"github.com/adaptiveteam/adaptive/daos/adaptiveCommunityUser"
	"github.com/adaptiveteam/adaptive/adaptive-engagements/community"
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
// StrategyObjective_AsInterfaceSlice casts each element to interface {}
func StrategyObjective_AsInterfaceSlice(sos []strategyObjective.StrategyObjective) (res pagination.InterfaceSlice) {
	for _, so := range sos {
		res = append(res, so)
	}
	return
}
// AdaptiveCommunityUser_AsInterfaceSlice casts each element to interface {}
func AdaptiveCommunityUser_AsInterfaceSlice(sos []adaptiveCommunityUser.AdaptiveCommunityUser) (res pagination.InterfaceSlice) {
	for _, so := range sos {
		res = append(res, so)
	}
	return
}

// StrategyObjectiveKey retrieves key of the instance
func StrategyObjectiveKey(i interface{}) daosCommon.DynamoDBKey {
	so := i.(strategyObjective.StrategyObjective)
	return map[string]*dynamodb.AttributeValue{ 
		string(strategyObjective.PlatformID): daosCommon.DynS(string(so.PlatformID)),
		string(strategyObjective.ID):         daosCommon.DynS(so.ID),
	}
}

// AdaptiveCommunityUserKey retrieves key of the instance
func AdaptiveCommunityUserKey(i interface{}) daosCommon.DynamoDBKey {
	acu := i.(adaptiveCommunityUser.AdaptiveCommunityUser)
	return map[string]*dynamodb.AttributeValue{ 
		string(adaptiveCommunityUser.ChannelID): daosCommon.DynS(string(acu.ChannelID)),
		string(adaptiveCommunityUser.UserID):    daosCommon.DynS(acu.UserID),
	}
}

// InterfaceQueryPager constructs a pager that will be yielding pages using 
// DynamoDB Query.
// sliceOfEntities = &[]models.StrategyObjective
// asInterfaceSlice = AsInterfaceSlice
// keyExtractor = StrategyObjectiveKey
// TODO: combine these into a type class
func InterfaceQueryPager(
	conn daosCommon.DynamoDBConnection, 
	queryInput dynamodb.QueryInput,
	queryable daosCommon.Queryable,
	) pagination.InterfacePager {
	return func() (sl pagination.InterfaceSlice, ip pagination.InterfacePager, err error) {
		if queryInput.Limit == nil || *queryInput.Limit > 0 {
			// var instances []models.StrategyObjective
			sliceOfEntities := queryable.PointerToSliceOfEntities()
			err = conn.Dynamo.QueryInternal(queryInput, sliceOfEntities)
			if err == nil {
				// log.Printf("sliceOfEntities=%v", sliceOfEntities)
				sl = queryable.AsInterfaceSlice(sliceOfEntities)
				if len(sl) > 0 {
					last := sl[len(sl) - 1]
					// log.Printf("last=%v", last)
					queryInput.ExclusiveStartKey = queryable.KeyExtractor(last)
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
			ip = InterfaceQueryPager(conn, queryInput, queryable)
		} else {
			ip = pagination.InterfacePagerPure()
		}
		return
	}
}

// StrategyObjective_QueryableImpl Queryable implementation for StrategyObjective
type StrategyObjective_QueryableImpl struct {}

// PointerToSliceOfEntities returns *[]T
func (StrategyObjective_QueryableImpl)PointerToSliceOfEntities() interface{} {
	var instances []models.StrategyObjective
	return &instances
}

// AsInterfaceSlice casts input to []T and then converts each entity to []interface{}
func (StrategyObjective_QueryableImpl)AsInterfaceSlice(i interface{}) pagination.InterfaceSlice {
	return StrategyObjective_AsInterfaceSlice(*i.(*[]models.StrategyObjective))
}

// KeyExtractor extracts Key from T
func (StrategyObjective_QueryableImpl)KeyExtractor(i interface{}) daosCommon.DynamoDBKey {
	return StrategyObjectiveKey(i)
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
		return InterfaceQueryPager(conn, queryInput, StrategyObjective_QueryableImpl{})
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

// SelectFromStrategyObjectiveJoinCommunityWhereUserIDOrInStrategyCommunityStream - stream-basedimplementation of the following SQL:
// SELECT * FROM _strategy_objective
// WHERE _strategy_objective.user_id=$userID 
//    OR IsUserInCommunity($userID, 'strategy')
func SelectFromStrategyObjectiveJoinCommunityWhereUserIDOrInStrategyCommunityStream(userID string) (stm daosCommon.InterfaceStream) {
	var isUserInStrategyCommunity daosCommon.InterfaceStream
	isUserInStrategyCommunity = SelectNonEmptyFromCommunityWhereUserIDCommunityIDStream(userID, community.Strategy)
	return isUserInStrategyCommunity.FlatMapF(
		func (b interface{}) (res daosCommon.InterfaceStream) {
			isUserInStrategyCommunity := b.(bool)
			if isUserInStrategyCommunity {
				res = StrategyObjective_ReadByPlatformIDStream()
			} else {
				res = SelectFromObjectivesJoinCommunityUsersWhereUserIDStream(userID)
			}
			return
		},
	)
}

// AdaptiveCommunityUser_QueryableImpl Queryable implementation for AdaptiveCommunityUser
type AdaptiveCommunityUser_QueryableImpl struct {}

// PointerToSliceOfEntities returns *[]T
func (AdaptiveCommunityUser_QueryableImpl)PointerToSliceOfEntities() interface{} {
	var instances []adaptiveCommunityUser.AdaptiveCommunityUser
	return &instances
}

// AsInterfaceSlice casts input to []T and then converts each entity to []interface{}
func (AdaptiveCommunityUser_QueryableImpl)AsInterfaceSlice(i interface{}) pagination.InterfaceSlice {
	return AdaptiveCommunityUser_AsInterfaceSlice(*i.(*[]adaptiveCommunityUser.AdaptiveCommunityUser))
}

// KeyExtractor extracts Key from T
func (AdaptiveCommunityUser_QueryableImpl)KeyExtractor(i interface{}) daosCommon.DynamoDBKey {
	return AdaptiveCommunityUserKey(i)
}


func AdaptiveCommunityUser_ReadByUserIDCommunityIDStream(userID string, communityID community.AdaptiveCommunity) (stm daosCommon.InterfaceStream) {
	userIDCommunityIDIndex := string(adaptiveCommunityUser.UserIDCommunityIDIndex)
	stm = func (conn daosCommon.DynamoDBConnection) pagination.InterfacePager {
		queryInput:= dynamodb.QueryInput{
			TableName:                 aws.String(adaptiveCommunityUser.TableName(conn.ClientID)),
			KeyConditionExpression:    aws.String(string(adaptiveCommunityUser.UserID + " = :a0" + " and " + adaptiveCommunityUser.CommunityID + " = :a1")),
			IndexName:                 &userIDCommunityIDIndex,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":a0": daosCommon.DynS(string(userID)),
				":a1": daosCommon.DynS(string(communityID)),
			},
			// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html
			ScanIndexForward:          aws.Bool(true),
		}	
		return InterfaceQueryPager(conn, queryInput, AdaptiveCommunityUser_QueryableImpl{})
	}
	return	
}

func AdaptiveCommunityUser_ReadByUserIDStream(userID string) (stm daosCommon.InterfaceStream) {
	userIDIndex := string(adaptiveCommunityUser.UserIDIndex)
	stm = func (conn daosCommon.DynamoDBConnection) pagination.InterfacePager {
		queryInput:= dynamodb.QueryInput{
			TableName:                 aws.String(adaptiveCommunityUser.TableName(conn.ClientID)),
			KeyConditionExpression:    aws.String(string(adaptiveCommunityUser.UserID + " = :a0")),
			IndexName:                 &userIDIndex,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":a0": daosCommon.DynS(string(userID)),
			},
			// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html
			ScanIndexForward:          aws.Bool(true),
		}
		return InterfaceQueryPager(conn, queryInput, AdaptiveCommunityUser_QueryableImpl{})
	}
	return	
}
// SelectNonEmptyFromCommunityWhereUserIDCommunityIDStream - stream-based implementation of SQL:
// SELECT NonEmpty(*) FROM _adaptive_community_user
// WHERE _adaptive_community_user.user_id=$userID AND _adaptive_community_user.community=$community
// NB: the stream contains exactly 1 element of boolean type.
func SelectNonEmptyFromCommunityWhereUserIDCommunityIDStream(userID string, communityID community.AdaptiveCommunity) (stm daosCommon.InterfaceStream) {
	s := AdaptiveCommunityUser_ReadByUserIDCommunityIDStream(userID, communityID)
	return s.NonEmpty()
}
// SelectFromObjectivesWhereCommunitiesStream - stream-based implementation of the following SQL:
// SELECT _strategy_objective.* 
// FROM _strategy_objective 
// WHERE _strategy_objective.community_ids CONTAINS SOME OF communityIDs
func SelectFromObjectivesWhereCommunitiesStream(communityIDs []string)  (stm daosCommon.InterfaceStream) {
	hasIntersectionWithCommunityIDs := core_utils_go.IsIntersectionNonEmpty(communityIDs)
	return StrategyObjective_ReadByPlatformIDStream().
		FilterF(
			func (conn daosCommon.DynamoDBConnection) func (i interface{}) (bool, error) {
				return func (i interface{}) (res bool, err error) {
					so := i.(strategyObjective.StrategyObjective)
					return hasIntersectionWithCommunityIDs(so.CapabilityCommunityIDs), nil
				}
			},
		)
}
// SelectFromObjectivesJoinCommunityUsersWhereUserIDStream - stream-based implementation of SQL:
// SELECT _strategy_objective.* 
// FROM _strategy_objective JOIN _adaptive_community_user ON _strategy_objective.community_ids CONTAINS _adaptive_community_user.id
// WHERE _adaptive_community_user.user_id=$userID
func SelectFromObjectivesJoinCommunityUsersWhereUserIDStream(userID string)  (res daosCommon.InterfaceStream) {
	communityUsers := AdaptiveCommunityUser_ReadByUserIDStream(userID)
	
	return communityUsers.
	Map(func (i interface{}) interface{} { 
		acu := i.(adaptiveCommunityUser.AdaptiveCommunityUser)
		return acu.CommunityID
	}).
	All().
	FlatMapF(func (i interface{}) daosCommon.InterfaceStream {
		communityIDs := i.([]string)
		return SelectFromObjectivesWhereCommunitiesStream(communityIDs)
	})
}
