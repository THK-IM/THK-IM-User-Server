package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	msgSdk "github.com/thk-im/thk-im-msgapi-server/pkg/sdk"
	"github.com/thk-im/thk-im-user-server/pkg/loader"
	"github.com/thk-im/thk-im-user-server/pkg/model"
)

type Context struct {
	*server.Context
	ContactChecker msgSdk.MsgCheckerApi
}

func (c *Context) UserModel() model.UserModel {
	return c.ModelMap["user"].(model.UserModel)
}

func (c *Context) UserOnlineRecordModel() model.UserOnlineRecordModel {
	return c.ModelMap["user_online_record"].(model.UserOnlineRecordModel)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	for _, api := range config.Sdks {
		if api.Name == "contact_check_api" && api.Endpoint != "" {
			c.ContactChecker = msgSdk.NewMsgCheckerApi(api, c.Logger())
			break
		}
	}
	c.ModelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
