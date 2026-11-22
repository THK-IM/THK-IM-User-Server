package sdk

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
)

type LoginApi = msg.LoginApi

func NewLoginApi(config conf.Sdk, logger *logrus.Entry) LoginApi {
	return msg.NewLoginApi(config, logger)
}
