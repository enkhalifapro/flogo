package calculatedelta

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"
	"gitlab.com/predictive-open-source/flogo-contrib/activity/calculatedelta/dto"
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

func getNewMsgs(db *buntdb.DB) ([]*dto.MsgDTO, error) {
	msgs := make([]*dto.MsgDTO, 0)
	db.CreateIndex("timeStamp", "*", buntdb.IndexJSON("timeStamp"))
	err := db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("timeStamp", func(key, value string) bool {
			// convert currentVal json string to map[string]interface
			msgDTO := dto.MsgDTO{}
			if value != "" {
				err := json.Unmarshal([]byte(value), &msgDTO)
				if err != nil {
					return false
				}
				msgs = append(msgs, &msgDTO)
			}
			return true
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func updateDeltas(msgs []*dto.MsgDTO) ([]*dto.MsgDTO, error) {
	db, err := buntdb.Open("./gateway.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// get changed values
	changedMsgs := make([]*dto.MsgDTO, 0)
	err = db.View(func(tx *buntdb.Tx) error {
		for _, msg := range msgs {
			key := fmt.Sprintf("%v_%v_delta", msg.ID, msg.Attribute)
			lastVal, err := tx.Get(key)
			if err != nil && err.Error() == "not found" {
				changedMsgs = append(changedMsgs, msg)
				continue
			}
			if err != nil {
				return err
			}
			if lastVal != msg.Value {
				changedMsgs = append(changedMsgs, msg)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// update changed messages Deltas
	err = db.Update(func(tx *buntdb.Tx) error {
		for _, msg := range changedMsgs {
			key := fmt.Sprintf("%v_%v_delta", msg.ID, msg.Attribute)
			_, _, err := tx.Set(key, msg.Value, nil)

			if err != nil {
				return err
			}
		}
		return nil
	})

	return changedMsgs, nil
}

func updateDeltasBuffer(msgType string, msgs []*dto.MsgDTO) error {
	db, err := buntdb.Open(fmt.Sprintf("./%v_buffer.db", msgType))
	if err != nil {
		return err
	}
	defer db.Close()

	// save deltas buffer
	err = db.Update(func(tx *buntdb.Tx) error {
		for _, msg := range msgs {
			key := fmt.Sprintf("%v_%v_delta", msg.ID, msg.Attribute)
			jsonVal, err := json.Marshal(&dto.DeltaBufferDTO{
				ID:        msg.ID,
				Topic:     msg.Topic,
				Attribute: msg.Attribute,
				Value:     msg.Value,
				TimeStamp: msg.TimeStamp,
			})
			if err != nil {
				return err
			}
			_, _, err = tx.Set(key, string(jsonVal), nil)
			if err != nil {
				return err
			}
		}
		return err
	})

	return err
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	// do eval
	msgType := context.GetInput("messageType").(string)
	if msgType == "" {
		return false, fmt.Errorf("invalid message type")
	}

	db, err := buntdb.Open(fmt.Sprintf("./%v.db", msgType))
	if err != nil {
		return false, err
	}
	defer db.Close()

	// process deltas of current msgQ messages
	newMsgs, err := getNewMsgs(db)
	if err != nil {
		return false, err
	}

	// update changed values
	changedMsgs, err := updateDeltas(newMsgs)
	if err != nil {
		return false, err
	}

	// update delta buffer
	err = updateDeltasBuffer(msgType, changedMsgs)
	if err != nil {
		return false, err
	}

	context.SetOutput("deltasNumber", len(changedMsgs))
	return true, nil
}
