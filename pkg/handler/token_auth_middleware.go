package handler

import (
	"github.com/gin-gonic/gin"

	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
	"github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func userTokenAuth(appCtx *app.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		auth, err := logic.NewUserLoginLogic(appCtx).ValidateToken(claims.GetToken())
		if err != nil {
			appCtx.Logger().Debugf("CheckUserToken: %v", err)
			baseDto.ResponseForbidden(context)
			context.Abort()
			return
		}
		context.Set(sdk.UidKey, auth.UID)
		context.Set("login_auth", auth)
		context.Next()
	}
}
