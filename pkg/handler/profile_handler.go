package handler

import (
	"github.com/gin-gonic/gin"
	base "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/logic"
	"github.com/thk-im/thk-im-user-server/pkg/sdk"
)

func logout(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet(middleware.ClaimsKey).(base.ThkClaims)
		if err := logic.NewUserLoginLogic(appCtx).Logout(claims.GetToken()); err != nil {
			base.ResponseInternalServerError(c, err)
			return
		}
		base.ResponseSuccess(c, nil)
	}
}
func changePassword(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Old string `json:"old_password" binding:"required"`
			New string `json:"new_password" binding:"required"`
		}
		if c.ShouldBindJSON(&req) != nil || len(req.New) > 72 {
			base.ResponseBadRequest(c)
			return
		}
		resp, err := logic.NewUserLoginLogic(appCtx).ChangePassword(c.MustGet("login_auth").(logic.Auth), req.Old, req.New)
		if err != nil {
			base.ResponseInternalServerError(c, err)
			return
		}
		base.ResponseSuccess(c, resp)
	}
}
func updateProfile(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UpdateProfileReq
		if c.ShouldBindJSON(&req) != nil {
			base.ResponseBadRequest(c)
			return
		}
		resp, err := logic.NewUserLoginLogic(appCtx).UpdateProfile(c.GetInt64(sdk.UidKey), req)
		if err != nil {
			base.ResponseInternalServerError(c, err)
			return
		}
		base.ResponseSuccess(c, resp)
	}
}

func ownProfile(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := logic.NewUserQueryLogic(appCtx).OwnProfile(c.GetInt64(sdk.UidKey))
		if err != nil {
			base.ResponseInternalServerError(c, err)
			return
		}
		base.ResponseSuccess(c, user)
	}
}
