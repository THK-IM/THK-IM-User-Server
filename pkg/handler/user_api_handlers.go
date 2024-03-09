package handler

import (
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
)

func RegisterUserApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()

	tokenAuth := userTokenAuth(appCtx)
	ipAuth := baseMiddleware.WhiteIpAuth(appCtx.Config().IpWhiteList, appCtx.Logger())

	loginGroup := httpEngine.Group("/login")
	loginGroup.POST("/register", userRegister(appCtx))         // 注册
	loginGroup.POST("/account", userAccountLogin(appCtx))      // Todo: 通过账号登录
	loginGroup.POST("/code", userCodeLogin(appCtx))            // Todo: 通过短信/邮件等验证码登录
	loginGroup.POST("/token", userTokenLogin(appCtx))          // 通过token登录
	loginGroup.POST("/third_part", userThirdPartLogin(appCtx)) // Todo: 第三方登录

	userGroup := httpEngine.Group("/user")
	userGroup.Use(tokenAuth)
	userGroup.GET("/:id", queryUser(appCtx)) // 根据id查询用户基础信息
	userGroup.POST("/:id/online_status", postUserOnlineStatus(appCtx))
	userGroup.GET("", queryUsers(appCtx))
	userGroup.GET("/search", searchUsers(appCtx))

	systemGroup := userGroup.Group("/system")
	systemGroup.Use(ipAuth)

}
