package models

type SendSmsRequest struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	Apicode       string   `json:"apicode"`
	Msisdn        []string `json:"msisdn"` // Can be string or []string
	Countrycode   string   `json:"countrycode"`
	Cli           string   `json:"cli"`
	Messagetype   string   `json:"messagetype"`
	Message       string   `json:"message"`
	Clienttransid string   `json:"clienttransid"`
	BillMsisdn    string   `json:"bill_msisdn"`
	TranType      string   `json:"tran_type"`
	RequestType   string   `json:"request_type"`
	RnCode        string   `json:"rn_code"`
}

type CheckCliRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Apicode       string `json:"apicode"`
	Cli           string `json:"cli"`
	Clienttransid string `json:"clienttransid"`
}

type CheckCreditRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Apicode       string `json:"apicode"`
	Clienttransid string `json:"clienttransid"`
}

type CheckDlrRequest struct {
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	Apicode         string   `json:"apicode"`
	Msisdn          []string `json:"msisdn"`
	Countrycode     string   `json:"countrycode"`
	Cli             string   `json:"cli"`
	Messagetype     string   `json:"messagetype"`
	Clienttransid   string   `json:"clienttransid"`
	Operatortransid string   `json:"operatortransid"`
}

type StatusInfo struct {
	StatusCode          string `json:"statusCode"`
	ErrorDescription    string `json:"errordescription"`
	ClientTransID       string `json:"clienttransid"`
	ServerReferenceCode string `json:"serverReferenceCode"`
}

type SendSmsResponse struct {
	StatusInfo StatusInfo `json:"statusInfo"`
}
