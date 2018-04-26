package calculatedelta

import (
	"fmt"
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
	// read first 2 messages by type
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("msgq.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// register query indexes

	// query
	msgType := context.GetInput("messageType").(string)
	microserviceName := context.GetInput("microserviceName").(string)

	msgs := make(map[string]string)
	db.CreateIndex("index", "*", buntdb.IndexJSON("createdAt"), buntdb.IndexJSON("type"), buntdb.IndexJSON("microservice"))
	//db.CreateIndex("type", "*", )
	err = db.View(func(tx *buntdb.Tx) error {
		fmt.Println("in 1")
		err := tx.AscendGreaterOrEqual("index", `{"type":"`+msgType+`","microservice":"`+microserviceName+`"}`, func(key, value string) bool {
			if len(msgs) < 2 {
				msgs[key] = value
			}
			return true
		})
		return err
	})
	if err != nil {
		return false, err
	}

	// if no new messages
	if len(msgs) < 1 {
		context.SetOutput("hasDelta", false)
		context.SetOutput("delta", "")
		return true, nil
	}

	// if it first message send all data
	if len(msgs) == 1 {
		context.SetOutput("hasDelta", true)
		for _, msg := range msgs {
			context.SetOutput("delta", msg)
		}
		return true, nil
	}

	// if we have 2 messages then calculate deltas

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
