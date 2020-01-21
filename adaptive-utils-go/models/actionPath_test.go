package models_test

import (
	"github.com/adaptiveteam/adaptive/adaptive-utils-go/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActionPath", func() {
	Context("Parsing", func(){
		It("should parse query", func(){
			urlbRel := "b?key=value"
			urlb := "/" + urlbRel
			url1 := "/a" + urlb
			a := models.ParseActionPath(url1)
			Ω(a.Path.Encode()).Should(Equal("/a/b"))
			h, t := a.HeadTail()
			Ω(h).Should(Equal("a"))
			p := models.ParseActionPath(urlb)
			Ω(t).Should(Equal(p.ToRelActionPath()))
			Ω(t.Encode()).Should(Equal(urlbRel))
		})
		It("should correctly parse path without values", func(){
			url := "/a/b"
			p := models.ParseActionPath(url)
			Ω(p.Path.Encode()).Should(Equal(url))
			Ω(len(p.Values)).Should(Equal(0))
			_, ok := p.Values["unknown key"]
			Ω(ok).Should(Equal(false))
		})
	})
})
