package handler

import (
	baseMiddleware "github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
)

func RegisterUserApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()

	tokenAuth := userTokenAuth(appCtx)
	ipAuth := baseMiddleware.WhiteIpAuth(appCtx.Config().IpWhiteList, appCtx.Logger())

	userGroup := httpEngine.Group("/user")
	userGroup.POST("/register", userRegister(appCtx))               // 注册
	userGroup.POST("/login/account", userAccountLogin(appCtx))      // Todo: 通过账号登录
	userGroup.POST("/login/code", userCodeLogin(appCtx))            // Todo: 通过短信/邮件等验证码登录
	userGroup.POST("/login/token", userTokenLogin(appCtx))          // 通过token登录
	userGroup.POST("/login/third_part", userThirdPartLogin(appCtx)) // Todo: 第三方登录

	queryGroup := userGroup.Group("/query")
	queryGroup.Use(tokenAuth)
	queryGroup.GET("/:id", queryUser(appCtx))         // 根据id查询用户基础信息
	queryGroup.GET("/batch", batchQueryUsers(appCtx)) // 根据[id]查询用户基础信息
	queryGroup.GET("", searchUsers(appCtx))           // 根据displayId搜索用户信息

	systemGroup := userGroup.Group("/system")
	systemGroup.Use(ipAuth)
	systemGroup.POST("/online", addUserOnlineRecord(appCtx)) // 更新用户上线记录

}
