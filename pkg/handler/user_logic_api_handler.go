package handler

import (
	"github.com/gin-gonic/gin"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
)

func userRegister(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.RegisterReq
		err := ctx.BindJSON(&req)
		if err != nil {
			appCtx.Logger().Error(err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.Register(req)
		if errReq != nil {
			appCtx.Logger().Error(errReq.Error())
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.LoginReq
		err := ctx.BindJSON(&req)
		if err != nil {
			appCtx.Logger().Error(err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.Login(req)
		if errReq != nil {
			appCtx.Logger().Error(errReq.Error())
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}
