package main

import (
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-user-server/pkg/app"
	"github.com/thk-im/thk-im-user-server/pkg/handler"
)

func main() {
	configPath := "etc/user_server.yaml"
	config := conf.LoadConfig(configPath)

	appCtx := &app.Context{}
	appCtx.Init(config)
	handler.RegisterUserApiHandlers(appCtx)

	appCtx.StartServe()
}
