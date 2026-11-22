package handler

import (
	"github.com/gin-gonic/gin"
	base "github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/errorx"
	msg "github.com/thk-im/thk-im-msgapi-server/pkg/dto"
	"github.com/thk-im/thk-im-msgapi-server/pkg/model"
	"github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
)

func checkMessage(contact sdk.MsgCheckerApi) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req msg.CheckMessageReq
		if err := c.ShouldBindJSON(&req); err != nil {
			base.ResponseBadRequest(c)
			return
		}
		if req.SessionType != model.SingleSessionType || req.FromUId <= 0 || req.MessageType <= 0 {
			base.ResponseSuccess(c, &msg.CheckMessageRes{})
			return
		}
		if req.EntityId <= 0 {
			base.ResponseBadRequest(c)
			return
		}
		if contact == nil {
			base.ResponseInternalServerError(c, errorx.ErrInternalServerError)
			return
		}
		claims := c.MustGet(base.ClaimsKey).(base.ThkClaims)
		result, err := contact.CheckMessage(&req, claims)
		if err != nil {
			base.ResponseInternalServerError(c, err)
			return
		}
		base.ResponseSuccess(c, result)
	}
}
