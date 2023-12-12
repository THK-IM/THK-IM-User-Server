package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"net/http"
	"time"
)

const (
	tokenLoginPath string = "/user/login/token"
	contentType    string = "application/json"
)

type (
	UserApi interface {
		LoginByToken(req dto.TokenLoginReq) (*dto.LoginRes, error)
	}

	defaultUserApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultUserApi) LoginByToken(req dto.TokenLoginReq) (*dto.LoginRes, error) {
	dataBytes, errJson := json.Marshal(req)
	if errJson != nil {
		d.logger.Error(errJson)
		return nil, errJson
	}
	url := fmt.Sprintf("%s%s", d.endpoint, tokenLoginPath)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", contentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		d.logger.Error(errRequest, string(dataBytes))
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Error(e, string(dataBytes))
		return nil, e
	} else {
		loginRes := &dto.LoginRes{}
		errJson = json.Unmarshal(res.Body(), loginRes)
		if errJson != nil {
			d.logger.Error(errJson, string(res.Body()))
			return nil, errJson
		} else {
			return loginRes, nil
		}
	}
}

func NewUserApi(sdk conf.Sdk, logger *logrus.Entry) UserApi {
	return defaultUserApi{
		endpoint: sdk.Endpoint,
		logger:   logger.WithField("rpc", sdk.Name),
		client: resty.New().
			SetTransport(&http.Transport{
				MaxIdleConns:    10,
				MaxConnsPerHost: 10,
				IdleConnTimeout: 30 * time.Second,
			}).
			SetTimeout(5 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(15 * time.Second).
			SetRetryMaxWaitTime(5 * time.Second),
	}
}
