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

var RESENT_API_KEY = os.Getenv("RESENT_API_KEY")
var RESENT_API_URL = os.Getenv("RESENT_API_URL")
var RESENT_FROM_EMAIL = os.Getenv("RESENT_FROM_EMAIL")

func SendIt(sendItContent TResendBody) int {
	bodyBytes, _ := json.Marshal(sendItContent)
	req, _ := http.NewRequest(
		"POST",
		RESENT_API_URL,
		bytes.NewBuffer(bodyBytes))

	req.Header.Set("Authorization", ALLOWED_API_KEY)
	sendItContent.From = RESENT_FROM_EMAIL

	res, _ := http.DefaultClient.Do(req)

	return res.StatusCode
}
