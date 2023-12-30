package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
	"strconv"
)

func queryUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserQueryLogic(appCtx)
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		id, errParam := strconv.Atoi(context.Param("id"))
		if errParam != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUser %v", errParam)
			baseDto.ResponseBadRequest(context)
			return
		}

		resp, err := l.QueryUserBasicInfoById(int64(id), claims)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("queryUser %v, %v", id, err)
			baseDto.ResponseInternalServerError(context, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("queryUser %v, %v", id, resp)
			baseDto.ResponseSuccess(context, resp)
		}
	}
}

func searchUsers(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserQueryLogic(appCtx)
	return func(context *gin.Context) {
		claims := context.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		displayId := context.Query("display_id")
		if displayId == "" {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("searchUsers %v", "display id is empty string")
			baseDto.ResponseBadRequest(context)
			return
		}
		resp, err := l.QueryUsers(displayId, claims)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("searchUsers %v %v", displayId, err)
			baseDto.ResponseInternalServerError(context, err)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("searchUsers %v %v", displayId, resp)
			baseDto.ResponseSuccess(context, resp)
		}
	}
}
