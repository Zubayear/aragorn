package presenter

import (
	"github.com/Zubayear/aragorn/pkg/models"
	"github.com/goccy/go-json"
)

func SendSmsResponse(statusCode, errorDescription, clientTransactionId, referenceCode string) string {
	m := &models.SendSmsResponse{
		StatusInfo: models.StatusInfo{
			StatusCode:          statusCode,
			ErrorDescription:    errorDescription,
			ClientTransID:       clientTransactionId,
			ServerReferenceCode: referenceCode,
		},
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bytes)
}
