package aws_utils_go

import (
	// "time"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	// dc "github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	// "os"
	"testing"
)

// Global test variables
var testingT *testing.T
var resource *dockertest.Resource
var resourcePool *dockertest.Pool
var err error

const (
	hostname = "localhost"
)

func testErrorHandler(err error, msg string) {
	if err != nil {
		assert.Fail(testingT, msg+" : "+fmt.Sprint(err))
	}
}

// Test Suite
func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	testingT = t
	RunSpecs(t, "Integration Test Suite")
}

// // Running this once before the test suite
// var _ = BeforeSuite(func() {
// 	fmt.Println("Starting localstack container ... ")
// 	resourcePool, err = dockertest.NewPool("")
// 	testErrorHandler(err, "Could not connect to docker")
// 	// Starting localstack docker container with port mappings
// 	// Lambdas in golang require 'LAMBDA_EXECUTOR=docker'
// 	// Privileged access is required to start docker inside the container
// 	resource, err = resourcePool.RunWithOptions(
// 		&dockertest.RunOptions{
// 			Repository: "localstack/localstack",
// 			Tag:        "0.10.5",
// 			PortBindings: map[dc.Port][]dc.PortBinding{
// 				// S3
// 				"4572/tcp": {{HostIP: "localhost", HostPort: "4572"}},
// 				// Dynamodb
// 				"4569/tcp": {{HostIP: "localhost", HostPort: "4569"}},
// 				// Lambda
// 				"4574/tcp": {{HostIP: "localhost", HostPort: "4574"}},
// 				// SNS
// 				"4575/tcp": {{HostIP: "localhost", HostPort: "4575"}},
// 			},
// 			// Env should be []string{} for python lambdas
// 			// should be []string{"LAMBDA_EXECUTOR=docker"}, for non-python lambdas
// 			Env:        []string{"LAMBDA_EXECUTOR=docker", "DEBUG=1"},
// 			Privileged: true,
// 		},
// 	)
// 	testErrorHandler(err, "Could not start docker container")

// 	s := NewS3(os.Getenv("AWS_REGION"), fmt.Sprintf("http://%s:4572", hostname), "test")

// 	// Ensuring container is ready to accept requests
// 	if err = resourcePool.Retry(func() error {
// 		_, err = s.ListBuckets()
// 		return err
// 	}); err != nil {
// 		testErrorHandler(err, "Could not connect to docker")
// 	}
// 	time.Sleep(10*time.Second)
// 	fmt.Println("Started localstack container ... ")
// })

// // Running this once after the test suite
// var _ = AfterSuite(func() {
// 	fmt.Println("Stopping localstack container ... ")
// 	// Once tests are done, kill and remove the container
// 	if resource != nil {
// 		err1 := resource.Expire(5)
// 		testErrorHandler(err1, "There was an error with expiring localstack resource")
// 		err2 := resource.Close()
// 		testErrorHandler(err2, "There was in error with closing localstack resource")
// 		// err3 := resourcePool.Purge(resource)
// 		// testErrorHandler(err3, "There was in error with stopping the container")
// 	}
// 	time.Sleep(30*time.Second)// to free up ports
// 	fmt.Println("Stopped localstack container")
// })
