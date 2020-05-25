package userEngagementScheduling

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Global test variables
var testingT *testing.T

func TestUserEngagementScheduling(t *testing.T) {
	RegisterFailHandler(Fail)
	testingT = t
	RunSpecs(t, "UserEngagementScheduling Suite")
}
