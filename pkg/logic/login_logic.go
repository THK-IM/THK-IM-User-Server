package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/o1egl/govatar"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	baseDto "github.com/thk-im/thk-im-base-server/dto"
	baseErr "github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/errorx"
	"github.com/thk-im/thk-im-user-server/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"image/color"
	"os"
	"strings"
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

func (l *UserLoginLogic) Register(req dto.RegisterReq, claims baseDto.ThkClaims) (*dto.RegisterRes, error) {
	id := l.appCtx.SnowflakeNode().Generate().Int64()
	if (req.Account == nil) != (req.Password == nil) {
		return nil, baseErr.ErrParamsError
	}
	if req.Account != nil {
		account := strings.ToLower(strings.TrimSpace(*req.Account))
		if account == "" || len(account) > 64 || *req.Password == "" || len(*req.Password) > 72 {
			return nil, baseErr.ErrParamsError
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashed := string(hash)
		req.Account, req.Password = &account, &hashed
	}
	if err := os.MkdirAll("tmp", 0700); err != nil {
		return nil, err
	}
	reqSex := req.Sex
	if reqSex == nil {
		reqSex = &sexMale
	}
	nickname := req.Nickname
	if nickname == nil {
		newNickName := utils.GetRandomString(8)
		nickname = &newNickName
	}
	avatarUrl := req.Avatar
	if avatarUrl == nil && l.appCtx.ObjectStorage() != nil {
		fileName := fmt.Sprintf("%d.jpg", id)
		filePath := fmt.Sprintf("tmp/%s", fileName)
		defer os.Remove(filePath)
		male := govatar.MALE
		if *reqSex == sexFemale {
			male = govatar.FEMALE
		}
		err := govatar.GenerateFileForUsername(male, *nickname, filePath)
		if err == nil {
			avatarKey := fmt.Sprintf("user/%d/avatar/%s", id, fileName)
			avatarUrl, err = l.appCtx.ObjectStorage().UploadObject(avatarKey, filePath)
			if err != nil {
				l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("Register %v %v", req, err)
			}
		} else {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("Register %v %v", req, err)
		}
	}

	var qrcodeUrl *string = nil
	qrFileName := fmt.Sprintf("%d-qrcode.png", id)
	qrFilePath := fmt.Sprintf("tmp/%s", qrFileName)
	defer os.Remove(qrFilePath)
	url := fmt.Sprintf("https://api.thkim.com/user/%d", id)
	errQrcode := qrcode.WriteColorFile(url, qrcode.Medium, 256, color.Black, color.White, qrFilePath)
	if errQrcode != nil {
		l.appCtx.Logger().Error(errQrcode)
	} else if l.appCtx.ObjectStorage() != nil {
		qrCodeKey := fmt.Sprintf("user/%d/qrcode/%s", id, qrFileName)
		qrcodeUrl, errQrcode = l.appCtx.ObjectStorage().UploadObject(qrCodeKey, qrFilePath)
		if errQrcode != nil {
			l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("Register %v %v", req, errQrcode)
		}
	}

	user, err := l.appCtx.UserModel().AddUser(id, req.Account, req.Password, nil, nickname,
		avatarUrl, qrcodeUrl, reqSex, req.Birthday, model.ChannelDefault,
	)
	if err != nil {
		return nil, err
	}

	token, errToken := l.issueToken(user.Id, 0)
	if errToken != nil {
		return nil, errToken
	}
	dtoUser := userModel2UserDto(user)
	resp := &dto.RegisterRes{
		User:  dtoUser,
		Token: token,
	}
	return resp, nil
}

func (l *UserLoginLogic) AccountLogin(req dto.AccountLoginReq, claims baseDto.ThkClaims) (*dto.LoginRes, error) {
	account, err := l.appCtx.UserModel().Authenticate(strings.ToLower(strings.TrimSpace(*req.Account)), *req.Password)
	if err != nil {
		return nil, err
	}
	user, err := getUserInfo(account.UserId, l.appCtx)
	if err != nil {
		return nil, err
	}
	token, err := l.issueToken(account.UserId, account.Version)
	return &dto.LoginRes{User: user, Token: token}, err
}

func (l *UserLoginLogic) TokenLogin(claims baseDto.ThkClaims) (*dto.LoginRes, error) {
	auth, err := l.ValidateToken(claims.GetToken())
	if err != nil {
		return nil, err
	}
	userInfo, errUser := getUserInfo(auth.UID, l.appCtx)
	if errUser != nil {
		l.appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("TokenLogin %v", err)
		return nil, errUser
	}
	loginRes := &dto.LoginRes{
		Token: claims.GetToken(),
		User:  userInfo,
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
		if user.Id == 0 {
			return nil, errorx.UserNotExisted
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
