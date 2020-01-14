package src_test

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	awsutils "github.com/adaptiveteam/adaptive/aws-utils-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

// Parameters for eventual checks in tests
var defaultEventuallyTimeout = time.Duration(120 * time.Second)
var defaultEventuallyPollingInterval = 30 * time.Second

var _ = Describe("IT tests", func() {
	SetDefaultEventuallyTimeout(defaultEventuallyTimeout)
	SetDefaultEventuallyPollingInterval(defaultEventuallyPollingInterval)
	d := awsutils.NewDynamo(awsRegion, "", "core-infra-tests")
	l := awsutils.NewLambda(awsRegion, "", "core-infra-tests")

	Context("Client Config", func() {
		It("should write the client config to table in", func() {
			// Config for a client
			config := models.ClientPlatformToken{
				ClientPlatformRequest: models.ClientPlatformRequest{
					Id:  "adaptive-test",
					Org: "adaptive",
				},
				ClientPlatform: models.ClientPlatform{
					PlatformName:  models.SlackPlatform,
					PlatformToken: adaptiveSlackToken,
				},
				ClientContact: models.ClientContact{
					ContactFirstName: "John",
					ContactLastName:  "Doe",
					ContactMail:      "john.doe@test.com",
				},
			}
			// Writing client config to the table
			err := d.PutTableEntry(config, clientConfigTable)
			Expect(err).To(BeNil())
		})

		It("should invoke user query lambda", func() {
			// Invoking user query lambda
			op, err := l.InvokeFunction(userQueryLambda, []byte("{}"), false)
			Expect(err).To(BeNil())
			// Asserting ok response
			Eventually(func() int64 {
				return *op.StatusCode
			}).Should(Equal(int64(200)))
			// Asserting on output payload
			Eventually(func() string {
				return string(op.Payload)
			}).Should(Equal("1"))
		})

		It("should write users to table", func() {

			getUsers := func() []models.User {
				var users []models.User
				// Scan users table
				err := d.ScanTable(usersTable, &users)
				Expect(err).To(BeNil())
				return users
			}

			// There are no user communities, hence there are no Adaptive associated users
			Eventually(func() int {
				return len(getUsers())
			}).Should(Equal(0))
		})
	})
})
