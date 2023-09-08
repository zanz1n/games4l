package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type UserCreateRequest struct {
	Body       []byte
	Duration   time.Duration
	StatusCode int
}

var mailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func MakeUserCreateRequest(
	body map[string]any,
	sig string,
	baseUrl string,
) (*UserCreateRequest, error) {
	payload, err := json.Marshal(body)

	digest := sha256.New()
	if _, err = digest.Write(payload); err != nil {
		return nil, errors.New("something went wrong with hashing algorithms: " + err.Error())
	}
	if _, err = digest.Write([]byte(sig)); err != nil {
		return nil, errors.New("something went wrong with hashing algorithms: " + err.Error())
	}

	sumHex := hex.EncodeToString(digest.Sum([]byte{}))
	start := time.Now()

	req, err := http.NewRequest(
		"POST",
		baseUrl+"/user",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, errors.New("failed to build the request entity: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Signature HEX "+sumHex)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to make the request: " + err.Error())
	}
	defer res.Body.Close()

	resBuf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read the response body: " + err.Error())
	}

	return &UserCreateRequest{
		Body:       identJson(resBuf),
		Duration:   time.Since(start),
		StatusCode: res.StatusCode,
	}, nil
}

func HandleUserCreate(
	args map[string]string,
	baseUrl string,
	webhookSig string,
) error {
	var (
		ok       bool
		username string
		email    string
		password string
		role     string
	)

	if username, ok = args["username"]; !ok {
		return errors.New("username argument must be provided")
	}
	if email, ok = args["email"]; !ok {
		return errors.New("email argument must be provided")
	}
	if password, ok = args["password"]; !ok {
		return errors.New("password argument must be provided")
	}
	if role, ok = args["role"]; !ok {
		return errors.New("role argument must be provided")
	}

	if role != "PACIENT" && role != "ADMIN" && role != "CLIENT" {
		return errors.New("role argument must be 'PACIENT' | 'ADMIN' | 'CLIENT'")
	}
	if !mailRegex.MatchString(email) {
		return errors.New("the provided email address is not valid")
	}

	reqBodyMap := JSON{
		"username": username,
		"email":    email,
		"password": password,
		"role":     role,
	}

	result, err := MakeUserCreateRequest(reqBodyMap, webhookSig, baseUrl)
	if err == nil {
		if result.StatusCode >= 200 && result.StatusCode < 300 {
			fmt.Println("User created: " + string(result.Body))
		} else {
			fmt.Println("User creation failed: " + string(result.Body))
		}
	} else {
		return err
	}

	return nil
}
