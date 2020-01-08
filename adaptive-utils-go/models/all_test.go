package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)


func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ActionPath Suite")
}
