package etc

import (
	"lazygo/common/auth"

	"github.com/weblazy/core/config"
	"github.com/weblazy/core/database/redis"
	"github.com/weblazy/core/rpcx"
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
