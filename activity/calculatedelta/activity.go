package calculatedelta

import (
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	isSame := reflect.DeepEqual(nil, nil)
	if !isSame {
		context.SetOutput("isSame", false)
		context.SetOutput("changes", nil)
		return true, nil
	}
	return true, nil
}

// CalculateDeltas get changes of 2 dictionaries
func CalculateDeltas(prev map[string]string, current map[string]string) map[string]string {
	diffs := make(map[string]string)

	// updated & deleted fields
	for pKey, pVal := range prev {
		cVal, ok := current[pKey]
		if ok && pVal != cVal {
			// field updated
			diffs[pKey] = cVal
		} else {
			// field deleted
			diffs[pKey] = "null"
		}
	}

	// created fields
	for cKey, cVal := range current {
		_, ok := prev[cKey]
		if !ok {
			// field created
			diffs[cKey] = cVal
		}
	}

	return diffs
}
