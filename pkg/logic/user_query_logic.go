package logic

import (
	"context"
	"encoding/json"
	"fmt"

	baseDto "github.com/thk-im/thk-im-base-server/dto"
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

func (l *UserQueryLogic) QueryUserBasicInfoById(id int64, claims baseDto.ThkClaims) (*dto.BasicUser, error) {
	userInfo, err := getUserInfo(id, l.appCtx)
	if err != nil {
		return nil, err
	}
	return userDto2UserBasicDto(userInfo), err
}

func (l *UserQueryLogic) BatchQueryUserBasicInfoByIds(ids []int64, claims baseDto.ThkClaims) (map[int64]*dto.BasicUser, error) {
	if len(ids) == 0 {
		return map[int64]*dto.BasicUser{}, nil
	}
	userInfoKeys := make([]string, 0)
	for _, id := range ids {
		userInfoKey := fmt.Sprintf(userInfoKeyFormatter, l.appCtx.Config().Name, id)
		userInfoKeys = append(userInfoKeys, userInfoKey)
	}
	interfaceUserInfos, err := l.appCtx.RedisCache().MGet(context.Background(), userInfoKeys...).Result()
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*dto.BasicUser)
	for _, interfaceUserInfo := range interfaceUserInfos {
		jsonUserInfo, ok := interfaceUserInfo.(string)
		if ok {
			user := &dto.User{}
			errJson := json.Unmarshal([]byte(jsonUserInfo), user)
			if errJson == nil && user.Id > 0 {
				baseUser := &dto.BasicUser{
					Id:        user.Id,
					DisplayId: user.DisplayId,
					Avatar:    user.Avatar,
					Nickname:  user.Nickname,
					Sex:       user.Sex,
				}
				userMap[user.Id] = baseUser
			}
		}
	}

	idsInDB := make([]int64, 0)
	for _, id := range ids {
		if userMap[id] == nil {
			idsInDB = append(idsInDB, id)
		}
	}

	dtoUsers := make([]*dto.User, 0)
	if len(idsInDB) > 0 {
		users, errDb := l.appCtx.UserModel().FineUsers(idsInDB)
		if errDb != nil {
			return nil, errDb
		}
		for _, user := range users {
			dtoUser := userModel2UserDto(user)
			dtoUsers = append(dtoUsers, dtoUser)
			baseUser := userDto2UserBasicDto(dtoUser)
			userMap[user.Id] = baseUser
		}
	}

	for _, user := range dtoUsers {
		if err := setUserInfoCache(user, l.appCtx); err != nil {
			l.appCtx.Logger().Warn(err)
		}
	}
	return userMap, nil
}

func (l *UserQueryLogic) QueryUserByDisplayId(displayId string, claims baseDto.ThkClaims) (*dto.BasicUser, error) {
	id, err := l.appCtx.UserModel().FindUIdByDisplayId(displayId)
	if err != nil {
		return nil, err
	}
	if id == nil {
		return nil, errorx.UserNotExisted
	}
	return l.QueryUserBasicInfoById(*id, nil)
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

func (l *UserQueryLogic) OwnProfile(uid int64) (*dto.User, error) {
	return getUserInfo(uid, l.appCtx)
}
