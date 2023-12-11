package logic

import (
	"fmt"
	"github.com/o1egl/govatar"
	"github.com/skip2/go-qrcode"
	"github.com/thk-im/thk-im-base-server/utils"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/dto"
	"github.com/thk-im/thk-im-user-server/pkg/model"
	"image/color"
	"time"
)

var (
	sexMale   int8 = 0
	sexFemale int8 = 1
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

	qrFileName := fmt.Sprintf("%s-%d-qrcode.png", *nickname, time.Now().UnixMilli()/1000)
	qrFilePath := fmt.Sprintf("tmp/%s", qrFileName)
	url := fmt.Sprintf("https://api.thkim.com/user/%d", id)
	errQrcode := qrcode.WriteColorFile(url, qrcode.Medium, 256, color.Black, color.White, qrFilePath)
	if errQrcode != nil {
		l.appCtx.Logger().Error(errQrcode)
	}
	qrCodeKey := fmt.Sprintf("avatar/%d/%s", id, qrFileName)
	qrAvatarUrl, errUpload := l.appCtx.ObjectStorage().UploadObject(qrCodeKey, qrFilePath)
	if errUpload != nil {
		l.appCtx.Logger().Error("upload object file error: ", errUpload)
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
		User:  l.userModel2Dto(user),
		Token: token,
	}
	return resp, nil
}

func (l *UserLoginLogic) Login(req dto.LoginReq) (*dto.LoginRes, error) {
	return nil, nil
}

func (l *UserLoginLogic) userModel2Dto(user *model.User) *dto.User {
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
