package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func userTokenAuth(appCtx *app.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		uId, err := utils.CheckUserToken(claims.GetToken(), appCtx.Config().Cipher)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("CheckUserToken err: %v", err)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		context.Set(sdk.UidKey, uId)
		context.Next()
	}
}
