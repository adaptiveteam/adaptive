package adaptive_utils_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

func ParseToCallback(ip string) (*models.MessageCallback, error) {
	list := strings.Split(ip, ":")
	if len(list) != 7 {
		return nil, errors.New("validation error: incorrect number of arguments in input " + ip)
	}
	return &models.MessageCallback{Module: list[0], Source: list[1], Topic: list[2], Action: list[3], Target: list[4], Month: list[5], Year: list[6]}, nil
}

// ParseToCallbackValue returns MessageCallback value rather than address.
func ParseToCallbackValue(ip string) (models.MessageCallback, error) {
	list := strings.Split(ip, ":")
	if len(list) != 7 {
		return models.MessageCallback{}, errors.New("validation error: incorrect number of arguments in input")
	}
	return models.MessageCallback{Module: list[0], Source: list[1], Topic: list[2], Action: list[3], Target: list[4], Month: list[5], Year: list[6]}, nil
}

// MessageCallbackParseUnsafe parses the identifier and panics in case of error
// @param namespace - is used to report error
func MessageCallbackParseUnsafe(callbackID, namespace string) models.MessageCallback {
	mc, err := ParseToCallbackValue(callbackID)
	core.ErrorHandler(err, namespace, fmt.Sprintf("MessageCallback: Could not parse"))
	return mc
}

func AddEng(eng models.UserEngagement, table string, d *awsutils.DynamoRequest, namespace string) {
	err := d.PutTableEntry(eng, table)
	core.ErrorHandler(err, namespace,
		fmt.Sprintf("Could not write to %s table for a new user with user engagement: %v", table, eng))
}

func UpdateEngAsAnswered(source, callbackId, table string, d *awsutils.DynamoRequest, namespace string) {
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(source),
		},
		"id": {
			S: aws.String(callbackId),
		},
	}
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":a": {
			N: aws.String("1"),
		},
	}
	updateExpression := "set answered = :a"
	err := d.UpdateTableEntry(exprAttributes, key, updateExpression, table)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update entry in %s table", table))
}

func UpdateEngAsIgnored(source, callbackId, table string, d *awsutils.DynamoRequest, namespace string) {
	key := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(source),
		},
		"id": {
			S: aws.String(callbackId),
		},
	}
	exprAttributes := map[string]*dynamodb.AttributeValue{
		":i": {
			N: aws.String("1"),
		},
	}
	updateExpression := "set ignored = :i"
	err := d.UpdateTableEntry(exprAttributes, key, updateExpression, table)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not update entry in %s table", table))
}

func UserEng(id, userId, table string, d *awsutils.DynamoRequest, namespace string) ebm.Message {
	params := map[string]*dynamodb.AttributeValue{
		"user_id": {
			S: aws.String(userId),
		},
		"id": {
			S: aws.String(id),
		},
	}
	var ue models.UserEngagement
	err := d.QueryTable(table, params, &ue)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not query %s table", table))
	var op ebm.Message
	err = json.Unmarshal([]byte(ue.Script), &op)
	core.ErrorHandler(err, namespace, fmt.Sprintf("Could not unmarshal query output to Message from surveys"))
	return op
}
