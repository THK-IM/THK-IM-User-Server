package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
)

func postUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserOnlineStatusReq
		if ctx.ShouldBindJSON(&req) != nil || req.UserId <= 0 || req.ConnId <= 0 || req.TimestampMs <= 0 || req.Platform == "" {
			baseDto.ResponseBadRequest(ctx)
			return
		}
		m := appCtx.UserOnlineRecordModel()
		var err error
		if req.IsOnline {
			err = m.CreateUserOnlineRecord(req.UserId, req.ConnId, req.Platform, req.TimestampMs)
		} else {
			err = m.UpdateUserOnlineRecord(req.UserId, req.ConnId, req.Platform, req.TimestampMs)
		}
		if err != nil {
			baseDto.ResponseInternalServerError(ctx, err)
			return
		}
		baseDto.ResponseSuccess(ctx, nil)
	}
}
