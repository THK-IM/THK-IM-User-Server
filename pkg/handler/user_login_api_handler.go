package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userRegister %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}
		resp, errReq := userLoginLogic.Register(req, claims)
		if errReq != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userRegister %v %v", req, errReq)
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("userRegister %v %v", req, resp)
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
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userAccountLogin %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.AccountLogin(req, claims)
		if errReq != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userAccountLogin %v %v", req, errReq)
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("userAccountLogin %v %v", req, resp)
			baseDto.ResponseSuccess(ctx, resp)
		}
	}
}

func userTokenLogin(appCtx *app.Context) gin.HandlerFunc {
	userLoginLogic := logic.NewUserLoginLogic(appCtx)
	return func(ctx *gin.Context) {
		claims := ctx.MustGet(baseMiddleware.ClaimsKey).(baseDto.ThkClaims)
		var req dto.TokenLoginReq
		err := ctx.BindJSON(&req)
		if err != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userTokenLogin %v", err)
			baseDto.ResponseBadRequest(ctx)
			return
		}

		resp, errReq := userLoginLogic.TokenLogin(req, claims)
		if errReq != nil {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("userTokenLogin %v %v", req, errReq)
			baseDto.ResponseInternalServerError(ctx, errReq)
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("userTokenLogin %v %v", req, resp)
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
