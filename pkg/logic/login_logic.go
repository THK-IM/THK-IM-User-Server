package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/o1egl/govatar"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/errorx"
	"github.com/thk-im/thk-im-user-server/pkg/model"
	"image/color"
	"time"
)

var (
	sexMale   int8 = 0
	sexFemale int8 = 1

	userInfoKeyFormatter = "%s:u_info:%d"
)

type UserLoginLogic struct {
	appCtx *app.Context
}

func NewUserLoginLogic(appCtx *app.Context) *UserLoginLogic {
	return &UserLoginLogic{appCtx: appCtx}
}

func (l *UserLoginLogic) Register(req dto.RegisterReq) (*dto.RegisterRes, error) {
	id := l.appCtx.SnowflakeNode().Generate().Int64()
	sex := req.Sex
	if sex == nil {
		sex = &sexMale
	}
	nickname := req.Nickname
	if nickname == nil {
		newNickName := utils.GetRandomString(8)
		nickname = &newNickName
	}
	avatarUrl := req.Avatar
	if avatarUrl == nil {
		fileName := fmt.Sprintf("%s-%d.jpg", *nickname, time.Now().UnixMilli()/1000)
		filePath := fmt.Sprintf("tmp/%s", fileName)
		male := govatar.MALE
		if *sex == sexFemale {
			male = govatar.FEMALE
		}
		err := govatar.GenerateFileForUsername(male, *nickname, filePath)
		if err == nil {
			avatarKey := fmt.Sprintf("avatar/%d/%s", id, fileName)
			avatarUrl, err = l.appCtx.ObjectStorage().UploadObject(avatarKey, filePath)
			if err != nil {
				l.appCtx.Logger().Error("upload object file error: ", err)
			}
		} else {
			l.appCtx.Logger().Error("go avatar generate file error: ", err)
		}
	}

	var qrAvatarUrl *string = nil
	qrFileName := fmt.Sprintf("%s-%d-qrcode.png", *nickname, time.Now().UnixMilli()/1000)
	qrFilePath := fmt.Sprintf("tmp/%s", qrFileName)
	url := fmt.Sprintf("https://api.thkim.com/user/%d", id)
	errQrcode := qrcode.WriteColorFile(url, qrcode.Medium, 256, color.Black, color.White, qrFilePath)
	if errQrcode != nil {
		l.appCtx.Logger().Error(errQrcode)
	} else {
		qrCodeKey := fmt.Sprintf("avatar/%d/%s", id, qrFileName)
		qrAvatarUrl, errQrcode = l.appCtx.ObjectStorage().UploadObject(qrCodeKey, qrFilePath)
		if errQrcode != nil {
			l.appCtx.Logger().Error("upload object file error: ", errQrcode)
		}

	}

	user, err := l.appCtx.UserModel().AddUser(id, req.Account, req.Password, nil, req.Nickname,
		avatarUrl, qrAvatarUrl, sex, nil, model.ChannelDefault,
	)
	if err != nil {
		return nil, err
	}

	token, errToken := utils.GenerateUserToken(user.Id, l.appCtx.Config().Name, l.appCtx.Config().Cipher)
	if err != nil {
		return nil, errToken
	}
	resp := &dto.RegisterRes{
		User:  userModel2UserDto(user),
		Token: token,
	}
	return resp, nil
}

func (l *UserLoginLogic) AccountLogin(req dto.AccountLoginReq) (*dto.LoginRes, error) {
	return nil, nil
}

func (l *UserLoginLogic) TokenLogin(req dto.TokenLoginReq) (*dto.LoginRes, error) {
	uId, err := utils.CheckUserToken(*req.Token, l.appCtx.Config().Cipher)
	if err != nil {
		l.appCtx.Logger().Error(err, req)
		return nil, errorx.UserTokenError
	}
	userInfo, errUser := getUserInfo(*uId, l.appCtx)
	if errUser != nil {
		l.appCtx.Logger().Error(err, req)
		return nil, errUser
	}
	loginRes := &dto.LoginRes{
		User: userInfo,
	}
	return loginRes, nil
}

func getUserInfo(uId int64, appCtx *app.Context) (*dto.User, error) {
	userInfoKey := fmt.Sprintf(userInfoKeyFormatter, appCtx.Config().Name, uId)
	userInfoJson, errCache := appCtx.RedisCache().Get(context.Background(), userInfoKey).Result()
	if errCache != nil && !errors.Is(errCache, redis.Nil) {
		appCtx.Logger().Error(errCache)
		return nil, errCache
	}
	if errors.Is(errCache, redis.Nil) {
		user, errDb := appCtx.UserModel().FindOne(uId)
		if errDb != nil {
			appCtx.Logger().Error(errDb)
			return nil, errDb
		}
		userInfo := userModel2UserDto(user)
		errCache = setUserInfoCache(userInfo, appCtx)
		if errCache != nil {
			appCtx.Logger().Error(errCache)
		}
		return userInfo, nil
	}
	userInfo := &dto.User{}
	errJson := json.Unmarshal([]byte(userInfoJson), userInfo)
	return userInfo, errJson
}

func setUserInfoCache(user *dto.User, appCtx *app.Context) error {
	userInfoKey := fmt.Sprintf(userInfoKeyFormatter, appCtx.Config().Name, user.Id)
	userInfoJson, errJson := json.Marshal(user)
	if errJson != nil {
		appCtx.Logger().Error(errJson)
		return errJson
	}
	return appCtx.RedisCache().Set(context.Background(), userInfoKey, string(userInfoJson), time.Hour*7*24).Err()
}

func userModel2UserDto(user *model.User) *dto.User {
	return &dto.User{
		Id:        user.Id,
		DisplayId: user.DisplayId,
		Avatar:    user.Avatar,
		Nickname:  user.Nickname,
		Qrcode:    user.Qrcode,
		Sex:       user.Sex,
		Birthday:  user.Birthday,
	}
}
