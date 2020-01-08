package models_test

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optional operations", func() {
	Context("Optional operations", func(){
		It("concat optional value", func(){
			prefix := "a"
			text := "b"
			Ω(models.ConcatPrefixOpt(prefix, "")).Should(Equal(""))
			Ω(models.ConcatPrefixOpt(prefix, text)).Should(Equal("ab"))
		})
	})
})
