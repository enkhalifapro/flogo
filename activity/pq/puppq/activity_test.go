package puppq

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	message := struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		LifeCycleID string `json:"lifecycleId"`
		DeviceID    string `json:"deviceId"`
		URL         string `json:"url"`
	}{
		ID:          "100",
		LifeCycleID: "1",
		DeviceID:    "remote_device_1",
		URL:         "www.testurl55.com",
	}

	jsonMessage, _ := json.Marshal(message)
	tc.SetInput("message", string(jsonMessage))

	_, err := act.Eval(tc)

	//check result attr

	// error should be nil
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	// result should be true
	result := tc.GetOutput("result").(bool)
	if result != true {
		t.Fail()
	} */
}
