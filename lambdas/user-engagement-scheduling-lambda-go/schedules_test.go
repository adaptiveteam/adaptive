package userEngagementScheduling

import (
	"time"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

)

var _ = Describe("Schedules", func() {
	Context("everyQuarterLastMonthEveryFridayAt12", func(){
		It("should trigger on first Friday", func(){
			res := everyQuarterLastMonthEveryFridayAt12.IsOnSchedule(
				time.Date(2020, 06, 05, 12, 05, 0, 0, time.UTC))
			Ω(res).Should(BeTrue())
		})
		It("should trigger on last Friday", func(){
			res := everyQuarterLastMonthEveryFridayAt12.IsOnSchedule(
				time.Date(2020, 06, 26, 12, 05, 0, 0, time.UTC))
			Ω(res).Should(BeTrue())
		})
		It("should not trigger on first Saturday", func(){
			res := everyQuarterLastMonthEveryFridayAt12.IsOnSchedule(
				time.Date(2020, 06, 06, 12, 05, 0, 0, time.UTC))
			Ω(res).Should(BeFalse())
		})
	})
})
