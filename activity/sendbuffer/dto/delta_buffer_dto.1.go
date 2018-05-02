package dto

type MQTTMsgDTO struct {
	Topic   string            `json:"topic"` // db key
	Payload map[string]string `json:"payload"`
}
