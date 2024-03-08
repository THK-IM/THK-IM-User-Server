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
	queryUserPath           string = "/user"
	PutUserOnlineStatusPath string = "/user/:id/online_status"
)

type (
	UserApi interface {
		PostUserOnlineStatus(req *dto.UserOnlineStatusReq, claims baseDto.ThkClaims) error
		QueryUsers(req *dto.QueryUsers, claims baseDto.ThkClaims) (map[int64]*dto.BasicUser, error)
	}

	defaultUserApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultUserApi) PostUserOnlineStatus(req *dto.UserOnlineStatusReq, claims baseDto.ThkClaims) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, PutUserOnlineStatusPath)
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", contentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := baseErrorx.NewErrorXFromResp(res)
		d.logger.Errorf("PostUserOnlineStatus: %v %v", req, e)
		return e
	} else {
		d.logger.Infof("PostUserOnlineStatus: %v %s", req, "success")
		return nil
	}
}

func (d defaultUserApi) QueryUsers(req *dto.QueryUsers, claims baseDto.ThkClaims) (map[int64]*dto.BasicUser, error) {
	url := fmt.Sprintf("%s%s", d.endpoint, queryUserPath)
	for i, id := range req.Ids {
		if i == 0 {
			url += fmt.Sprintf("?ids=%d", id)
		} else {
			url += fmt.Sprintf("&ids=%d", id)
		}
	}
	request := d.client.R()
	for k, v := range claims {
		vs := v.(string)
		request.SetHeader(k, vs)
	}
	res, errRequest := request.
		SetHeader("Content-Type", contentType).
		Get(url)
	if errRequest != nil {
		d.logger.Errorf("QueryUsers %v %v", claims, errRequest)
		return nil, errRequest
	}
	if res.StatusCode() != http.StatusOK {
		errRes := baseErrorx.NewErrorXFromResp(res)
		d.logger.Errorf("QueryUsers: %v %v", claims, errRes)
		return nil, errRes
	} else {
		resp := make(map[int64]*dto.BasicUser)
		e := json.Unmarshal(res.Body(), &resp)
		if e != nil {
			d.logger.Errorf("QueryUsers: %v %v", claims, e)
			return nil, e
		} else {
			d.logger.Infof("QueryUsers: %v %v", claims, resp)
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
