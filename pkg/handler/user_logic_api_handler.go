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

func userAccountLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.AccountLoginReq
		err := ctx.BindJSON(&req)
		if err != nil {
			appCtx.Logger().Error(err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.AccountLogin(req)
		if errReq != nil {
			appCtx.Logger().Error(errReq.Error())
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userTokenLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.TokenLoginReq
		err := ctx.BindJSON(&req)
		if err != nil {
			appCtx.Logger().Error(err.Error())
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.TokenLogin(req)
		if errReq != nil {
			appCtx.Logger().Error(errReq.Error())
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userCodeLogin(appCtx *app.Context) gin.HandlerFunc {
	// TODO: to be implemented
	// userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		baseDto.ResponseBadRequest(ctx)
	}
}

func userThirdPartLogin(appCtx *app.Context) gin.HandlerFunc {
	// TODO: to be implemented
	// userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		baseDto.ResponseBadRequest(ctx)
	}
}
