package updatepq

import (
	"encoding/json"
	"fmt"
	"log"

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
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("/home/pi/gateway/msgq.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	key := context.GetInput("key").(string)
	fieldName := context.GetInput("fieldName").(string)
	fieldValue := context.GetInput("fieldValue").(string)
	fmt.Println(fieldValue)
	err = db.Update(func(tx *buntdb.Tx) error {
		// get current message
		val, err := tx.Get(key)
		if err != nil {
			return err
		}

		// update
		msg := make(map[string]interface{})
		err = json.Unmarshal([]byte(val), &msg)
		if err != nil {
			return err
		}
		msg[fieldName] = fieldValue

		// save new message
		json, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(key, string(json), nil)
		return err
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
