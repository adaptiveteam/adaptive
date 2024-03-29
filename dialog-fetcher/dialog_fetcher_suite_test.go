package fetch_dialog

import (
	"fmt"
	"testing"
	// "time"

	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	// "github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"

	// . "github.com/adaptiveteam/adaptive/aws-utils-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Global test variables
var testingT *testing.T
var resource *dockertest.Resource
var resourcePool *dockertest.Pool

const (
	hostname         = "localhost"
	dynamoDbEndpoint = "http://" + hostname + ":4569"
)

func TestDialogFetcher(t *testing.T) {
	if !testing.Short() {
		RegisterFailHandler(Fail)
		testingT = t
		RunSpecs(t, "DialogFetcher Suite")
	}
}

func testErrorHandler(err error, msg string) {
	if err != nil {
		assert.Fail(testingT, msg+" : "+fmt.Sprint(err))
	}
}

// Running this once before the test suite
var _ = BeforeSuite(func() {
	// var err error
	// fmt.Println("Starting localstack container ... ")
	// resourcePool, err = dockertest.NewPool("")
	// testErrorHandler(err, "Could not connect to docker")
	// // Starting localstack docker container with port mappings
	// // Lambdas in golang require 'LAMBDA_EXECUTOR=docker'
	// // Privileged access is required to start docker inside the container
	// resource, err = resourcePool.RunWithOptions(
	// 	&dockertest.RunOptions{
	// 		Repository: "localstack/localstack",
	// 		Tag:        "0.10.5",
	// 		PortBindings: map[docker.Port][]docker.PortBinding{
	// 			// S3
	// 			"4572/tcp": {{HostIP: hostname, HostPort: "4572"}},
	// 			// Dynamodb
	// 			"4569/tcp": {{HostIP: hostname, HostPort: "4569"}},
	// 			// Lambda
	// 			"4574/tcp": {{HostIP: hostname, HostPort: "4574"}},
	// 			// SNS
	// 			"4575/tcp": {{HostIP: hostname, HostPort: "4575"}},
	// 		},
	// 		// Env should be []string{} for python lambdas
	// 		// should be []string{"LAMBDA_EXECUTOR=docker"}, for non-python lambdas
	// 		Env:        []string{"LAMBDA_EXECUTOR=docker", "DEBUG=1"},
	// 		Privileged: true,
	// 	},
	// )
	// if err != nil {
	// 	fmt.Printf("Ignoring error when creating docker pool: %+v", err)
	// 	resource = nil
	// }
	// // testErrorHandler(err, "Could not start docker container")

	// // Ensuring container is ready to accept requests
	// if err = resourcePool.Retry(func() error {
	// 	s := NewS3(testAwsRegion(), fmt.Sprintf("http://%s:4572", hostname), "test")
	// 	_, err = s.ListBuckets()
	// 	return err
	// }); err != nil {
	// 	testErrorHandler(err, "Could not connect to docker")
	// }
	// time.Sleep(10*time.Second)

	// fmt.Println("Started localstack container")
	if !testing.Short() {
		addDialogSchema()
		err2 := addMockData(localStackDao())
		if err2 != nil {
			testingT.Errorf("Problem creating mock data: %+v", err2)
		}
	}
})

func localStackDynamoDB() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("foo", "var", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.UsEast1RegionID),
		DisableSSL:       aws.Bool(true),
	}))
	conf := &aws.Config{
		Endpoint: aws.String(dynamoDbEndpoint),
	}
	return dynamodb.New(sess, conf)
}

func localStackDynamoRequest() *awsutils.DynamoRequest {
	return awsutils.NewDynamo(testAwsRegion(), dynamoDbEndpoint, "dialogs-testing")
}

func localStackDao() DAO {
	return NewDAO(
		localStackDynamoRequest(), SchemaDialogTable,
	)
}

func addDialogSchema() error {
	db := localStackDynamoDB()
	return localStackInitializeSchema(db)
}

func testAwsRegion() string {
	return "us-east-1" // we don't need real region in tests getNonEmptyEnv("AWS_REGION")
}
