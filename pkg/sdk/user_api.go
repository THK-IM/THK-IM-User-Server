package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
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
		LoginByToken(req dto.TokenLoginReq, claims baseDto.ThkClaims) (*dto.LoginRes, error)
	}

	defaultUserApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultUserApi) LoginByToken(claims baseDto.ThkClaims) (*dto.LoginRes, error) {
	url := fmt.Sprintf("%s%s", d.endpoint, tokenLoginPath)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", contentType).
		Post(url)
	if errRequest != nil {
		d.logger.Errorf("LoginByToken %v %v", req, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := &errorx.ErrorX{}
		e := json.Unmarshal(res.Body(), errRes)
		if e != nil {
			d.logger.Errorf("LoginByToken: %v %v", req, e)
			return nil, e
		} else {
			return nil, errRes
		}
	} else {
		resp := &dto.LoginRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("LoginByToken: %v %v", req, e)
			return nil, e
		} else {
			d.logger.Infof("LoginByToken: %v %v", req, resp)
			return resp, nil
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
