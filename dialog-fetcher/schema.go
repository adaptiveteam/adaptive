package fetch_dialog


import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	SchemaDialogTable string = "dialogs"
	SchemaDialogAliasesTable string = "dialogs_alias"
)

var (
	provisionedThroughput = dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}
)

func createDialogTable(db *dynamodb.DynamoDB) error {
	schemaDialogTable := SchemaDialogTable
	dialog_id_field_name := "dialog_id"
	dialog_id_field_type := dynamodb.ScalarAttributeTypeS
	dialog_id := dynamodb.AttributeDefinition{
		AttributeName: &dialog_id_field_name,
		AttributeType: &dialog_id_field_type,
	}
	context_field_name := "context"
	context_field_type := dynamodb.ScalarAttributeTypeS
	context := dynamodb.AttributeDefinition{
		AttributeName: &context_field_name,
		AttributeType: &context_field_type,
	}
	subject_field_name := "subject"
	subject_field_type := dynamodb.ScalarAttributeTypeS
	subject := dynamodb.AttributeDefinition{
		AttributeName: &subject_field_name,
		AttributeType: &subject_field_type,
	}
	context_subject_index_name := "context-subject-index"
	projectionType := "ALL"
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		TableName: &schemaDialogTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			&dialog_id,
			&context,
			&subject,
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(dialog_id_field_name),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &provisionedThroughput,
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			&dynamodb.GlobalSecondaryIndex{
				IndexName: &context_subject_index_name,
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("context"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("subject"),
						KeyType:       aws.String("RANGE"),
					},
				},
				ProvisionedThroughput: &provisionedThroughput,
				Projection: &dynamodb.Projection{ProjectionType: &projectionType},
			},
		},
	})
	if err != nil {
		fmt.Printf("Couldn't create table %s: %v\n", schemaDialogTable, err)
	} else {
		fmt.Printf("Created table %s\n", schemaDialogTable)
	}
	return err
}
func createDialogAliasTable(db *dynamodb.DynamoDB) error {
	schemaDialogAliasesTable := SchemaDialogAliasesTable

	application_alias_field_name := "application_alias"
	application_alias_field_type := dynamodb.ScalarAttributeTypeS
	application_alias := dynamodb.AttributeDefinition{
		AttributeName: &application_alias_field_name,
		AttributeType: &application_alias_field_type,
	}
	
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		TableName: &schemaDialogAliasesTable,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{&application_alias},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(application_alias_field_name),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &provisionedThroughput,
	})
	if err != nil {
		fmt.Printf("Couldn't create table %s: %v\n", schemaDialogAliasesTable, err)
	} else {
		fmt.Printf("Created table %s\n", schemaDialogAliasesTable)
	}
	return err
}

func localStackInitializeSchema(db *dynamodb.DynamoDB) error {
	err := createDialogTable(db)
	if err == nil {
		err = createDialogAliasTable(db)
	}
	return err
}
