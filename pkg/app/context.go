package app

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-user-server/pkg/loader"
	"github.com/thk-im/thk-im-user-server/pkg/model"
)

type Context struct {
	*server.Context
	modelMap map[string]interface{}
}

func (c *Context) UserModel() model.UserModel {
	return c.modelMap["user"].(model.UserModel)
}

func (c *Context) UserOnlineRecordModel() model.UserOnlineRecordModel {
	return c.Context.ModelMap["user_online_record"].(model.UserOnlineRecordModel)
}

func (c *Context) Init(config *conf.Config) {
	c.Context = &server.Context{}
	c.Context.Init(config)
	c.modelMap = loader.LoadModels(c.Config().Models, c.Database(), c.Logger(), c.SnowflakeNode())
	err := loader.LoadTables(c.Config().Models, c.Database())
	if err != nil {
		panic(err)
	}
}
