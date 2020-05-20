package fetch_dialog

import (
	"testing"
	"reflect"
	. "github.com/onsi/ginkgo"
)

func Init() {
	testing.Init()
	if !testing.Short() {
		describeAllTests() 
	}
}

func describeAllTests(){ Describe("Fetching dialog using DAO", func() {
	Context("FetchDialogBySubject", func(){
		It("should fetch a predefined dialog", func(){
			dao := localStackDao()
			result1, err := dao.FetchByContextSubject(knownDialogContext, mockSubject)
			if err != nil { testingT.Errorf("Error in FetchByContextSubject(%v)", err)}
			result2, err := dao.FetchByAlias(mockPackageName, mockAlias, mockSubject)
			if err != nil { testingT.Errorf("Error in FetchByAlias(%v)", err)}
			if !reflect.DeepEqual(result1, result2) {
				testingT.Errorf("Expected result1=%v\n result2=%v, to be the same.", result1, result2)
			}
		})
	})
	Context("InMemoryDao", func(){
		It("should fetch a dialog from mock", func(){
			CheckDialogMocking(testingT)
		})
	})
	Context("FetchDialogByMock", func(){
		It("should fetch a dialog from mock", func(){
			CheckDialogEntryEquality(testingT)
		})
	})
})
}
