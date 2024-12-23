package models

type Message struct {
	Source        string
	Destination   string
	Body          string
	TransactionId int64
	DataCoding    int16
}
