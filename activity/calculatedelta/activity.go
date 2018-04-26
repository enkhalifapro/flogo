package calculatedelta

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"
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
	msgType := context.GetInput("messageType").(string)
	microserviceName := context.GetInput("microserviceName").(string)

	db, err := buntdb.Open("msgq.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	deltaDb, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer deltaDb.Close()

	// get previous delta
	prevValStr := ""
	err = deltaDb.View(func(tx *buntdb.Tx) error {
		prevValStr, err = tx.Get(microserviceName + "_delta")
		if err != nil && err.Error() != "not found" {
			return err
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	// get current values from msgQ
	db.CreateIndex("index", "*", buntdb.IndexJSON("createdAt"), buntdb.IndexJSON("type"), buntdb.IndexJSON("microservice"))
	currentKey, currentValStr := "", ""
	err = db.View(func(tx *buntdb.Tx) error {
		err := tx.AscendGreaterOrEqual("index", `{"type":"`+msgType+`","microservice":"`+microserviceName+`"}`, func(key, value string) bool {
			// pick the first one
			if currentKey == "" {
				currentKey = key
				currentValStr = value
			}
			return true
		})
		return err
	})
	if err != nil {
		return false, err
	}

	// if no new changes avoid compare and return
	if currentKey == "" {
		context.SetOutput("hasDelta", false)
		context.SetOutput("delta", "")
		return true, nil
	}

	// convert prevVal json string to map[string]interface
	prevVal := make(map[string]interface{})
	if prevValStr != "" {
		err := json.Unmarshal([]byte(prevValStr), &prevVal)
		if err != nil {
			return false, err
		}
	}

	// convert currentVal json string to map[string]interface
	currentVal := make(map[string]interface{})
	if currentValStr != "" {
		err := json.Unmarshal([]byte(currentValStr), &currentVal)
		if err != nil {
			return false, err
		}
	}

	// if no prevVal send the entire result
	if prevValStr == "" {
		context.SetOutput("hasDelta", true)
		context.SetOutput("delta", currentVal)
		return true, nil
	}

	delta := CalculateDeltas(prevVal, currentVal)

	//delta is empty return empty
	if len(delta) < 1 {
		context.SetOutput("hasDelta", false)
		context.SetOutput("delta", "currentVal")
		return true, nil
	}

	// convert delta to json str
	deltaStr, err := json.Marshal(delta)
	context.SetOutput("hasDelta", true)
	if err != nil {
		return false, err
	}
	context.SetOutput("delta", deltaStr)
	return true, nil
}

// CalculateDeltas get changes of 2 dictionaries
func CalculateDeltas(prev map[string]interface{}, current map[string]interface{}) map[string]interface{} {
	diffs := make(map[string]interface{})

	// updated & deleted fields
	for pKey, pVal := range prev {

		if _, ok := pVal.(string); ok { // string val
			cVal, ok := current[pKey]
			if ok && pVal != cVal {
				// field updated
				diffs[pKey] = cVal.(string)
			}
			if !ok {
				// field deleted
				diffs[pKey] = "null"
			}
		} else if val, ok := pVal.(map[string]interface{}); ok { // object val
			cVal, isExist := current[pKey]
			if isExist && !reflect.DeepEqual(pVal, cVal) {
				// field updated
				itemDiff := CalculateDeltas(val, current[pKey].(map[string]interface{}))
				diffs[pKey] = itemDiff
			}
			if !isExist {
				// field deleted
				diffs[pKey] = "null"
			}

		} else { // if array val
			cVal, ok := current[pKey]
			if ok && !reflect.DeepEqual(pVal, cVal) {
				itemDiffs := make([]interface{}, 0)
				for i := 0; i < len(pVal.([]interface{})); i++ {
					// field updated
					itemDiff := CalculateDeltas(pVal.([]interface{})[i].(map[string]interface{}), cVal.([]interface{})[i].(map[string]interface{}))
					itemDiffs = append(itemDiffs, itemDiff)
				}
				diffs[pKey] = itemDiffs
			}
			if !ok {
				// field deleted
				diffs[pKey] = "null"
			}
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
