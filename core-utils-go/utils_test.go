package core_utils_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIfThenElse(t *testing.T) {
	iteT := IfThenElse(true, 1, 2)
	assert.True(t, iteT == 1)

	iteF := IfThenElse(false, 1, 2)
	assert.True(t, iteF == 2)
}
