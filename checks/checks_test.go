package checks

import (
	"reflect"
	"testing"

	bt "github.com/adaptiveteam/business-time"
)

func TestCheckMap_Evaluate(t *testing.T) {

	expectedResults := CheckResultMap{
		"testOne":   true,
		"testTwo":   true,
		"testThree": false,
		"testFour":  false,
		"testFive":  true,
		"testSix":   false,
	}

	functionMap := CheckFunctionMap{
		"testOne":   ReturnsTrue,
		"testTwo":   ReturnsTrue,
		"testThree": ReturnsFalse,
		"testFour":  ReturnsFalse,
		"testFive":  ReturnsTrue,
		"testSix":   ReturnsFalse,
	}

	type args struct {
		userID string
		date   bt.Date
	}
	tests := []struct {
		name      string
		functions CheckFunctionMap
		args      args
		wantRv    CheckResultMap
	}{
		{
			name:      "testOne",
			functions: functionMap,
			args: struct {
				userID string
				date   bt.Date
			}{
				userID: "",
				date:   bt.NewDate(2019, 1, 1),
			},
			wantRv: expectedResults,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRv := tt.functions.Evaluate(tt.args.userID, tt.args.date); !reflect.DeepEqual(gotRv, tt.wantRv) {
				t.Errorf("CheckFunctionMap.Evaluate() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}

func TestCheckResultMap_CheckResult(t *testing.T) {
	functionMap := CheckFunctionMap{
		"testOne":   ReturnsTrue,
		"testTwo":   ReturnsTrue,
		"testThree": ReturnsFalse,
		"testFour":  ReturnsFalse,
		"testFive":  ReturnsTrue,
		"testSix":   ReturnsFalse,
	}

	checkResultMap := functionMap.Evaluate("ctcreel",bt.NewDate(2019, 1, 1),)

	type args struct {
		functionName string
		want         bool
	}
	tests := []struct {
		name    string
		crm     CheckResultMap
		args    args
		wantRv  bool
		wantErr bool
	}{
		{
			name:"basic",
			crm:checkResultMap,
			args:args{
				functionName:"testOne",
				want:true,
			},
			wantRv:true,
			wantErr:false,
		},
		{
			name:"error",
			crm:checkResultMap,
			args:args{
				functionName:"testNone",
				want:true,
			},
			wantRv:false,
			wantErr:true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRv, err := tt.crm.CheckResult(tt.args.functionName, tt.args.want)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckResultMap.CheckResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRv != tt.wantRv {
				t.Errorf("CheckResultMap.CheckResult() = %v, want %v", gotRv, tt.wantRv)
			}
		})
	}
}
