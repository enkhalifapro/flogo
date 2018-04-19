package puppq

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/go-mangos/mangos/protocol/push"
	"github.com/go-mangos/mangos/transport/all"
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

	// get message input
	message := context.GetInput("message").(map[string]interface{})
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return false, err
	}

	// push message
	err = pushMsg(jsonMessage)
	if err != nil {
		return false, err
	}
	context.SetOutput("result", true)

	return true, nil
}

func pushMsg(msg []byte) error {
	url := "tcp://127.0.0.1:7000"

	// Create Push Client
	pushSocket, err := push.NewSocket()
	if err != nil {
		return err
	}
	defer pushSocket.Close()
	all.AddTransports(pushSocket)

	// Client dials Server
	if err = pushSocket.Dial(url); err != nil {
		fmt.Printf("\nClient dial failed: %v", err)
		return err
	}

	// Client sends message
	if err = pushSocket.Send(msg); err != nil {
		fmt.Printf("\nClient send failed: %v", err)
		return err
	}
	fmt.Printf("\nClient sending: %s\n", msg)
	return nil
}
