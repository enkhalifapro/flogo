package dto

type DeltaBufferDTO struct {
	ID        string `json:"id"` // deviceID
	Topic     string `json:"topic"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	TimeStamp int64  `json:"timeStamp"`
}
