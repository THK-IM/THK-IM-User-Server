package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErrorx "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"net/http"
	"time"
)

const (
	tokenLoginPath string = "/login/token"
)

type (
	LoginApi interface {
		LoginByToken(claims baseDto.ThkClaims) (*dto.LoginRes, error)
	}

	defaultLoginApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultLoginApi) LoginByToken(claims baseDto.ThkClaims) (*dto.LoginRes, error) {
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
		d.logger.Errorf("LoginByToken %v %v", claims, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := baseErrorx.NewErrorXFromResp(res)
		d.logger.Errorf("LoginByToken: %v %v", claims, errRes)
		return nil, errRes
	} else {
		resp := &dto.LoginRes{}
		e := json.Unmarshal(res.Body(), resp)
		if e != nil {
			d.logger.Errorf("LoginByToken: %v %v", claims, e)
			return nil, e
		} else {
			d.logger.Infof("LoginByToken: %v %v", claims, resp)
			return resp, nil
		}
	}
}

func NewLoginApi(sdk conf.Sdk, logger *logrus.Entry) LoginApi {
	return defaultLoginApi{
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
