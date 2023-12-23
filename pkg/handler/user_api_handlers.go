package handler

import (
	"github.com/thk-im/thk-im-user-server/pkg/app"
)

func RegisterUserApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()

	userGroup := httpEngine.Group("/user")
	userGroup.POST("/register", userRegister(appCtx))               // 注册
	userGroup.POST("/login/account", userAccountLogin(appCtx))      // Todo: 通过账号登录
	userGroup.POST("/login/code", userCodeLogin(appCtx))            // Todo: 通过短信/邮件等验证码登录
	userGroup.POST("/login/token", userTokenLogin(appCtx))          // 通过token登录
	userGroup.POST("/login/third_part", userThirdPartLogin(appCtx)) // Todo: 第三方登录
	userGroup.GET("/:id", queryUser(appCtx))                        // 根据id查询用户基础信息
	userGroup.GET("", searchUsers(appCtx))                          // 搜索用户信息
	userGroup.POST("/online", addUserOnlineRecord(appCtx))          // 更新用户上线记录

}
