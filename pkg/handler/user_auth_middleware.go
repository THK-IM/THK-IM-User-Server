package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-user-server/pkg/app"
)

func userTokenAuth(appCtx app.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get(middleware.TokenKey)
		if token == "" {
			dto.ResponseUnauthorized(context)
			return
		}

		userId, err := utils.CheckUserToken(token, appCtx.Config().Cipher)
		if err != nil {
			dto.ResponseUnauthorized(context)
			return
		}

		context.Set(middleware.UidKey, *userId)
		context.Next()
	}
}
