package core_utils_go


import (
	"github.com/Merovius/go-misc/lazy"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lazy(t *testing.T) {
	count := 0
	a := lazy.Int(func () int { 
		count ++
		return count
	})

	b := lazy.Int(func () int { 
		count ++
		return count
	})

	assert.Equal(t, a(), a())
	assert.Equal(t, 3, a() + b())
	assert.Equal(t, 3, a() + b())

	assert.Equal(t, 2, count)
}
