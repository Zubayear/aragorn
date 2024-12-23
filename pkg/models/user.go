package models

type User struct {
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	MtPort        string  `json:"mtPort"`
	Status        int     `json:"status"`
	MidExpiryTime float32 `json:"midExpiryTime"`
	DcApi         string  `json:"dcApi"`
	CpApi         string  `json:"cpApi"`
	SpApi         string  `json:"spApi"`
}
