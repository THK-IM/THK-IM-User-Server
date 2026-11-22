package errorx

import (
	"github.com/thk-im/thk-im-base-server/errorx"
)

var UserTokenError = errorx.NewErrorX(4001001, "user token error")
var UserNotExisted = errorx.NewErrorX(4001002, "user not existed")

var InvalidCredentials = errorx.NewErrorX(4001003, "invalid account or password")
