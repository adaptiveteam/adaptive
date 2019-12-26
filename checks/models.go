package checks

import (
	bt "github.com/adaptiveteam/business-time"
	"log"
)

type CheckFunction func(userID string, date bt.Date) (rv bool)

type CheckFunctionMap map[string]CheckFunction

type CheckResultMap map[string]bool

func (cf CheckFunction)And(other CheckFunction) CheckFunction {
	return func(userID string, date bt.Date) (rv bool) {
		return cf(userID, date) && other(userID, date)
	}
}

func (cf CheckFunction)Or(other CheckFunction) CheckFunction {
	return func(userID string, date bt.Date) (rv bool) {
		return cf(userID, date) || other(userID, date)
	}
}

func Not(other CheckFunction) CheckFunction {
	return func(userID string, date bt.Date) (rv bool) {
		return !other(userID, date)
	}
}

func AndMany(functions []CheckFunction) CheckFunction {
	return func(userID string, date bt.Date) (rv bool) {
		for _, f := range functions {
			if !f(userID, date){
				return false
			}
		}
		return true
	}
}

func OrMany(functions []CheckFunction) CheckFunction {
	return func(userID string, date bt.Date) (rv bool) {
		for _, f := range functions {
			if f(userID, date){
				return true
			}
		}
		return false
	}
}

func (m CheckFunctionMap)GetUnsafe(name string) (f CheckFunction) {
	f = m[name] 
	if f == nil {
		log.Panicln(name, " is not a valid check")
	}
	return
}

func (m CheckFunctionMap)Keys() (keys []string) {
    keys = make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
	}
	return
}

func SliceToSet(s []string) (set map[string]struct{}) {
	set = make(map[string]struct{})
    for _, k := range s {
        set[k] = struct{}{}
	}
	return
}

func (m CheckFunctionMap)KeySet() (map[string]struct{}) {
	return SliceToSet(m.Keys())
}

var True  CheckFunction = func(_ string, _ bt.Date) bool { return true }
var False CheckFunction = func(_ string, _ bt.Date) bool { return false }
