package sendbuffer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/buntdb"
	"gitlab.com/predictive-open-source/flogo-contrib/activity/sendbuffer/dto"
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

func sendMQTTBatch(broker string, msgs []*dto.DeltaBufferDTO) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	// not make it input
	opts.SetClientID("RANDOM_CLINTID")
	opts.SetUsername("")
	opts.SetPassword("")
	client := mqtt.NewClient(opts)
	fmt.Printf("MQTT Publisher connecting")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("MQTT Publisher connected, sending message")

	for _, msg := range msgs {
		payload := make(map[string]string)
		payload[msg.Attribute] = msg.Value
		payload["timeStamp"] = strconv.FormatInt(msg.TimeStamp, 10)

		token := client.Publish(msg.Topic, byte(0), false, payload)
		token.Wait()
	}

	client.Disconnect(250)
	fmt.Println("MQTT Publisher disconnected")
	return nil
}

func getNewMsgsBuffer(db *buntdb.DB) ([]*dto.DeltaBufferDTO, error) {
	msgs := make([]*dto.DeltaBufferDTO, 0)
	db.CreateIndex("timeStamp", "*", buntdb.IndexJSON("timeStamp"))
	err := db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("timeStamp", func(key, value string) bool {
			// convert currentVal json string DeltaBufferDTO struct
			msgBufferDTO := dto.DeltaBufferDTO{}
			if value != "" {
				err := json.Unmarshal([]byte(value), &msgBufferDTO)
				if err != nil {
					return false
				}
				msgs = append(msgs, &msgBufferDTO)
			}
			msgBufferDTO.Key = key
			return true
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	msgType := context.GetInput("messageType").(string)
	if msgType == "" {
		return false, fmt.Errorf("invalid message type")
	}

	broker := context.GetInput("broker").(string)
	if broker == "" {
		return false, fmt.Errorf("invalid broker")
	}

	db, err := buntdb.Open(fmt.Sprintf("./%v_buffer.db", msgType))
	if err != nil {
		return false, err
	}
	defer db.Close()

	// get current messages buffer
	newMsgsBuffer, err := getNewMsgsBuffer(db)
	if err != nil {
		return false, err
	}

	// send MQTT messages
	sendMQTTBatch(broker, newMsgsBuffer)

	sentNumber := len(newMsgsBuffer)

	context.SetOutput("sentNumber", sentNumber)
	return true, nil
}
