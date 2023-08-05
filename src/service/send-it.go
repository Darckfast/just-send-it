package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type TResendBody struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

var RESEND_API_KEY = os.Getenv("RESEND_API_KEY")
var RESEND_API_URL = os.Getenv("RESEND_API_URL")
var RESEND_FROM_EMAIL = os.Getenv("RESEND_FROM_EMAIL")

func SendIt(sendItContent *TResendBody) *http.Response {
	sendItContent.From = RESEND_FROM_EMAIL
	bodyBytes, _ := json.Marshal(sendItContent)

	req, _ := http.NewRequest(
		"POST",
		RESEND_API_URL,
		bytes.NewBuffer(bodyBytes))

	req.Header.Set("Authorization", RESEND_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	return res
}
