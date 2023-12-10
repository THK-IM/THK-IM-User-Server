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
	uri := "/user/register"
	url := fmt.Sprintf("%s%s", getUserApiEndpoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	task := httptest.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *httptest.Result {
		startTime := time.Now().UnixMilli()
		registerReq := &dto.RegisterReq{}
		requestJson, _ := json.Marshal(registerReq)
		requestBody := bytes.NewReader(requestJson)
		req, errReq := http.NewRequest("POST", url, requestBody)
		req.Header.Set("Content-Type", contentType)
		if errReq != nil {
			duration := time.Now().UnixMilli() - startTime
			return httptest.NewHttpTestResult(index, -2, 0, duration, errReq)
		}
		response, errHttp := client.Do(req)
		if errHttp != nil {
			duration := time.Now().UnixMilli() - startTime
			return httptest.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			resBytes, err := io.ReadAll(response.Body)
			if err != nil {
				duration := time.Now().UnixMilli() - startTime
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

	for {
		select {
		case success, opened := <-successChan:
			if !opened {
				t.Fail()
			}
			if success {
				t.Skip()
			} else {
				t.Fail()
			}
			return
		}
	}
}
