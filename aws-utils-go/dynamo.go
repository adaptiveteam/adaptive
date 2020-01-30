package aws_utils_go

import (
	"encoding/json"
	"fmt"
	"github.com/adaptiveteam/adaptive/core-utils-go/logger"
	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"strings"
)

type DynamoTable struct {
	Name string  `json:"name"`
	Arn  *string `json:"arn"`
}

type DynamoIndexExpression struct {
	IndexName  string                 `json:"index_name"`
	Condition  string                 `json:"condition"`
	Attributes map[string]interface{} `json:"attributes"`
}

type DynamoRequest struct {
	svc dynamodbiface.DynamoDBAPI
	log *logger.Logger
}

func NewDynamo(region, endpoint, namespace string) *DynamoRequest {
	session, config := sess(region, endpoint)
	return &DynamoRequest{
		svc: dynamodb.New(session, config),
		log: logger.WithNamespace(fmt.Sprintf("adaptive.dynamo.%s", namespace)),
	}
}

func NewDax(region, endpoint, namespace string) (*DynamoRequest, error) {
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{endpoint}
	cfg.Region = region
	d, err2 := dax.New(cfg)
	return &DynamoRequest{
		svc: d,
		log: logger.WithNamespace(fmt.Sprintf("adaptive.dax.%s", namespace)),
	}, err2
}

func (d *DynamoRequest) errorLog(err error) {
	d.log.Errorf(err.Error())
}

func (d *DynamoRequest) ListTables(input *dynamodb.ListTablesInput) (*dynamodb.ListTablesOutput, error) {
	return d.svc.ListTables(input)
}

func (d *DynamoRequest) CreateTable(input *dynamodb.CreateTableInput) (err error) {
	_, err = d.svc.CreateTable(input)
	return
}

func (d *DynamoRequest) DescribeTable(name string) (*DynamoTable, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	}
	print(input, true)
	op, err2 := d.svc.DescribeTable(input)
	if err2 != nil {
		d.errorLog(err2)
		return nil, err2
	}
	return &DynamoTable{Name: name, Arn: op.Table.TableArn}, nil
}
// PutTableEntry -
// deprecated. The order of arguments in this function is not consistent with the other functions.
// more important table name should go first.
func (d *DynamoRequest) PutTableEntry(item interface{}, table string) (err error) {
	return d.PutItemIntoTable(table, item)
}

// PutItemIntoTable puts a new entry to the table
// it has consistent order of arguments.
func (d *DynamoRequest) PutItemIntoTable(tableName string, item interface{}) (err error) {
	if tableName == "" {
		err = fmt.Errorf("Table name is empty")
	} else {
		av, err2 := dynamodbattribute.MarshalMap(item)
		err = err2
		if err == nil {
			input := &dynamodb.PutItemInput{
				Item:      av,
				TableName: aws.String(tableName),
			}
			// var o *dynamodb.PutItemOutput
			_, err = d.svc.PutItem(input)
			// o.
		}
	}
	if err != nil {
		d.errorLog(err)
	}
	return
}

func (d *DynamoRequest) PutTableEntryWithCondition(item interface{}, table string, conditional string) (err error) {
	av, err2 := dynamodbattribute.MarshalMap(item)
	err = err2
	if err == nil {
		input := &dynamodb.PutItemInput{
			Item:                av,
			TableName:           aws.String(table),
			ConditionExpression: aws.String(conditional),
		}
		_, err = d.svc.PutItem(input)
	}
	if err != nil {
		d.errorLog(err)
	}
	return nil
}

func (d *DynamoRequest) UpdateTableEntry(exprAttributes, key map[string]*dynamodb.AttributeValue, updateExpression, table string) error {
	return d.UpdateItemInTable(table, key, updateExpression, exprAttributes)
}

func (d *DynamoRequest) UpdateItemInTable(tableName string, 
	key map[string]*dynamodb.AttributeValue, updateExpression string, exprAttributes map[string]*dynamodb.AttributeValue, 
) (err error) {
	input := dynamodb.UpdateItemInput{
		ExpressionAttributeValues: exprAttributes,
		TableName:                 aws.String(tableName),
		Key:                       key,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(updateExpression),
	}
	err = d.UpdateItemInternal(input)
	return
}

// UpdateItemInternal updates item using Dynamo directly
func (d *DynamoRequest) UpdateItemInternal(input dynamodb.UpdateItemInput) error {
	_, err2 := d.svc.UpdateItem(&input)
	if err2 != nil {
		d.errorLog(err2)
	}
	return err2
}

// QueryTable reads single item
// deprecated. Use GetItemFromTable or GetItemOrEmptyFromTable
func (d *DynamoRequest) QueryTable(table string, params map[string]*dynamodb.AttributeValue, out interface{}) (err error) {
	return d.GetItemFromTable(table, params, out)
}

// GetItemFromTable reads single item identified by id
// It makes sure that the item is not empty.
// If the item was not found, it returns an error - not found.
// @param out - should be non-nil pointer (https://github.com/aws/aws-sdk-go/blob/5f3810f647bffb7ed2654bf5ff0fe7b3a5ad530d/service/dynamodb/dynamodbattribute/decode.go#L86)
func (d *DynamoRequest) GetItemFromTable(table string, idParams map[string]*dynamodb.AttributeValue, out interface{}) (err error) {
	var found bool
	found, err = d.GetItemOrEmptyFromTable(table, idParams, out)
	if err == nil {
		if !found {
			err = fmt.Errorf("[NOT FOUND] in table %s value not found, idParams=%s", table, showIDParams(idParams))
		}
	}
	return
}

func showIDParams(idParams map[string]*dynamodb.AttributeValue) string {
	return strings.ReplaceAll(fmt.Sprintf("%v", idParams), "\n", " ")
}
// GetItemOrEmptyFromTable reads single item identified by id
// @param out - should be non-nil pointer (https://github.com/aws/aws-sdk-go/blob/5f3810f647bffb7ed2654bf5ff0fe7b3a5ad530d/service/dynamodb/dynamodbattribute/decode.go#L86)
func (d *DynamoRequest) GetItemOrEmptyFromTable(table string, 
	idParams map[string]*dynamodb.AttributeValue,
	out interface{}) (found bool, err error) {
	var result *dynamodb.GetItemOutput
	result, err = d.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       idParams,
	})
	// result.ConsumedCapacity.GoString()
	if err == nil {
		found = len(result.Item) > 0
		if found {
			err = dynamodbattribute.UnmarshalMap(result.Item, out)
		}
	}
	if err != nil {
		d.errorLog(err)
	}
	return
}

func attrsMapped(attrs map[string]string) map[string]*string {
	var mapped = map[string]*string{}
	for k, v := range attrs {
		mapped[k] = aws.String(v)
	}
	return mapped
}

func (d *DynamoRequest) QueryTableWithExpr(table string, condExpr string, attrNames map[string]string,
	params map[string]*dynamodb.AttributeValue, scanForward bool, limit int, out interface{}) (err error) {
	qi := &dynamodb.QueryInput{
		TableName:                 aws.String(table),
		KeyConditionExpression:    aws.String(condExpr),
		ExpressionAttributeValues: params,
		ScanIndexForward:          aws.Bool(scanForward),
	}
	// This is to avoid the following error: 'ValidationException: ExpressionAttributeNames must not be empty'
	if len(attrNames) > 0 {
		qi.ExpressionAttributeNames = attrsMapped(attrNames)
	}
	if limit > 0 {
		qi.Limit = aws.Int64(int64(limit))
	}
	result, err2 := d.svc.Query(qi)
	err = err2
	if err == nil {
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, out)
	}
	if err != nil {
		d.errorLog(err)
	}
	return
}

func (d *DynamoRequest) QueryTableWithIndex(table string, indexExpr DynamoIndexExpression, attrNames map[string]string,
	scanForward bool, limit int, out interface{}) (err error) {
	var m = map[string]*dynamodb.AttributeValue{}
	var attrNamesMapped = attrsMapped(attrNames)

	for k, v := range indexExpr.Attributes {
		dynAttr, err2 := dynamodbattribute.Marshal(v)

		if err2 != nil {
			d.errorLog(err2)
			return err2
		}
		m[k] = dynAttr
	}

	ip := &dynamodb.QueryInput{
		TableName:                 aws.String(table),
		KeyConditionExpression:    aws.String(indexExpr.Condition),
		ExpressionAttributeValues: m,
		// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Query.html
		ScanIndexForward: aws.Bool(scanForward),
	}

	if indexExpr.IndexName != "" {
		ip.IndexName = aws.String(indexExpr.IndexName)
	}

	if len(attrNamesMapped) > 0 {
		ip.ExpressionAttributeNames = attrNamesMapped
	}

	if limit > 0 {
		ip.Limit = aws.Int64(int64(limit))
	}

	result, err2 := d.svc.Query(ip)
	err = err2
	if err == nil {
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, out)
	}
	if err != nil {
		d.errorLog(err)
	}
	return
}

func (d *DynamoRequest) ScanTable(table string, out interface{}) (err error) {
	res, err2 := d.svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	err = err2
	if err == nil {
		err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &out)
	}
	if err != nil {
		d.errorLog(err)
	}
	return
}

func (d *DynamoRequest) DeleteEntry(table string, params map[string]*dynamodb.AttributeValue) error {
	_, err2 := d.svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key:       params,
	})
	if err2 != nil {
		d.errorLog(err2)
	}
	return err2
}

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to a generic interface
func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue, out interface{}) error {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)
	for k, v := range attribute {
		var dbAttr dynamodb.AttributeValue
		bytes, err2 := v.MarshalJSON()
		if err2 != nil {
			return err2
		}
		err2 = json.Unmarshal(bytes, &dbAttr)
		if err2 != nil {
			return err2
		}
		dbAttrMap[k] = &dbAttr
	}
	return dynamodbattribute.UnmarshalMap(dbAttrMap, out)
}
