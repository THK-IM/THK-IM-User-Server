package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
)

const (
	TokenKey    = msg.TokenKey
	UidKey      = msg.UidKey
	PlatformKey = msg.PlatformKey
)

func UserTokenAuth(api LoginApi, logger *logrus.Entry) gin.HandlerFunc {
	return msg.UserTokenAuth(api, logger)
}
