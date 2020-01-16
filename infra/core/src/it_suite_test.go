package src_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

// Test Suite
func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	testingT = t
	RunSpecs(t, "Adaptive Integration Test Suite")
}

const (
	adaptiveSlackToken = "xoxb-436528929141-492802537186-3tbN7QlbieTa27P6ROdsOoTj"
)

// Global test variables
var testingT *testing.T
var terraformOptions *terraform.Options

var clientConfigTable string
var userQueryLambda string
var usersTable string
var engagementsTable string

var awsRegion = "us-west-2"
var adaptiveBotRealName = "adaptive"

func tearDown() {
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	_, err := terraform.DestroyE(testingT, terraformOptions)
	if err != nil {
		log.Printf("Error tearing down terraform: %+v\n", err)
		terraform.Destroy(testingT, terraformOptions)
	}
}

// Running this once before the test suite
var _ = BeforeSuite(func() {
	fmt.Println("Initializing terraform ... ")
	// initializing terraform
	testVars := map[string]interface{}{
		// Hard-coding client id because this ensures no residual resources are left behind
		// If they are, test fails on the next run
		"client_id":   "test",
		"aws_region":  awsRegion,
		"environment": "test",
	}
	terraformOptions = &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../terraform/",
		// Variables to pass to our Terraform code using -var options
		Vars: testVars,
		// Remote state store configuration
		BackendConfig: map[string]interface{}{
			"bucket":         "test-adaptive-core-infra-remote-state",
			"key":            "test-terraform.tfstate",
			"region":         awsRegion,
			"dynamodb_table": "test-adaptive-core-infra-remote-state",
		},
		Lock: true,
		LockTimeout: "10s",
		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION":  awsRegion,
			"TF_VAR_RDS_HOST":     "localhost",
			"TF_VAR_RDS_USER":     "no",
			"TF_VAR_RDS_PASSWORD": "no",
			"TF_VAR_RDS_PORT":     "no",
			"TF_VAR_RDS_DB_NAME":  "no",
		},
	}

	//tearDown()
	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	msg, err := terraform.InitAndApplyE(testingT, terraformOptions)
	if err != nil {
		tearDown()
		assert.Fail(testingT, "There was an error with getting up infrastructure: "+msg)
	}

	// Run `terraform output` to get the value of an output variable
	clientConfigTable = terraform.Output(testingT, terraformOptions, "client_config_table_name")
	userQueryLambda = terraform.Output(testingT, terraformOptions, "user_query_lambda_name")
	usersTable = terraform.Output(testingT, terraformOptions, "users_table_name")
	engagementsTable = terraform.Output(testingT, terraformOptions, "user_engagements_table_name")
})

// Running this once after the test suite
var _ = AfterSuite(func() {
	fmt.Println("Destroying terraform ... ")
	tearDown()
})
