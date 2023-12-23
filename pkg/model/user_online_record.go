package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"gorm.io/gorm"
)

type (
	UserOnlineRecord struct {
		UserId      int64  `gorm:"user_id"`
		Platform    string `gorm:"platform"`
		ConnId      int64  `gorm:"conn_id"`
		OnlineTime  int64  `gorm:"online_time"`
		OfflineTime int64  `gorm:"offline_time"`
	}

	UserOnlineRecordModel interface {
		GetUserLastOnlineRecord(userId int64) (*UserOnlineRecord, error)
		GetUsersOnlineRecords(userId int64) ([]*UserOnlineRecord, error)
		CreateUserOnlineRecord(userId, connId int64, platform string, onlineTime int64) error
		UpdateUserOnlineRecord(userId, connId int64, platform string, offlineTime int64) error
	}

	defaultUserOnlineRecordModel struct {
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
		shards        int64
	}
)

func (d defaultUserOnlineRecordModel) GetUserLastOnlineRecord(userId int64) (*UserOnlineRecord, error) {
	tableName := d.genUserOnlineRecordTable(userId)
	usersOnlineRecord := &UserOnlineRecord{}
	sql := fmt.Sprintf("select * from %s where user_id = ? order by online_time desc limit 0, 1", tableName)
	err := d.db.Raw(sql, userId).Scan(&usersOnlineRecord).Error
	return usersOnlineRecord, err
}

func (d defaultUserOnlineRecordModel) GetUsersOnlineRecords(userId int64) ([]*UserOnlineRecord, error) {
	usersOnlineRecords := make([]*UserOnlineRecord, 0)
	tableName := d.genUserOnlineRecordTable(userId)
	sql := fmt.Sprintf("select * from %s where user_id = ?", tableName)
	err := d.db.Raw(sql, userId).Scan(&usersOnlineRecords).Error
	return usersOnlineRecords, err
}

func (d defaultUserOnlineRecordModel) CreateUserOnlineRecord(userId, connId int64, platform string, onlineTime int64) error {
	tableName := d.genUserOnlineRecordTable(userId)
	usersOnlineRecord := &UserOnlineRecord{
		UserId:      userId,
		Platform:    platform,
		ConnId:      connId,
		OnlineTime:  onlineTime,
		OfflineTime: 0,
	}
	return d.db.Table(tableName).Create(usersOnlineRecord).Error
}

func (d defaultUserOnlineRecordModel) UpdateUserOnlineRecord(userId, connId int64, platform string, offlineTime int64) error {
	tableName := d.genUserOnlineRecordTable(userId)
	sql := fmt.Sprintf("update %s set offline_time = ? where user_id = ? and conn_id = ?", tableName)
	return d.db.Exec(sql, offlineTime, userId, connId).Error
}

func (d defaultUserOnlineRecordModel) genUserOnlineRecordTable(userId int64) string {
	return fmt.Sprintf("user_online_record_%d", userId%(d.shards))
}

func NewUserOnlineRecordModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserOnlineRecordModel {
	return defaultUserOnlineRecordModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
