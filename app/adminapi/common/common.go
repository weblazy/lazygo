package common

import (
	"github.com/weblazy/core/database/redis"
	"github.com/weblazy/core/database/sqlx"
)

var (
	Conn     sqlx.SqlConn
	BizRedis *redis.Redis
)
