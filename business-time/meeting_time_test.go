package business_time

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)


func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LocalTime Suite")
}

var _ = Describe("LocalTime", func() {
	Context("Parsing", func(){
		It("should parse identifier", func(){
			Ω(ParseLocalTimeID("1845")).Should(Equal( MeetingTime(18, 45)))
		})
		It("should render AM", func(){
			Ω(MeetingTime(18, 45).ToUserFriendly()).Should(Equal( "06:45 PM"))
		})
	})
})
