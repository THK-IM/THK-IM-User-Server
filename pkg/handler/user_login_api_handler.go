package handler

import (
	"github.com/gin-gonic/gin"

	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
)

func userRegister(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.RegisterReq
		err := ctx.BindJSON(&req)
		if err != nil {
			baseDto.ResponseBadRequest(ctx)
			return
		}
		resp, errReq := userLoginLogic.Register(req, claims)
		if errReq != nil {
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userAccountLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.AccountLoginReq
		err := ctx.BindJSON(&req)
		if err != nil {
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.AccountLogin(req, claims)
		if errReq != nil {
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userTokenLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		resp, errReq := userLoginLogic.TokenLogin(claims)
		if errReq != nil {
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
