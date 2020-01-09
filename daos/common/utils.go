package common
// This file is handcrafted!
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strconv"
)

func DynString(str string) (attr *dynamodb.AttributeValue) {
	return DynS(str)
}
// DynS wraps a string into dynamo db attribute value
// if the string is empty, uses `NULL: true`
func DynS(str string) (attr *dynamodb.AttributeValue) {
	if str == "" {
		attr = &dynamodb.AttributeValue{NULL: aws.Bool(true)} 
	} else {
		attr = &dynamodb.AttributeValue{S: aws.String(str)}
	}
	return 
}
func DynN(i int) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(i))}
}
func DynBOOL(b bool) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{BOOL: &b}
}
func DynSS(list []string) (attr *dynamodb.AttributeValue) {
	return &dynamodb.AttributeValue{SS: aws.StringSlice(list)}
}
func StringArraysEqual(arr1 []string, arr2 []string) (res bool) {
	res = len(arr1) == len(arr2)
	if res {
		for i := range arr1 {
			res = res && arr1[i] == arr2[i]
		}
	}
	return
}
