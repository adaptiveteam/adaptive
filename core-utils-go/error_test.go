package core_utils_go

import (
	"github.com/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicRecovery(t *testing.T) {
	a, b, err2 := myFunc()
	assert.Equal(t, 1, a)
	assert.Equal(t, 0, b)
	assert.NotEmpty(t, err2)
	assert.Contains(t, err2.Error(), "CHECK")
}

func myFunc() (a, b int, err error) {
	defer RecoverToErrorVar("myFunc", &err)
	a = 1
	if a == 1 {
		panic(errors.New("Error in myFunc. Code = CHECK"))
	}
	b = 1
	return
}

func myFunc2() (a, b int, err error) {
	defer func (name string) {
		err2 := recover()
		if err2 != nil {
			// log.Printf("RecoverToErrorVar2 (%s) (err=%+v), (err2: %+v\n", name, *err, err2)
			switch err2.(type) {
			case error:
				err3 := err2.(error)
				err4 := errors.Wrapf(err3, "%s: Recover from panic", name)
				err = err4
			case string:
				err3 := err2.(string)
				err4 := errors.New(name + ": Recover from string-panic: " + err3)
				err = err4
			// default:
			// 	err4 := errors.New(fmt.Sprintf("%s: Recover from unknown-panic: %+v", name, err2))
			// 	err = &err4
			}
		}
		// if *err == nil && err2 != nil {
		// 	// log.Printf("RecoverToErrorVar3: could not assign to err var\n")
		// }
	}("myFunc")
	a = 1
	if a == 1 {
		panic(errors.New("Error in myFunc. Code = CHECK"))
	}
	b = 1
	return
}
