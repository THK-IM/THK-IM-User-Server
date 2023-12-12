package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
	"strconv"
)

func queryUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserQueryLogic(appCtx)
	return func(context *gin.Context) {
		id, errParam := strconv.Atoi(context.Param("id"))
		if errParam != nil {
			appCtx.Logger().Error(errParam.Error())
			baseDto.ResponseBadRequest(context)
			return
		}

		resp, err := l.QueryUserBasicInfoById(int64(id))
		if err != nil {
			baseDto.ResponseInternalServerError(context, err)
		} else {
			baseDto.ResponseSuccess(context, resp)
		}
	}
}

func searchUsers(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserQueryLogic(appCtx)
	return func(context *gin.Context) {
		displayId := context.Query("id")
		if displayId == "" {
			baseDto.ResponseBadRequest(context)
			return
		}
		resp, err := l.QueryUserByDisplayId(displayId)
		if err != nil {
			baseDto.ResponseInternalServerError(context, err)
		} else {
			baseDto.ResponseSuccess(context, resp)
		}
	}
}
