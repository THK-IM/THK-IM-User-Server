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
		}
	}
	return nil
}
