package entities

type User struct {
	Id                         uint64
	Username, Password, MtPort string
	Status                     int
	MidExpiryTime              float32
	DcApi, CpApi, SpApi        string
}
