package etc

import (
	"github.com/weblazy/core/logx"
	"github.com/weblazy/core/rpcx"
)

type (
	Config struct {
		RpcServerConf rpcx.RpcServerConf
		Log           logx.Config
	}
)

var Conf Config
