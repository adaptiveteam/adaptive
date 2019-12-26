package core_utils_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListInAAndB(t *testing.T) {
	listA := []string{"a", "b", "c", "d"}
	listB := []string{"c", "d", "e", "f"}

	assert.True(t, UnorderedEqual(InAAndB(listA, listB), []string{"c", "d"}))
}

func TestListInAButNotB(t *testing.T) {
	listA := []string{"a", "b", "c", "d"}
	listB := []string{"c", "d", "e", "f"}

	assert.True(t, UnorderedEqual(InAButNotB(listA, listB), []string{"a", "b"}))
}

func TestListInBButNotA(t *testing.T) {
	listA := []string{"a", "b", "c", "d"}
	listB := []string{"c", "d", "e", "f"}

	assert.True(t, UnorderedEqual(InBButNotA(listA, listB), []string{"e", "f"}))
}

func TestPickDialogOption(t *testing.T) {
	listA := []string{"a", "b", "c", "d"}
	for i := 0; i < 10; i++ {
		value := []string{RandomString(listA)}
		assert.True(t, UnorderedEqual(InAAndB(value, listA), value))
	}
}

func TestUnique(t *testing.T) {
	listA := []string{"a", "b", "c", "d", "a", "c"}
	assert.True(t, UnorderedEqual(Distinct(listA), []string{"a", "b", "c", "d"}))
}
