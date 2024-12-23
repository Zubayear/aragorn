package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/Zubayear/aragorn/api/presenter"
	"github.com/Zubayear/aragorn/pkg/models"
	"github.com/Zubayear/aragorn/pkg/sms"
	"github.com/Zubayear/aragorn/pkg/user"
	"github.com/goccy/go-json"
	"github.com/rs/xid"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CpHandler struct {
	userRepo user.Repository
}

func NewCpHandler(userRepo user.Repository) *CpHandler {
	return &CpHandler{userRepo: userRepo}
}

func (ch *CpHandler) HandleRequest(userService user.Service, smsService sms.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			// w.Header().Set("Content-Type", "application/json")
			response := presenter.SendSmsResponse("1011", "Something went wrong. try again", "", "")
			http.Error(w, response, 200)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				response := presenter.SendSmsResponse(
					"1011",
					"Something went wrong. try again",
					"",
					"",
				)
				http.Error(w, response, 200)
				return
			}
		}(r.Body)
		var reqBody models.SendSmsRequest

		err = json.Unmarshal(bytes, &reqBody)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				reqBody.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}

		switch reqBody.Apicode {
		case "5":
			handleApiCode5(r.Context(), w, reqBody, userService, smsService)
		case "6":
			handleApiCode5(r.Context(), w, reqBody, userService, smsService)
		case "4":
			handleApiCode4(r.Context(), w, reqBody, userService, smsService)
		case "3":
			handleApiCode3(r.Context(), w, reqBody, userService, smsService)
		default:
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				reqBody.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
	}
}

func handleApiCode3(
	ctx context.Context,
	w http.ResponseWriter,
	body models.SendSmsRequest,
	userService user.Service,
	smsService sms.Service,
) {
	if !(len(body.Clienttransid) >= 10 && len(body.Clienttransid) <= 36 && regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(body.Clienttransid)) {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType != "T" && body.TranType != "P" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype != "1" && body.Messagetype != "3" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "1" && len(body.Message) > 1600 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "3" && len(body.Message) > 700 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.RequestType == "S" && len(body.Msisdn) > 1 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "T" && body.RequestType == "B" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && body.Messagetype == "1" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && len(body.Msisdn) > 1000 {
		response := presenter.SendSmsResponse(
			"1054",
			"MSISDN Limit Exceeded",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "88") && len(body.Cli) != 13 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "01") && len(body.Cli) != 11 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	userFromCache, err := userService.GetUserFromCache(ctx, "user:"+body.Username)

	if err.Error() == "user not found" {
		userFromDb, err := userService.GetUser(ctx, body.Username)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1002",
				"Invalid Username",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromDb.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse(
				"1003",
				"Invalid Password",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromDb.MtPort, body.Cli) {
			response := presenter.SendSmsResponse(
				"1006",
				"CLI/Masking Invalid",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if userFromDb.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromDb.C15CampaignCpapi == "false" {
			response := presenter.SendSmsResponse(
				"1017",
				"API Not allowed for user",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		value, err := json.Marshal(userFromDb)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		isSaved := userService.SaveToCache(ctx, "user:"+body.Username, value)
		if !isSaved {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
	} else {
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromCache.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse("1003", "Invalid Password", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromCache.MtPort, body.Cli) {
			response := presenter.SendSmsResponse("1006", "CLI/Masking Invalid", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.CpApi == "false" {
			response := presenter.SendSmsResponse("1017", "API Not allowed for user", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	duplicateTranKey := body.Username + ":" + body.Clienttransid

	if userService.KeyExists(ctx, duplicateTranKey) {
		response := presenter.SendSmsResponse(
			"1005",
			"Duplicate Transaction Id",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	} else {
		isSaved := userService.SaveToCache(ctx, duplicateTranKey, []byte(""))
		if !isSaved {
			response := presenter.SendSmsResponse("1011", "Something went wrong. try again", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	// call stored procedure
	numbers := strings.Join(body.Msisdn, ",")
	refId := xid.New()
	apiCode, _ := strconv.Atoi(body.Apicode)
	var msgType int
	if body.TranType == "T" {
		msgType = 1
	} else {
		msgType = 2
	}
	var unicode int
	if body.Messagetype == "3" {
		unicode = 1
	} else {
		unicode = 0
	}
	var msgLen int
	if body.Messagetype == "3" {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 70.0))
	} else {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 160.0))
	}
	result := smsService.ProcessNumbersCp(
		ctx,
		body.Cli,
		body.Username,
		"",
		numbers,
		apiCode,
		msgType,
		unicode,
		msgLen,
		refId.String(),
		body.Message,
		refId.String(),
	)

	var spResponse models.CpSpResponse
	_ = json.Unmarshal([]byte(result), &spResponse)

	if spResponse.Message == "no_valid_numbers" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}
	if spResponse.Message == "not_enough_credit" {
		response := presenter.SendSmsResponse(
			"1008",
			"Insufficient Balance",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}
	var messages []string
	if spResponse.OffnetNumbers != "" {
		off := strings.Split(spResponse.OffnetNumbers, ",")
		for _, v := range off {
			message := models.Message{
				Source:        spResponse.Longcode,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}
	if spResponse.OnnetNumbers != "" {
		onn := strings.Split(spResponse.OnnetNumbers, ",")
		for _, v := range onn {
			message := models.Message{
				Source:        body.Cli,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}

	refCode := strconv.FormatInt(spResponse.RefCode, 10)
	response := presenter.SendSmsResponse("1000", "Success", body.Clienttransid, refCode)
	http.Error(w, response, 200)
	return
}

func handleApiCode4(
	ctx context.Context,
	w http.ResponseWriter,
	body models.SendSmsRequest,
	userService user.Service,
	smsService sms.Service,
) {
	// ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	if !(len(body.Clienttransid) >= 10 && len(body.Clienttransid) <= 36 && regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(body.Clienttransid)) {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType != "T" && body.TranType != "P" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype != "1" && body.Messagetype != "3" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "1" && len(body.Message) > 1600 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "3" && len(body.Message) > 700 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.RequestType == "S" && len(body.Msisdn) > 1 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "T" && body.RequestType == "B" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && body.Messagetype == "1" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && len(body.Msisdn) > 1000 {
		response := presenter.SendSmsResponse(
			"1054",
			"MSISDN Limit Exceeded",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "88") && len(body.Cli) != 13 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "01") && len(body.Cli) != 11 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	// var dbResult, redisResult
	// credential check
	userFromCache, err := userService.GetUserFromCache(ctx, "user:"+body.Username)

	if err.Error() == "user not found" {
		userFromDb, err := userService.GetUser(ctx, body.Username)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1002",
				"Invalid Username",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromDb.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse(
				"1003",
				"Invalid Password",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromDb.MtPort, body.Cli) {
			response := presenter.SendSmsResponse(
				"1006",
				"CLI/Masking Invalid",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if userFromDb.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromDb.C15CampaignCpapi == "false" {
			response := presenter.SendSmsResponse(
				"1017",
				"API Not allowed for user",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		value, err := json.Marshal(userFromDb)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		isSaved := userService.SaveToCache(ctx, "user:"+body.Username, value)
		if !isSaved {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
	} else {
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromCache.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse("1003", "Invalid Password", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromCache.MtPort, body.Cli) {
			response := presenter.SendSmsResponse("1006", "CLI/Masking Invalid", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.CpApi == "false" {
			response := presenter.SendSmsResponse("1017", "API Not allowed for user", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	duplicateTranKey := body.Username + ":" + body.Clienttransid

	if userService.KeyExists(ctx, duplicateTranKey) {
		response := presenter.SendSmsResponse(
			"1005",
			"Duplicate Transaction Id",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	} else {
		isSaved := userService.SaveToCache(ctx, duplicateTranKey, []byte(""))
		if !isSaved {
			response := presenter.SendSmsResponse("1011", "Something went wrong. try again", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	// call stored procedure
	numbers := strings.Join(body.Msisdn, ",")
	refId := xid.New()
	apiCode, _ := strconv.Atoi(body.Apicode)
	var msgType int
	if body.TranType == "T" {
		msgType = 1
	} else {
		msgType = 2
	}
	var unicode int
	if body.Messagetype == "3" {
		unicode = 1
	} else {
		unicode = 0
	}
	var msgLen int
	if body.Messagetype == "3" {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 70.0))
	} else {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 160.0))
	}
	result := smsService.ProcessNumbersCp(
		ctx,
		body.Cli,
		body.Username,
		"",
		numbers,
		apiCode,
		msgType,
		unicode,
		msgLen,
		refId.String(),
		body.Message,
		refId.String(),
	)

	var spResponse models.CpSpResponse
	_ = json.Unmarshal([]byte(result), &spResponse)

	if spResponse.Message == "no_valid_numbers" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}
	if spResponse.Message == "not_enough_credit" {
		response := presenter.SendSmsResponse(
			"1008",
			"Insufficient Balance",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}
	var messages []string
	if spResponse.OffnetNumbers != "" {
		off := strings.Split(spResponse.OffnetNumbers, ",")
		for _, v := range off {
			message := models.Message{
				Source:        spResponse.Longcode,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}
	if spResponse.OnnetNumbers != "" {
		onn := strings.Split(spResponse.OnnetNumbers, ",")
		for _, v := range onn {
			message := models.Message{
				Source:        body.Cli,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}

	refCode := strconv.FormatInt(spResponse.RefCode, 10)
	response := presenter.SendSmsResponse("1000", "Success", body.Clienttransid, refCode)
	http.Error(w, response, 200)
	return
}

func handleApiCode5(
	ctx context.Context,
	w http.ResponseWriter,
	body models.SendSmsRequest,
	userService user.Service,
	smsService sms.Service,
) {
	// ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	if !(len(body.Clienttransid) >= 10 && len(body.Clienttransid) <= 36 && regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(body.Clienttransid)) {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType != "T" && body.TranType != "P" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype != "1" && body.Messagetype != "3" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "1" && len(body.Message) > 1600 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.Messagetype == "3" && len(body.Message) > 700 {
		response := presenter.SendSmsResponse(
			"1012",
			"Message Length Exceeds",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if body.RequestType == "S" && len(body.Msisdn) > 1 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "T" && body.RequestType == "B" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && body.Messagetype == "1" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if body.TranType == "P" && len(body.Msisdn) > 1000 {
		response := presenter.SendSmsResponse(
			"1054",
			"MSISDN Limit Exceeded",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "88") && len(body.Cli) != 13 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	if strings.HasPrefix(body.Cli, "01") && len(body.Cli) != 11 {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}

	// var dbResult, redisResult
	// credential check
	userFromCache, err := userService.GetUserFromCache(ctx, "user:"+body.Username)

	if err.Error() == "user not found" {
		userFromDb, err := userService.GetUser(ctx, body.Username)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1002",
				"Invalid Username",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromDb.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse(
				"1003",
				"Invalid Password",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromDb.MtPort, body.Cli) {
			response := presenter.SendSmsResponse(
				"1006",
				"CLI/Masking Invalid",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		if userFromDb.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromDb.C15CampaignCpapi == "false" {
			response := presenter.SendSmsResponse(
				"1017",
				"API Not allowed for user",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		value, err := json.Marshal(userFromDb)
		if err != nil {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
		isSaved := userService.SaveToCache(ctx, "user:"+body.Username, value)
		if !isSaved {
			response := presenter.SendSmsResponse(
				"1011",
				"Something went wrong. try again",
				body.Clienttransid,
				"",
			)
			http.Error(w, response, 200)
			return
		}
	} else {
		rawPasswordHash := hashMD5(body.Password)
		if !compareMD5Hashes(userFromCache.Password, rawPasswordHash) {
			response := presenter.SendSmsResponse("1003", "Invalid Password", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if !strings.Contains(userFromCache.MtPort, body.Cli) {
			response := presenter.SendSmsResponse("1006", "CLI/Masking Invalid", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.Status == -1 {
			response := presenter.SendSmsResponse("1007", "Account Barred", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
		if userFromCache.CpApi == "false" {
			response := presenter.SendSmsResponse("1017", "API Not allowed for user", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	duplicateTranKey := body.Username + ":" + body.Clienttransid

	if userService.KeyExists(ctx, duplicateTranKey) {
		response := presenter.SendSmsResponse(
			"1005",
			"Duplicate Transaction Id",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	} else {
		isSaved := userService.SaveToCache(ctx, duplicateTranKey, []byte(""))
		if !isSaved {
			response := presenter.SendSmsResponse("1011", "Something went wrong. try again", body.Clienttransid, "")
			http.Error(w, response, 200)
			return
		}
	}

	// call stored procedure
	numbers := strings.Join(body.Msisdn, ",")
	refId := xid.New()
	apiCode, _ := strconv.Atoi(body.Apicode)
	var msgType int
	if body.TranType == "T" {
		msgType = 1
	} else {
		msgType = 2
	}
	var unicode int
	if body.Messagetype == "3" {
		unicode = 1
	} else {
		unicode = 0
	}
	var msgLen int
	if body.Messagetype == "3" {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 70.0))
	} else {
		msgLen = int(math.Ceil(float64(len(body.Message)) / 160.0))
	}
	result := smsService.ProcessNumbersCp(
		ctx,
		body.Cli,
		body.Username,
		"",
		numbers,
		apiCode,
		msgType,
		unicode,
		msgLen,
		refId.String(),
		body.Message,
		refId.String(),
	)

	var spResponse models.CpSpResponse
	_ = json.Unmarshal([]byte(result), &spResponse)

	if spResponse.Message == "no_valid_numbers" {
		response := presenter.SendSmsResponse("1005", "Invalid Parameter", body.Clienttransid, "")
		http.Error(w, response, 200)
		return
	}
	if spResponse.Message == "not_enough_credit" {
		response := presenter.SendSmsResponse(
			"1008",
			"Insufficient Balance",
			body.Clienttransid,
			"",
		)
		http.Error(w, response, 200)
		return
	}
	var messages []string
	if spResponse.OffnetNumbers != "" {
		off := strings.Split(spResponse.OffnetNumbers, ",")
		for _, v := range off {
			message := models.Message{
				Source:        spResponse.Longcode,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}
	if spResponse.OnnetNumbers != "" {
		onn := strings.Split(spResponse.OnnetNumbers, ",")
		for _, v := range onn {
			message := models.Message{
				Source:        body.Cli,
				Destination:   v,
				Body:          body.Message,
				TransactionId: spResponse.RefCode,
				DataCoding:    3,
			}
			b, _ := json.Marshal(message)
			messages = append(messages, string(b))
		}
	}

	refCode := strconv.FormatInt(spResponse.RefCode, 10)
	response := presenter.SendSmsResponse("1000", "Success", body.Clienttransid, refCode)
	http.Error(w, response, 200)
	return
}

// hashMD5 hashes a string using MD5 and returns the hash as a hexadecimal string
func hashMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// compareMD5Hashes compares two MD5 hashes and returns true if they match, false otherwise
func compareMD5Hashes(hashedPassword, rawPasswordHash string) bool {
	return hashedPassword == rawPasswordHash
}
