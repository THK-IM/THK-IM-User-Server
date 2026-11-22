package sdk

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
)

const PutUserOnlineStatusPath = msg.PutUserOnlineStatusPath

type UserApi = msg.UserApi

func NewUserApi(config conf.Sdk, logger *logrus.Entry) UserApi { return msg.NewUserApi(config, logger) }
