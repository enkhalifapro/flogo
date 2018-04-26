package calculatedelta

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
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
	}() */

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("messageType", "MICROSERVICE_READ_VALUES")
	tc.SetInput("microserviceName", "CMA")
	act.Eval(tc)

	//check result attr
}

func getPrevValues() map[string]interface{} {
	str := `{
	"devices": [{
		"Comment": "test5",
		"External Name": "",
		"Firmware Version": "\"7.2/1.8\"",
		"Model": "MA 1000 \u0026 2000",
		"Name": "",
		"Revision": "A00",
		"Serial Number": "#09377B7",
		"Type": "Master",
		"index": "cntrl:0"
	}, {
		"Comment": "",
		"DL AGC Attenuation Value": "N/A",
		"DL AGC Status": "N/A",
		"DL DCA Manual Override": "N/A",
		"DL Input Power Status": "N/A",
		"DL Power": "N/A",
		"DL Power Interface Type": "N/A",
		"External Name": "",
		"Firmware Version": "N/A",
		"Name": "BU 4                ",
		"Revision": "N/A",
		"Serial Number": "#FFFFFFF",
		"UL Attenuation Value": "N/A",
		"index": "cntrl:0#p5-base_unit#bu:1-4"
	}],
	"inventory": {
		"cntrl:0": "\"7.2/1.8\"",
		"cntrl:0#p2-riu": "N/A",
		"cntrl:0#p2-riu#p1-btsc": "N/A",
		"cntrl:0#p2-riu#p10-btsc": "N/A",
		"cntrl:0#p2-riu#p11-btsc": "N/A",
		"cntrl:0#p2-riu#p12-btsc": "N/A",
		"cntrl:0#p2-riu#p2-btsc": "N/A",
		"cntrl:0#p2-riu#p2-cmu": "N/A",
		"cntrl:0#p2-riu#p3-btsc": "N/A",
		"cntrl:0#p2-riu#p4-btsc": "N/A",
		"cntrl:0#p2-riu#p5-btsc": "N/A",
		"cntrl:0#p2-riu#p6-btsc": "N/A",
		"cntrl:0#p2-riu#p7-btsc": "N/A",
		"cntrl:0#p2-riu#p8-btsc": "N/A",
		"cntrl:0#p2-riu#p9-btsc": "N/A",
		"cntrl:0#p3-och:unit": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:1": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:2": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:3": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:4": "N/A",
		"cntrl:0#p3-och_unit#och:1-4": "N/A",
		"cntrl:0#p3-och_unit#och:5-8": "N/A",
		"cntrl:0#p4-och:unit": "N/A",
		"cntrl:0#p4-och:unit#och:1-4#gx:1": "N/A",
		"cntrl:0#p4-och_unit#och:1-4": "N/A",
		"cntrl:0#p5-base:unit": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:1": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:2": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:3": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:4": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:1": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:2": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:3": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:4": "N/A",
		"cntrl:0#p6-och:unit": "N/A",
		"cntrl:0#p6-och_unit#och:1-4": "N/A"
	}
}`

	obj := make(map[string]interface{})
	json.Unmarshal([]byte(str), &obj)
	return obj
}

func getCurrentValues() map[string]interface{} {
	str := `{
	"devices": [{
		"Comment": "test5New",
		"External Name": "",
		"Firmware Version": "\"7.2/1.8\"",
		"Model": "MA 1000 \u0026 2000",
		"Name": "",
		"Revision": "A00",
		"Serial Number": "#09377B7",
		"Type": "Master",
		"index": "cntrl:0"
	}, {
		"Comment": "with comment",
		"DL AGC Attenuation Value": "N/A",
		"DL AGC Status": "N/A",
		"DL DCA Manual Override": "N/A",
		"DL Input Power Status": "N/A",
		"DL Power": "N/A",
		"DL Power Interface Type": "N/A",
		"External Name": "",
		"Firmware Version": "N/A",
		"Name": "BU 4                ",
		"Revision": "N/A",
		"Serial Number": "#FFFFFFF",
		"UL Attenuation Value": "N/A",
		"index": "cntrl:0#p5-base_unit#bu:1-4"
	}],
	"inventory": {
		"cntrl:0": "\"7.2/1.9\"",
		"cntrl:0#p2-riu": "N/A",
		"cntrl:0#p2-riu#p1-btsc": "N/A",
		"cntrl:0#p2-riu#p10-btsc": "N/A",
		"cntrl:0#p2-riu#p11-btsc": "N/A",
		"cntrl:0#p2-riu#p12-btsc": "N/A",
		"cntrl:0#p2-riu#p2-btsc": "N/A",
		"cntrl:0#p2-riu#p2-cmu": "N/A",
		"cntrl:0#p2-riu#p3-btsc": "N/A",
		"cntrl:0#p2-riu#p4-btsc": "N/A",
		"cntrl:0#p2-riu#p5-btsc": "N/A",
		"cntrl:0#p2-riu#p6-btsc": "N/A",
		"cntrl:0#p2-riu#p7-btsc": "N/A",
		"cntrl:0#p2-riu#p8-btsc": "N/A",
		"cntrl:0#p2-riu#p9-btsc": "N/A",
		"cntrl:0#p3-och:unit": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:1": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:2": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:3": "N/A",
		"cntrl:0#p3-och:unit#och:1-4#gx:4": "N/A",
		"cntrl:0#p3-och_unit#och:1-4": "N/A",
		"cntrl:0#p3-och_unit#och:5-8": "N/A",
		"cntrl:0#p4-och:unit": "N/A",
		"cntrl:0#p4-och:unit#och:1-4#gx:1": "N/A",
		"cntrl:0#p4-och_unit#och:1-4": "N/A",
		"cntrl:0#p5-base:unit": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:1": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:2": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:3": "N/A",
		"cntrl:0#p5-base_unit#bu:1-4#hx:4": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:1": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:2": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:3": "N/A",
		"cntrl:0#p5-base_unit#bu:5-8#hx:4": "N/A",
		"cntrl:0#p6-och:unit": "N/A",
		"cntrl:0#p6-och_unit#och:1-4": "N/A"
	}
}`
	obj := make(map[string]interface{})
	json.Unmarshal([]byte(str), &obj)
	return obj
}

func TestCalculateDeltas(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	//setup attrs

	prev := getPrevValues()
	current := getCurrentValues()

	diffs := CalculateDeltas(prev, current)

	//check result attr
	if diffs["devices"].([]interface{})[0].(map[string]interface{})["Comment"] != "test5New" {
		t.Fail()
	}
	if diffs["devices"].([]interface{})[1].(map[string]interface{})["Comment"] != "with comment" {
		t.Fail()
	}
	if diffs["inventory"].(map[string]interface{})["cntrl:0"] != "\"7.2/1.9\"" {
		t.Fail()
	}
}
