package models

import "fmt"

type KvPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// DistinctKvPairs returns distinct KvPairs based on the combination of key and value present them
func DistinctKvPairs(kvSlice []KvPair) (list []KvPair) {
	keys := make(map[string]bool)
	for _, entry := range kvSlice {
		keyValueAppended := fmt.Sprintf("%s:%s", entry.Key, entry.Value)
		if _, ok := keys[keyValueAppended]; !ok {
			keys[keyValueAppended] = true
			list = append(list, entry)
		}
	}
	return
}
