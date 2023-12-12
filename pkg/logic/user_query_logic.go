package logic

import (
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/errorx"
)

type UserQueryLogic struct {
	appCtx *app.Context
}

func NewUserQueryLogic(appCtx *app.Context) *UserQueryLogic {
	return &UserQueryLogic{appCtx: appCtx}
}

func (l *UserQueryLogic) QueryUserBasicInfoById(id int64) (*dto.BasicUser, error) {
	userInfo, err := getUserInfo(id, l.appCtx)
	if err != nil {
		return nil, err
	}
	return userDto2UserBasicDto(userInfo), err
}

func (l *UserQueryLogic) QueryUserByDisplayId(displayId string) (*dto.BasicUser, error) {
	id, err := l.appCtx.UserModel().FindUIdByDisplayId(displayId)
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, errorx.UserNotExisted
	}
	return l.QueryUserBasicInfoById(*id)
}

func userDto2UserBasicDto(user *dto.User) *dto.BasicUser {
	return &dto.BasicUser{
		Id:        user.Id,
		DisplayId: user.DisplayId,
		Avatar:    user.Avatar,
		Nickname:  user.Nickname,
		Sex:       user.Sex,
	}
}