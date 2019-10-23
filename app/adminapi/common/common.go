package common

import (
	"lazygo/core/database/redis"
	"lazygo/core/database/sqlx"
)

var (
	Conn     sqlx.SqlConn
	BizRedis *redis.Redis
)
