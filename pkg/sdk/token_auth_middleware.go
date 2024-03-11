package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
)

const (
	TokenKey    = "Token"
	UidKey      = "Uid"
	PlatformKey = "Platform"
)

func UserTokenAuth(loginApi LoginApi, logger *logrus.Entry) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		loginRsp, err := loginApi.LoginByToken(claims)
		if err != nil {
			logger.WithFields(logrus.Fields(claims)).Errorf("UserTokenAuth: %v", err)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		if loginRsp.User == nil {
			logger.WithFields(logrus.Fields(claims)).Errorf("UserTokenAuth: %v", loginRsp)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		context.Set(UidKey, loginRsp.User.Id)
		context.Next()
	}
}
