package loader

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"github.com/thk-im/thk-im-user-server/pkg/model"
	"gorm.io/gorm"
	"os"
)

func LoadModels(modeConfigs []conf.Model, database *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node) map[string]interface{} {
	modelMap := make(map[string]interface{}, 0)
	for _, ms := range modeConfigs {
		var m interface{}
		if ms.Name == "user" {
			m = model.NewUserModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "user_online_record" {
			m = model.NewUserOnlineRecordModel(database, logger, snowflakeNode, ms.Shards)
		}
		modelMap[ms.Name] = m
	}
	return modelMap
}

func LoadTables(modeConfigs []conf.Model, database *gorm.DB) error {
	for _, ms := range modeConfigs {
		path := fmt.Sprintf("./sql/%s.sql", ms.Name)
		buffer, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for i := int64(0); i < ms.Shards; i++ {
			sql := fmt.Sprintf(string(buffer), fmt.Sprintf("%d", i))
			err = database.Exec(sql).Error
			if err != nil {
				return err
			}
			if ms.Name == "account" {
				// An ambiguous legacy account must be resolved before it can authenticate.
				q := fmt.Sprintf("INSERT INTO account_login_0 (account,user_id) SELECT LOWER(TRIM(a.account)),a.user_id FROM account_%d a LEFT JOIN account_login_0 l ON l.account=LOWER(TRIM(a.account)) AND l.user_id=a.user_id WHERE a.channel=? AND l.user_id IS NULL", i)
				if err := database.Exec(q, model.ChannelDefault).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
