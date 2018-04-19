package time

import (
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/labstack/gommon/log"
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

	// Get input from context
	zone := context.GetInput("zone").(string)

	// calculate time
	loc, err := time.LoadLocation(zone)
	if err != nil {
		log.Errorf(err.Error())
		return false, err
	}

	time := time.Now().In(loc)

	context.SetOutput("time", time.String())
	return true, nil
}
