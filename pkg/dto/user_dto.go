package dto

type RegisterReq struct {
	Account  *string `json:"account"`
	Password *string `json:"password"`
	Nickname *string `json:"nickname"`
	Sex      *int8   `json:"sex"`
	Avatar   *string `json:"avatar"`
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
	Avatar    *string `json:"avatar"`
	Nickname  *string `json:"nickname"`
	Qrcode    *string `json:"qrcode"`
	Sex       *int8   `json:"sex"`
	Birthday  *int64  `json:"birthday"`
}
