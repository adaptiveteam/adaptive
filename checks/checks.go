package checks

import (
	"time"
	// "strings"
	"math/rand"
	"log"
	"fmt"
	bt "github.com/adaptiveteam/adaptive/business-time"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
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
	deltaT time.Duration
	err    error
}


func runningtime(output *checkResultOut) (*checkResultOut, time.Time) {
	output.deltaT = -1 // in case of an error
    return output, time.Now()
}

func track(output *checkResultOut, startTime time.Time) {
	output.deltaT = time.Since(startTime)
}

func SafeWrapCheckFunction(input checkResultIn) (output checkResultOut) {
	defer core.RecoverToErrorVar("go " + input.name, &output.err)
	output.name = input.name
	defer track(runningtime(&output))
	output.result = input.function(input.userID, input.date)
	return
}

func (functions CheckFunctionMap) Evaluate(userID string, date bt.Date) (rv CheckResultMap) {
	rv = make(CheckResultMap, 0)
	id := rand.Int()
	start := time.Now()
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
	outputs := []checkResultOut{}
	// log.Printf("Evaluate(%d): # functions to wait for: %d. Functions are: %v\n", id, len(names), names)
	for range filteredFunctionNames {
		out := <-outChannel
		outputs = append(outputs, out)
		if out.err != nil {
			log.Printf("Evaluate(%d): Error in CheckFunctionMap) Evaluate of function %s: %+v\n", id, out.name, out.err)
		}
		if out.deltaT > 50 * time.Millisecond {
			log.Printf("CheckFunctionMap.Evaluate(%d): Evaluate of function %s took %v ms\n", id, out.name, out.deltaT / time.Millisecond)
		}
		rv[out.name] = out.result
		delete(names, out.name)
		// log.Printf("Evaluate(%d): # functions left to wait for: %d. Functions are: %v\n", id, len(names), names)
	}
	deltaTime := time.Since(start)
	if deltaTime > time.Second {
		log.Printf("CheckFunctionMap.Evaluate took %v ms. Here is the list of individual checks:\n", deltaTime / time.Millisecond)
		for _, out := range outputs {
			log.Printf("%s: %v ms\n", out.name, out.deltaT / time.Millisecond)
		}
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