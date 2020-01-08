package models

import (
	"gotest.tools/assert"
	"testing"
)

func TestDistinctKvPairsWithSameKeyValues(t *testing.T) {
	kvPairSlice1 := []KvPair{
		{Key: "key1", Value: "value1",},
		{Key: "key2", Value: "value2",},
	}
	kvPairSlice2 := []KvPair{
		{Key: "key2", Value: "value2",},
		{Key: "key3", Value: "value3",},
	}

	res := DistinctKvPairs(append(kvPairSlice1, kvPairSlice2...))
	assert.Assert(t, len(res) == 3)
}

func TestDistinctKvPairsWithDifferentKeyValues(t *testing.T) {
	kvPairSlice1 := []KvPair{
		{Key: "key1", Value: "value1",},
		{Key: "key2", Value: "value2",},
	}
	kvPairSlice2 := []KvPair{
		{Key: "key2", Value: "value1",},
		{Key: "key3", Value: "value3",},
	}

	res := DistinctKvPairs(append(kvPairSlice1, kvPairSlice2...))
	assert.Assert(t, len(res) == 4)
}
