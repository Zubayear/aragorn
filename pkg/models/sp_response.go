package models

type CpSpResponse struct {
	Message       string `json:"msg"`
	OnnetNumbers  string `json:"onnetNumbers"`
	Longcode      string `json:"lc"`
	RefCode       int64  `json:"rc"`
	OffnetNumbers string `json:"offnetNumbers"`
	RequestId     string `json:"reqId"`
}
