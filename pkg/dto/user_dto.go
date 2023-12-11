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
	*User `json:"user"`
	Token string `json:"token"`
}

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginRes struct {
	*User `json:"user"`
	Token string `json:"token"`
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
