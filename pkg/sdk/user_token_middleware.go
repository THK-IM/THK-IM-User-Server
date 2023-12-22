package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
)

const (
	TokenKey    = "Token"
	UidKey      = "Uid"
	PlatformKey = "Platform"
)

func UserTokenAuth(userApi UserApi, logger *logrus.Entry) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader(TokenKey)
		platform := context.GetHeader(PlatformKey)
		if token == "" {
			logger.Error("UserTokenAuth: token in header is empty string")
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		req := dto.TokenLoginReq{Token: token, Platform: platform}
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		userInfo, err := userApi.LoginByToken(req, claims)
		if err != nil {
			logger.Errorf("UserTokenAuth: %v %v", req, err)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		if userInfo.User == nil {
			logger.Errorf("UserTokenAuth: %v %v", req, userInfo)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		context.Set(UidKey, userInfo.User.Id)
		context.Next()
	}
}
