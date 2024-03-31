package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-user-server/pkg/app"
)

func postUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		baseDto.ResponseBadRequest(ctx)
	}
}
