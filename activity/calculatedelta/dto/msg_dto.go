package dto

type MsgDTO struct {
	Type      string `json:"type"`
	Topic     string `json:"topic"`
	ID        string `json:"id"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	TimeStamp int64  `json:"timeStamp"`
	KeyHash   string `json:"-"` // to generate hash once
}
