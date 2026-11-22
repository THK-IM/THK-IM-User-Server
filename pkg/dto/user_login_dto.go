package dto

import msg "github.com/thk-im/thk-im-msgapi-server/pkg/dto"

type RegisterReq struct {
	Account  *string `json:"account"`
	Password *string `json:"password"`
	Nickname *string `json:"nickname"`
	Sex      *int8   `json:"sex"`
	Avatar   *string `json:"avatar"`
	Birthday *int64  `json:"birthday"`
}

type RegisterRes struct {
	User  *User  `json:"user,omitempty"`
	Token string `json:"token"`
}

type AccountLoginReq struct {
	Account  *string `json:"account" binding:"required"`
	Password *string `json:"password" binding:"required"`
	Platform *string `json:"platform" binding:"required"`
}

type ThirdPartLoginReq struct {
	Channel    *string `json:"channel" binding:"required"`
	ThirdToken *string `json:"third_token" binding:"required"`
	Platform   *string `json:"platform" binding:"required"`
}

type LoginRes = msg.LoginRes
type User = msg.User
type UserOnlineStatusReq = msg.UserOnlineStatusReq

type UpdateProfileReq struct {
	Nickname *string `json:"nickname" binding:"omitempty,min=1,max=64"`
	Avatar   *string `json:"avatar"`
	Sex      *int8   `json:"sex" binding:"omitempty,oneof=0 1"`
	Birthday *int64  `json:"birthday"`
}
