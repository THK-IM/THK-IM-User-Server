package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thk-im/thk-im-base-server/httptest"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"io"
	"net/http"
	"testing"
	"time"
)

func getUserApiEndpoint() string {
	return "http://localhost:10000"
}

func TestUserRegister(t *testing.T) {
	now := time.Now().UnixMilli()
	uri := "/user/register"
	url := fmt.Sprintf("%s%s", getUserApiEndpoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	accounts := make([]*string, 0)
	passwords := make([]*string, 0)
	for i := 0; i < count; i++ {
		var account, password *string
		if i%3 != 0 {
			account = nil
			password = nil
		} else {
			act := fmt.Sprintf("act-%d-%d", now%10000, i)
			pwd := fmt.Sprintf("p-%d", i)
			account = &act
			password = &pwd
		}
		accounts = append(accounts, account)
		passwords = append(passwords, password)
	}
	task := httptest.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *httptest.Result {
		startTime := time.Now().UnixMilli()
		registerReq := &dto.RegisterReq{
			Account:  accounts[index],
			Password: passwords[index],
		}
		requestJson, _ := json.Marshal(registerReq)
		requestBody := bytes.NewReader(requestJson)
		req, errReq := http.NewRequest("POST", url, requestBody)
		req.Header.Set("Content-Type", contentType)
		if errReq != nil {
			duration := time.Now().UnixMilli() - startTime
			return httptest.NewHttpTestResult(index, -2, 0, duration, errReq)
		}
		response, errHttp := client.Do(req)
		duration := time.Now().UnixMilli() - startTime
		if errHttp != nil {
			return httptest.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			if response.StatusCode >= 400 {
				return httptest.NewHttpTestResult(index, response.StatusCode, 0, duration, errHttp)
			} else {
				resBytes, err := io.ReadAll(response.Body)
				if err != nil {
					return httptest.NewHttpTestResult(index, 500, 0, duration, errHttp)
				}
				registerResp := &dto.RegisterRes{}
				err = json.Unmarshal(resBytes, registerResp)
				duration := time.Now().UnixMilli() - startTime
				if err != nil {
					return httptest.NewHttpTestResult(index, 500, 0, duration, errHttp)
				} else {
					return httptest.NewHttpTestResult(index, 200, int64(len(resBytes)), duration, nil)
				}
			}

		}
	}, func(task *httptest.Task) {
		task.PrintResults()
		for _, result := range task.Results() {
			if result.StatusCode() != http.StatusOK {
				successChan <- false
				return
			}
		}
		successChan <- true
		return
	})
	task.Start()

	responseCnt := 0
	responseSuccessCnt := 0
	for {
		select {
		case success, _ := <-successChan:
			responseCnt++
			if success {
				responseSuccessCnt++
			}
			if responseCnt == count {
				if responseCnt == responseSuccessCnt {
					t.Skip()
				} else {
					t.Fail()
				}
			}
			return
		}
	}
}
