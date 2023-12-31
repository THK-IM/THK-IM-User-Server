package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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
			if errJson != nil {
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

	mSetStrings := make([]string, 0)
	for _, dtoUser := range dtoUsers {
		value, errJson := json.Marshal(dtoUser)
		if errJson == nil {
			key := fmt.Sprintf(userInfoKeyFormatter, l.appCtx.Config().Name, dtoUser.Id)
			mSetStrings = append(mSetStrings, key)
			mSetStrings = append(mSetStrings, string(value))
		}
	}

	if len(mSetStrings) > 0 {
		result, errSet := l.appCtx.RedisCache().MSet(context.Background(), mSetStrings).Result()
		if errSet != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("BatchQueryUserBasicInfoByIds mset %s", errSet.Error())
		} else {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("BatchQueryUserBasicInfoByIds mset %s", result)
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
