package calculatedelta

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

// add 3 new messages
func prepareQTestDB() error {
	db, err := buntdb.Open("./MICROSERVICE_READ_VALUES.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *buntdb.Tx) error {
		// set first row
		_, _, err = tx.Set("5ae6216b415dc4b0816502bf", `{"type":"MICROSERVICE_READ_VALUES","topic":"event/OEM/MODEL/NAME/VERSION/00000000-aaaaaaaa-11111111","id":"00000000-aaaaaaaa-11111111","attribute":"Type","value":"N/A","timeStamp":1525031275}`, nil)
		if err != nil {
			return err
		}

		// set second row
		_, _, err = tx.Set("5ae6216b415dc4b081650516", `{"type":"MICROSERVICE_READ_VALUES","topic":"event/OEM/MODEL/NAME/VERSION/00000000-aaaaaaaa-11111111","id":"00000000-aaaaaaaa-11111111","attribute":"Revision","value":"N/A","timeStamp":1525031275}`, nil)
		if err != nil {
			return err
		}

		// set third row
		_, _, err = tx.Set("5ae6216b415dc4b08165056b", `{"type":"MICROSERVICE_READ_VALUES","topic":"event/OEM/MODEL/NAME/VERSION/00000000-aaaaaaaa-11111111","id":"00000000-aaaaaaaa-11111111","attribute":"Serial Number","value":"#FFFFFFF","timeStamp":1525031275}`, nil)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// add 3 new messages
func prepareGatewayTestDB() error {
	db, err := buntdb.Open("./gateway.db")
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *buntdb.Tx) error {
		// set first row
		_, _, err = tx.Set("dc3070c8", "N/A", nil)
		return err
	})
	return err
}

func TestEval(t *testing.T) {

	/* defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}() */

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	err := prepareQTestDB()
	if err != nil {
		t.Fail()
	}
	err = prepareGatewayTestDB()
	if err != nil {
		t.Fail()
	}
	tc.SetInput("messageType", "MICROSERVICE_READ_VALUES")
	done, err := act.Eval(tc)

	//check result attr
	if done != true {
		t.Fail()
	}
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// buffer count should be 2
	db, err := buntdb.Open("./MICROSERVICE_READ_VALUES_buffer.db")
	if err != nil {
		t.Fail()
	}
	err = db.View(func(tx *buntdb.Tx) error {
		// set first row
		count, err := tx.Len()
		if err != nil {
			t.Fail()
		}
		if count != 2 {
			t.Fail()
		}
		return err
	})
}
