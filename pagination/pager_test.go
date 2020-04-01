package pagination_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/adaptiveteam/adaptive/pagination"
	
)

func TestConcat(t *testing.T){
	p123 := pagination.InterfacePagerPure(1,2,3)
	p456 := pagination.InterfacePagerPure(4,5,6)
	p123456 := pagination.InterfacePagerConcat(p123, p456)
	res, err2 := p123456.Drain()
	assert.Nil(t, err2)
	ints := res.AsIntSlice()
	assert.Equal(t, ints, []int{1,2,3,4,5,6})
}

func TestFlatMap(t *testing.T){
	p123 := pagination.InterfacePagerPure(1,2,3)
	p122436 := p123.FlatMap(func (i interface{})pagination.InterfacePager {
		ii := i.(int)
		return pagination.InterfacePagerPure(ii, 2 * ii)
	} )
	res, err2 := p122436.Drain()
	assert.Nil(t, err2)
	ints := res.AsIntSlice()
	assert.Equal(t, ints, []int{1,2,2,4,3,6})
}

