package etc

import (
	"lazygo/common/auth"
	"lazygo/core/config"
	"lazygo/core/database/redis"
	"lazygo/core/rpcx"
)

type (
	Config struct {
		config.ApiConfig
		ConfigureRpc    rpcx.RpcClientConf
		RedisConf       redis.RedisConf
		MysqlDataSource string
		AuthList        []auth.RedisNode
	}
)

var Conf Config
