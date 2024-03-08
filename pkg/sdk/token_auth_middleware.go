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
		userInfoResp, err := loginApi.LoginByToken(claims)
		if err != nil {
			logger.WithFields(logrus.Fields(claims)).Errorf("UserTokenAuth: %v", err)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		if userInfoResp.Id <= 0 {
			logger.WithFields(logrus.Fields(claims)).Errorf("UserTokenAuth: %v", userInfoResp)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		context.Set(UidKey, userInfoResp.User.Id)
		context.Next()
	}
}
