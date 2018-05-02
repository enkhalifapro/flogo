package dto

type DeltaBufferDTO struct {
	Key       string `json:"-"`  // db key
	ID        string `json:"id"` // deviceID
	Topic     string `json:"topic"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
	TimeStamp int64  `json:"timeStamp"`
}
