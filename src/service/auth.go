package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var API_URL = os.Getenv("KV_REST_API_URL")
var API_TOKEN = os.Getenv("KV_REST_API_TOKEN")
var SIGN_SECRET = os.Getenv("SIGN_SECRET")
var ALLOWED_API_KEY = os.Getenv("ALLOWED_API_KEY")

type TSessionResponse struct {
	Result string
}

func checkIfExist(sessionKey string) string {
	client := &http.Client{}

	sessionBody := []string{
		"HGET", sessionKey, "token",
	}

	bodyBytes, _ := json.Marshal(sessionBody)

	req, _ := http.NewRequest("POST", API_URL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var sessionResponse TSessionResponse

	decoder.Decode(&sessionResponse)

	return sessionResponse.Result
}

func RemoveSession(sessionKey string) {
	if sessionKey == "" {
		return
	}

	client := &http.Client{}

	deleteBody := []string{
		"DEL", sessionKey,
	}

	bodyBytes, _ := json.Marshal(deleteBody)

	req, _ := http.NewRequest("POST", API_URL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	client.Do(req)
}

func getJWTFromToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(SIGN_SECRET), nil
	})
}

func ValidateSession(rawSessionString string) bool {
	sessionString := strings.Replace(rawSessionString, "Bearer ", "", 1)

	token, err := getJWTFromToken(sessionString)

	if err != nil {
		fmt.Println("invalid token", err)
		return false
	}

	sessionKey, _ := token.Claims.GetSubject()

	storedSession := checkIfExist("session:" + sessionKey)
	if storedSession != "" && storedSession == sessionString {
		return true
	}

	return false
}

func CreateSession(apiKey string) string {
	hash := sha256.New()
	hash.Write([]byte(apiKey))

	sessionHash := hex.EncodeToString(hash.Sum(nil))
	sessionKey := "session:" + sessionHash
	storedSession := checkIfExist(sessionKey)
	isStoredSessionValid, _ := getJWTFromToken(storedSession)

	if storedSession != "" && isStoredSessionValid.Valid {
		return storedSession
	}

	RemoveSession(storedSession)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   sessionHash,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "just-send-it",
	})

	tokenString, _ := token.SignedString([]byte(SIGN_SECRET))

	client := &http.Client{}
	sessionBody := [][]string{
		{"HSET", sessionKey, "token", tokenString},
		{"EXPIRE", sessionKey, "3600"},
	}

	bodyBytes, _ := json.Marshal(sessionBody)
	req, _ := http.NewRequest("POST", API_URL+"/multi-exec", bytes.NewBuffer(bodyBytes))

	req.Header.Set("Authorization", "Bearer "+API_TOKEN)

	client.Do(req)

	return tokenString
}
