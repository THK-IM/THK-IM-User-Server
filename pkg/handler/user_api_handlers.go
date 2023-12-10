package handler

import "github.com/thk-im/thk-im-user-server/pkg/app"

func RegisterUserApiHandlers(appCtx *app.Context) {
	httpEngine := appCtx.HttpEngine()

	userGroup := httpEngine.Group("/user")
	userGroup.POST("/register", userRegister(appCtx))
	userGroup.POST("/login", userLogin(appCtx))

	userGroup.GET("/:id", queryUser(appCtx))
	userGroup.GET("", queryUsers(appCtx))

}
