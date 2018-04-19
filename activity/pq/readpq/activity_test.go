package readpq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"log"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"

	"gopkg.in/mgo.v2/bson"
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

func TestEval(t *testing.T) {

	/* defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}() */

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup

	// save test object
	//setup attrs
	message := struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		LifeCycleID string `json:"lifecycleId"`
		DeviceID    string `json:"deviceId"`
		URL         string `json:"url"`
	}{
		ID:          "100",
		Type:        "INSTALL_REQUEST",
		LifeCycleID: "1",
		DeviceID:    "remote_device_1",
		URL:         "www.testurl.com",
	}

	jsonMessage, _ := json.Marshal(message)

	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("msgq.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	messageID := bson.NewObjectId().Hex()
	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(messageID, string(jsonMessage), nil)
		return err
	})

	tc.SetInput("messageType", "INSTALL_REQUEST")

	_, err = act.Eval(tc)

	// error should be nil
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	messageValue := tc.GetOutput("message").(map[string]interface{})

	// type should be INSTALL_REQUEST
	if messageValue["type"] != "INSTALL_REQUEST" {
		fmt.Println(err)
		t.Fail()
	}

	// lifecycleId should be 1
	if messageValue["lifecycleId"] != "1" {
		fmt.Println(err)
		t.Fail()
	}

	messageKey := tc.GetOutput("messageKey").(string)

	// key should be valid objectID
	if bson.IsObjectIdHex(messageKey) == false {
		fmt.Println(err)
		t.Fail()
	}

}
