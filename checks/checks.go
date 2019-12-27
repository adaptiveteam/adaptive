package checks

import (
	// "strings"
	"math/rand"
	"log"
	"fmt"
	bt "github.com/adaptiveteam/adaptive/business-time"
)

type checkResultIn struct {
	name     string
	function CheckFunction
	userID   string
	date     bt.Date
}

type checkResultOut struct {
	name   string
	result bool
	err    error
}

func SafeWrapCheckFunction(input checkResultIn) (output checkResultOut) {
	defer recoverToErrorVar("go " + input.name, &output.err)
	output.name = input.name
	output.result = input.function(input.userID, input.date)
	return
}

func (functions CheckFunctionMap) Evaluate(userID string, date bt.Date) (rv CheckResultMap) {
	rv = make(CheckResultMap, 0)
	id := rand.Int()
	// Now run all of the checks in parallel

	var filteredFunctionNames []string
	for k := range functions {
		// if !strings.HasPrefix(k, "ObjectivesExist") {
			filteredFunctionNames = append(filteredFunctionNames, k)
		// }
	}
	outChannel := make(chan checkResultOut, len(functions))
	for _, k := range filteredFunctionNames {
		v := functions[k]
		input1 := checkResultIn{
			name:     k,
			function: v,
			userID:   userID,
			date:     date,
		}
		go func(input checkResultIn, outChannel chan checkResultOut) {
			outChannel <- SafeWrapCheckFunction(input)
		}(input1, outChannel)
	}
	names := SliceToSet(filteredFunctionNames)
	// log.Printf("Evaluate(%d): # functions to wait for: %d. Functions are: %v\n", id, len(names), names)
	for range filteredFunctionNames {
		out := <-outChannel
		if out.err != nil {
			log.Printf("Evaluate(%d): Error in CheckFunctionMap) Evaluate of function %s: %+v\n", id, out.name, out.err)
		}
		rv[out.name] = out.result
		delete(names, out.name)
		// log.Printf("Evaluate(%d): # functions left to wait for: %d. Functions are: %v\n", id, len(names), names)
	}

	return rv
}


func (crm CheckResultMap) CheckResult(functionName string, want bool) (rv bool, err error) {
	if val, ok := crm[functionName]; ok {
		rv = val==want
		err = nil
	} else {
		rv = false
		err = fmt.Errorf("check function is not in list %v",functionName)
	}
	return rv,err
}