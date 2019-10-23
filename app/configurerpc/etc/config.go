package etc

import (
	"lazygo/core/logx"
	"lazygo/core/rpcx"
)

type (
	Config struct {
		RpcServerConf rpcx.RpcServerConf
		Log           logx.Config
	}
)

var Conf Config
