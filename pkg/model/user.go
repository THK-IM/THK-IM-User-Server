package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
	"hash/crc32"
	"strconv"
	"time"
)

const ChannelDefault = "thk_im"
const ChannelWechat = "wechat"
const ChannelApple = "apple"

type (
	Account struct {
		Id         int64   `gorm:"id"`
		UserId     int64   `gorm:"user_id"`
		Account    string  `gorm:"account"`
		Password   *string `gorm:"password"`
		Channel    string  `gorm:"channel"`
		CreateTime int64   `gorm:"create_time"`
		UpdateTime int64   `gorm:"update_time"`
	}
	User struct {
		Id         int64   `gorm:"id"`
		DisplayId  string  `gorm:"display_id"`
		Nickname   *string `gorm:"nickname"`
		Phone      *string `gorm:"phone"`
		Sex        *int8   `gorm:"sex"`
		Birthday   *int64  `gorm:"birthday"`
		Avatar     *string `gorm:"avatar"`
		Qrcode     *string `gorm:"qrcode"`
		CreateTime int64   `gorm:"create_time"`
		UpdateTime int64   `gorm:"update_time"`
	}

	UserDisplay struct {
		DisplayId string `gorm:"display_id"`
		Id        int64  `gorm:"id"`
	}

	UserModel interface {
		AddUser(id int64, account, password, phone, nickname, avatar, qrcode *string, sex *int8, birthday *int64, channel string) (*User, error)
		UpdateUser(id int64, phone, nickname, avatar, qrcode *string, sex *int8, birthday *int64) error
		FindOne(id int64) (*User, error)
		FindUIdByDisplayId(displayId string) (*int64, error)
		FineUsers(ids []int64) ([]*User, error)
	}

	defaultUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserModel) AddUser(id int64, account, password, phone, nickname, avatar, qrcode *string, sex *int8, birthday *int64, channel string) (user *User, err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	displayId := strconv.FormatInt(id, 36)
	displayIdTableName := d.genUserDisplayIdTableName(displayId)

	now := time.Now().UnixMilli()
	user = &User{
		Id:         id,
		DisplayId:  displayId,
		Nickname:   nickname,
		Phone:      phone,
		Sex:        sex,
		Birthday:   birthday,
		Avatar:     avatar,
		Qrcode:     qrcode,
		CreateTime: now,
		UpdateTime: now,
	}
	userDisplay := &UserDisplay{
		DisplayId: displayId,
		Id:        id,
	}
	err = tx.Table(displayIdTableName).Create(userDisplay).Error
	if err != nil {
		return nil, err
	}
	err = tx.Table(d.genUserTableName(id)).Create(user).Error
	if err != nil {
		return nil, err
	}
	if account != nil {
		userAccount := &Account{
			UserId:     id,
			Account:    *account,
			Password:   password,
			Channel:    channel,
			CreateTime: now,
			UpdateTime: now,
		}
		err = tx.Table(d.genAccountTableName(id)).Create(userAccount).Error
	}
	return user, err
}

func (d defaultUserModel) UpdateUser(id int64, phone, nickname, avatar, qrcode *string, sex *int8, birthday *int64) error {
	if phone == nil && avatar == nil && nickname == nil && qrcode == nil && sex == nil && birthday == nil {
		return nil
	}
	updateMap := make(map[string]interface{})
	if phone != nil {
		updateMap["phone"] = phone
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

func (d defaultUserModel) FindOne(id int64) (*User, error) {
	tableName := d.genUserTableName(id)
	sql := fmt.Sprintf("select * from %s where id = ?", tableName)
	user := &User{}
	err := d.db.Table(tableName).Raw(sql, id).Scan(user).Error
	return user, err
}

func (d defaultUserModel) FindUIdByDisplayId(displayId string) (*int64, error) {
	tableName := d.genUserDisplayIdTableName(displayId)
	sql := fmt.Sprintf("select * from %s where display_id = ?", tableName)
	user := &UserDisplay{}
	err := d.db.Table(tableName).Raw(sql, displayId).Scan(user).Error
	if err != nil {
		return nil, err
	} else {
		if user.Id == 0 {
			return nil, nil
		} else {
			return &user.Id, nil
		}
	}
}

func (d defaultUserModel) FineUsers(ids []int64) ([]*User, error) {
	results := make([]*User, 0)
	if len(ids) == 0 {
		return results, nil
	}
	tableMap := make(map[string][]int64)
	for _, id := range ids {
		tableName := d.genUserTableName(id)
		if tableMap[tableName] == nil {
			tableMap[tableName] = []int64{id}
		} else {
			tableMap[tableName] = append(tableMap[tableName], id)
		}
	}
	for tableName, tableIds := range tableMap {
		users := make([]*User, 0)
		sql := fmt.Sprintf("select * from %s where id in ?", tableName)
		err := d.db.Raw(sql, tableIds).Scan(&users).Error
		if err != nil {
			return nil, err
		}
		results = append(results, users...)
	}
	return results, nil
}

func (d defaultUserModel) genUserDisplayIdTableName(displayId string) string {
	sum := int64(crc32.ChecksumIEEE([]byte(displayId)))
	return fmt.Sprintf("user_display_%d", sum%d.shards)
}

func (d defaultUserModel) genAccountTableName(id int64) string {
	return fmt.Sprintf("account_%d", id%d.shards)
}

func (d defaultUserModel) genUserTableName(id int64) string {
	return fmt.Sprintf("user_%d", id%d.shards)
}

func NewUserModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserModel {
	return defaultUserModel{
		shards:        shards,
		logger:        logger,
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}
