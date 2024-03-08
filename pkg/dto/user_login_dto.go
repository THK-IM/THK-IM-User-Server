package dto

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

type LoginRes struct {
	Id   int64 `json:"id"`
	User *User `json:"user,omitempty"`
}

type User struct {
	Id        int64   `json:"id"`
	DisplayId string  `json:"display_id"`
	Avatar    *string `json:"avatar,omitempty"`
	Nickname  *string `json:"nickname,omitempty"`
	Qrcode    *string `json:"qrcode,omitempty"`
	Sex       *int8   `json:"sex,omitempty"`
	Birthday  *int64  `json:"birthday,omitempty"`
}

type UserOnlineStatusReq struct {
	UserId      int64  `json:"user_id"`
	IsOnline    bool   `json:"is_online"`
	TimestampMs int64  `json:"timestamp_ms"`
	ConnId      int64  `json:"conn_id"`
	Platform    string `json:"platform"` // Android/Ios/Web/Linux/Mac/Windows
}
