package strategy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	daosCommon "github.com/adaptiveteam/adaptive/daos/common"
	"github.com/adaptiveteam/adaptive/daos/strategyObjective"
	"github.com/adaptiveteam/adaptive/daos/userObjective"
)

// StrategyObjectivePager allows lazy loading of StrategyObjective-s
type StrategyObjectivePager interface {
	// Next returns a portion of StrategyObjective-s and a pager instance that
	// is capable of returning further pages.
	// If the slice is empty, then there are no further StrategyObjective-s.
	// the returned StrategyObjectivePager will always yield empty slice.
	Next() ([]models.StrategyObjective, StrategyObjectivePager, error)
}
type StrategyObjectiveEffectfulPredicate = func(so models.StrategyObjective) func(conn daosCommon.DynamoDBConnection) (bool, error)

// StrategyObjectiveConst constructs predicate that always return true or false.
func StrategyObjectiveBoolConst(r bool) StrategyObjectiveEffectfulPredicate {
	return func(so models.StrategyObjective) func(conn daosCommon.DynamoDBConnection) (bool, error) {
		return func(conn daosCommon.DynamoDBConnection) (bool, error) {
			return r, nil
		}
	}
}
// StrategyObjectiveStream provides a couple of convenience methods to shrink
// the stream.
type StrategyObjectiveStream interface {
	Run(conn daosCommon.DynamoDBConnection) StrategyObjectivePager
	Limit(limit int) StrategyObjectiveStream
	Filter(pred StrategyObjectiveEffectfulPredicate) StrategyObjectiveStream
}

// IsStrategyObjectiveNotCompleted a predicate
func IsStrategyObjectiveNotCompleted(so models.StrategyObjective) func(conn daosCommon.DynamoDBConnection) (res bool, err error) {
	return func(conn daosCommon.DynamoDBConnection) (res bool, err error) {
		var uos []userObjective.UserObjective
		uos, err = userObjective.ReadOrEmpty(so.ID)(conn)
		res = len(uos) > 0 && uos[0].Completed == 0
		return
	}
}

type strategyObjectiveStreamImpl struct {
	dynamodb.QueryInput
	Filters []StrategyObjectiveEffectfulPredicate
}

func ReadByPlatformIDStream(platformID daosCommon.PlatformID) (stm StrategyObjectiveStream) {
	platformIDIndex := string(strategyObjective.PlatformIDIndex)
	stm = strategyObjectiveStreamImpl{
		QueryInput: dynamodb.QueryInput{
			KeyConditionExpression:    aws.String(string(strategyObjective.PlatformID + " = :a0")),
			IndexName:                 &platformIDIndex,
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":a0": daosCommon.DynS(string(platformID)),
			},
			// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html
			ScanIndexForward:          aws.Bool(true),
		},
		Filters: []StrategyObjectiveEffectfulPredicate{},
	}

	return	
}

func (s strategyObjectiveStreamImpl) Run(conn daosCommon.DynamoDBConnection) StrategyObjectivePager {
	return strategyObjectivePagerImpl{
		strategyObjectiveStreamImpl: s,
		DynamoDBConnection: conn,
	}
}

func (s strategyObjectiveStreamImpl) Limit(limit int) (res StrategyObjectiveStream) {
	qi := s.QueryInput
	var l int64
	l = int64(limit)
	qi.Limit = &l
	res = strategyObjectiveStreamImpl{
		QueryInput: qi,
		Filters: s.Filters,
	}
	return
}

func (s strategyObjectiveStreamImpl) Filter(pred StrategyObjectiveEffectfulPredicate) (res StrategyObjectiveStream) {
	return strategyObjectiveStreamImpl{
		QueryInput: s.QueryInput,
		Filters: append(s.Filters, pred),
	}
}

type strategyObjectivePagerImpl struct {
	strategyObjectiveStreamImpl
	daosCommon.DynamoDBConnection
}

func (s strategyObjectivePagerImpl) Next() (res []models.StrategyObjective, next StrategyObjectivePager, err error) {
	if s.QueryInput.Limit == nil || *s.QueryInput.Limit > 0 {
		qi := s.QueryInput
		conn := s.DynamoDBConnection
		qi.TableName = aws.String(strategyObjective.TableName(conn.ClientID))
		var instances []models.StrategyObjective
		err = conn.Dynamo.QueryInternal(qi, &instances)
		if err == nil {
			if len(instances) > 0 {
				last := instances[len(instances) - 1]
				qi.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{ 
					string(strategyObjective.PlatformID): daosCommon.DynS(string(last.PlatformID)),
					string(strategyObjective.ID): daosCommon.DynS(last.ID),
				}
				nextImpl := strategyObjectivePagerImpl{
					strategyObjectiveStreamImpl: strategyObjectiveStreamImpl{
						QueryInput: qi,
						Filters: s.Filters,
					},
					DynamoDBConnection: conn,
				}
				for _, i := range instances{
					for _, f := range s.Filters {
						var fl bool
						fl, err = f(i)(conn)
						if err != nil {
							return
						}
						if fl {
							res = append(res, i)
						}
					}
				}
				if len(res) == 0 {
					return next.Next() // NB! Recursion probably without tailcall optimization.
				} else if s.QueryInput.Limit != nil {
					if int64(len(res)) > *s.QueryInput.Limit {
						res = res[0:*s.QueryInput.Limit]
						*nextImpl.strategyObjectiveStreamImpl.QueryInput.Limit = 0
					} else {
						*nextImpl.strategyObjectiveStreamImpl.QueryInput.Limit = *nextImpl.QueryInput.Limit - int64(len(res))
					}
				}
				next = nextImpl
			} else {
				next = s // res is empty; err is nil; subsequent calls should yield the same
			}
		}
	} else {
		next = s // res is empty; err is nil; subsequent calls should yield the same
	} 
	return
}

func SelectFromStrategyObjectiveJoinUserObjectiveWhereNotCompletedStream() func(conn daosCommon.DynamoDBConnection) (StrategyObjectiveStream, err error) {
	return func(conn daosCommon.DynamoDBConnection) (StrategyObjectiveStream, err error) {
		var allObjs []models.StrategyObjective
		allObjs, err = strategyObjective.ReadByPlatformID(conn.PlatformID)(conn)
		if err == nil {
			for _, so := range allObjs {
				var uos []userObjective.UserObjective
				uos, err = userObjective.ReadOrEmpty(so.ID)(conn)
				if len(uos) > 0 && uos[0].Completed == 0 {
					// objs = append(objs, so)
				}
			}
		}
		return
	}
}

