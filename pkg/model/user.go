package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type (
	User struct {
		Id         int64  `gorm:"id"`
		DisplayId  int64  `gorm:"display_id"`
		Nickname   string `gorm:"nickname"`
		Sex        int8   `gorm:"sex"`
		Birthday   int64  `gorm:"birthday"`
		Avatar     string `gorm:"avatar"`
		Qrcode     string `gorm:"qrcode"`
		CreateTime int64  `gorm:"create_time"`
		UpdateTime int64  `gorm:"update_time"`
	}

	UserModel interface {
		AddUser(displayId int64, nickname, avatar, qrcode string, sex int8, birthday int64) (*User, error)
		UpdateUser(id int64, displayId *int64, nickname, avatar, qrcode *string, sex *int8, birthday *int64) error
	}

	defaultUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserModel) AddUser(displayId int64, nickname, avatar, qrcode string, sex int8, birthday int64) (*User, error) {
	now := time.Now().UnixMilli()
	user := &User{
		DisplayId:  displayId,
		Nickname:   nickname,
		Sex:        sex,
		Birthday:   birthday,
		Avatar:     avatar,
		Qrcode:     qrcode,
		CreateTime: now,
		UpdateTime: now,
	}
	err := d.db.Table(d.genUserTableName(0)).Create(user).Error
	return user, err
}

func (d defaultUserModel) UpdateUser(id int64, displayId *int64, nickname, avatar, qrcode *string, sex *int8, birthday *int64) error {
	if displayId == nil && avatar == nil && nickname == nil && qrcode == nil && sex == nil && birthday == nil {
		return nil
	}
	updateMap := make(map[string]interface{})
	if displayId != nil {
		updateMap["display_id"] = *displayId
	}
	if avatar != nil {
		updateMap["avatar"] = *avatar
	}
	if nickname != nil {
		updateMap["nickname"] = *nickname
	}
	if qrcode != nil {
		updateMap["qrcode"] = *qrcode
	}
	if sex != nil {
		updateMap["sex"] = *sex
	}
	if birthday != nil {
		updateMap["birthday"] = *birthday
	}
	updateMap["update_time"] = time.Now().UnixMilli()
	return d.db.Table(d.genUserTableName(id)).Where("id = ?", id).Updates(updateMap).Error
}

func (d defaultUserModel) genUserTableName(id int64) string {
	return fmt.Sprintf("user_%d", id/d.shards)
}

func NewUserModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserModel {
	return defaultUserModel{
		shards:        shards,
		logger:        logger,
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}
