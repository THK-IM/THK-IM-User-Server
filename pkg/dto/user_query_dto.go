package dto

type BasicUser struct {
	Id        int64   `json:"id"`
	DisplayId string  `json:"display_id"`
	Avatar    *string `json:"avatar,omitempty"`
	Nickname  *string `json:"nickname,omitempty"`
	Sex       *int8   `json:"sex,omitempty"`
}
