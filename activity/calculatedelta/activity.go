package calculatedelta

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/tidwall/buntdb"
	"gitlab.com/predictive-open-source/flogo-contrib/activity/calculatedelta/dto"
	"gitlab.com/predictive-open-source/flogo-contrib/activity/calculatedelta/hash"
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

func getNewMsgs(db *buntdb.DB) ([]string, []*dto.MsgDTO, error) {
	msgs := make([]*dto.MsgDTO, 0)
	db.CreateIndex("timeStamp", "*", buntdb.IndexJSON("timeStamp"))
	IDs := make([]string, 0)
	err := db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("timeStamp", func(key, value string) bool {
			IDs = append(IDs, key)
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
		return nil, nil, err
	}
	return IDs, msgs, nil
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
		fnvHasher := hash.NewFNV1aHelper()
		for _, msg := range msgs {
			keyHash := fnvHasher.GetHashString([]byte(fmt.Sprintf("%v_%v_delta", msg.ID, msg.Attribute)))
			msg.KeyHash = keyHash
			lastVal, err := tx.Get(keyHash)
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
			fmt.Println(msg.KeyHash)
			_, _, err := tx.Set(msg.KeyHash, msg.Value, nil)

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
	fnvHasher := hash.NewFNV1aHelper()
	err = db.Update(func(tx *buntdb.Tx) error {
		for _, msg := range msgs {
			keyHash := fnvHasher.GetHashString([]byte(fmt.Sprintf("%v_%v_delta", msg.ID, msg.Attribute)))
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
			_, _, err = tx.Set(keyHash, string(jsonVal), nil)
			if err != nil {
				return err
			}
		}
		return err
	})

	return err
}

// NOTE: don't delete all update it to delete processed messages only
func clearMsgQ(db *buntdb.DB, processedIDs []string) error {
	err := db.Update(func(tx *buntdb.Tx) error {
		for _, id := range processedIDs {
			_, err := tx.Delete(id)
			if err != nil {
				return err
			}
		}
		return nil
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
	IDs, newMsgs, err := getNewMsgs(db)
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

	// clear processed msg IDs
	err = clearMsgQ(db, IDs)
	if err != nil {
		return false, err
	}

	context.SetOutput("deltasNumber", len(changedMsgs))
	return true, nil
}
