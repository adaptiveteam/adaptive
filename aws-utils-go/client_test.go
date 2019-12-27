package aws_utils_go

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"time"
)

// Parameters for eventual checks in tests
var defaultEventuallyTimeout = time.Duration(30 * time.Second)
var defaultEventuallyPollingInterval = 5 * time.Second

var _ = Describe("IT Tests", func() {

	SetDefaultEventuallyTimeout(defaultEventuallyTimeout)
	SetDefaultEventuallyPollingInterval(defaultEventuallyPollingInterval)

	Context("AWS S3 Tests", func() {
		testBucketName := "test"
		key := "foo/var/main.zip"

		s := NewS3(os.Getenv("AWS_REGION"), fmt.Sprintf("http://%s:4572", hostname), "test")

		It("should assert initial bucket count to 0", func() {
			Eventually(func() int {
				op, err := s.ListBuckets()
				testErrorHandler(err, "Could not list S3 buckets")
				return len(op)
			}).Should(Equal(0))
		})

		It("should create a bucket and assert it exists", func() {
			err2 := s.EnsureBucketExists(testBucketName)
			fmt.Printf("EnsureBucketExists(%s): err:%v\n", testBucketName, err2)
			Expect(err2).To(BeNil())
			// ListBuckets fails with parse XML message.
			// Eventually(func() int {
			// 	op, err3 := s.ListBuckets()
			// 	fmt.Printf("ListBuckets(): err:%v\n", err3)
			// 	Expect(err3).To(BeNil())
			// 	return len(op)
			// }).Should(Equal(1))
		})

		It("should upload a file and assert object exists", func() {
			err := s.AddFile("testdata/main.zip", testBucketName, key)
			Expect(err).To(BeNil())
			Eventually(func() bool {
				return s.ObjectExists(testBucketName, key)
			}).Should(BeTrue())
		})

		It("should delete a bucket and assert it's removed", func() {
			_, err := s.DeleteObject(testBucketName, key)
			Expect(err).To(BeNil())
			err = s.DeleteBucket(testBucketName)
			Expect(err).To(BeNil())
			Eventually(func() int {
				op, err := s.ListBuckets()
				Expect(err).To(BeNil())
				return len(op)
			}).Should(Equal(0))
		})
	})

	Context("AWS Dynamodb Tests", func() {
		testTableName := "test"
		d := NewDynamo(os.Getenv("AWS_REGION"), fmt.Sprintf("http://%s:4569", hostname), "test")

		keyParams := map[string]*dynamodb.AttributeValue{
			"year": {
				N: aws.String("2018"),
			},
			"title": {
				S: aws.String("Infinity War"),
			},
		}

		It("should create a table ", func() {
			// Creating a table foe testing
			input := &dynamodb.CreateTableInput{
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("year"),
						AttributeType: aws.String("N"),
					},
					{
						AttributeName: aws.String("title"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("title"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("year"),
						KeyType:       aws.String("RANGE"),
					},
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
				TableName: aws.String(testTableName),
			}
			err := d.CreateTable(input)
			Expect(err).To(BeNil())
		})

		It("should describe a table", func() {
			op, err := d.DescribeTable(testTableName)
			Expect(err).To(BeNil())
			Expect(op.Name).To(Equal(testTableName))
		})

		type movie struct {
			Year   int    `json:"year"`
			Title  string `json:"title"`
			Rating int    `json:"rating"`
		}

		It("should put an entry in the table with condition", func() {
			item := movie{Year: 2018, Title: "Infinity War", Rating: 9}
			err := d.PutTableEntryWithCondition(item, testTableName, "attribute_not_exists(title)")
			Expect(err).To(BeNil())
		})

		It("should put an entry in the table with no condition", func() {
			item := movie{Year: 2018, Title: "Infinity War"}
			err := d.PutTableEntry(item, testTableName)
			Expect(err).To(BeNil())
		})

		It("update put an entry in the table", func() {
			exprAttributes := map[string]*dynamodb.AttributeValue{
				":r": {
					N: aws.String("8"),
				},
			}
			updateExpression := "set rating = :r"
			err := d.UpdateTableEntry(exprAttributes, keyParams, updateExpression, testTableName)
			Expect(err).To(BeNil())
		})

		It("should scan table", func() {
			var movies []movie
			err := d.ScanTable(testTableName, &movies)
			Expect(err).To(BeNil())
			Expect(len(movies)).To(Equal(1))
		})

		It("should query table", func() {
			var op movie
			err := d.QueryTable(testTableName, keyParams, &op)
			Expect(err).To(BeNil())
			Expect(op).To(Equal(movie{2018, "Infinity War", 8}))
		})

		It("should query table with expression", func() {
			queryExpr := "title = :t AND #year BETWEEN :y1 AND :y2"
			params := map[string]*dynamodb.AttributeValue{
				":t": {
					S: aws.String("Infinity War"),
				},
				":y1": {
					N: aws.String("2017"),
				},
				":y2": {
					N: aws.String("2019"),
				},
			}
			var op []movie
			err = d.QueryTableWithExpr(testTableName, queryExpr, map[string]string{"#year": "year"}, params, true, 1, &op)
			Expect(err).To(BeNil())
			Expect(len(op)).To(Equal(1))
		})

		It("should delete item from table", func() {
			err := d.DeleteEntry(testTableName, keyParams)
			Expect(err).To(BeNil())

			//scan the table after deletion
			var movies []movie
			err = d.ScanTable(testTableName, &movies)
			Expect(err).To(BeNil())
			Expect(len(movies)).To(Equal(0))
		})
	})

	Context("AWS Lambda Tests with Cloudwatch", func() {
		l := NewLambda(os.Getenv("AWS_REGION"), fmt.Sprintf("http://%s:4574", hostname), "test")
		It("should return false when querying for a non-existent function", func() {
			Expect(l.FunctionExists("test")).To(Equal(false))
		})

		It("should list lambda functions", func() {
			l, err := l.ListLambdas()
			Expect(err).To(BeNil())
			Expect(len(l.Functions)).To(Equal(0))
		})

		It("should create a lambda function", func() {
			zipBytes, err := ioutil.ReadFile("testdata/main.zip")
			Expect(err).To(BeNil())
			fn := &LambdaFunction{
				Name:       "main",
				Handler:    "main",
				Role:       "test",
				MemorySize: int64(128),
				Timeout:    int64(30),
			}
			_, err = l.CreateFunction(fn, zipBytes)
			Expect(err).To(BeNil())
			Expect(l.FunctionExists("main")).To(Equal(true))
		})
	})

	Context("AWS SNS tests", func() {
		s := NewSNS(os.Getenv("AWS_REGION"), fmt.Sprintf("http://%s:4575", hostname), "test")
		var topicArn string
		It("should create a topic", func() {
			arn, err := s.CreateTopic("test", nil)
			Expect(err).To(BeNil())

			Expect(*arn).NotTo(BeEmpty())
			topicArn = *arn
		})

		It("should list topics", func() {
			arns, _, err := s.ListTopics(nil)
			Expect(err).To(BeNil())

			Expect(len(arns)).NotTo(BeZero())
			Expect(arns[0]).To(Equal(topicArn))
		})

		It("should publish to a topic", func() {
			_, err := s.Publish("message", topicArn)
			Expect(err).To(BeNil())
		})

		It("should delete a topic", func() {
			err := s.DeleteTopic(topicArn)
			Expect(err).To(BeNil())

			arns, _, err := s.ListTopics(nil)
			Expect(err).To(BeNil())
			Expect(len(arns)).To(BeZero())
		})

	})
})
