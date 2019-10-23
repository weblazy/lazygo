package main

import (
	"flag"
	"lazygo/app/adminapi/common"
	"lazygo/app/adminapi/controller"
	"lazygo/app/adminapi/etc"
	"lazygo/common/auth"
	"lazygo/common/rpcclient"
	"lazygo/core/apix"
	"lazygo/core/config"
	"lazygo/core/database/redis"
	"lazygo/core/database/sqlx"
	"lazygo/core/logx"
	_ "net/http/pprof"
)

var configFile = flag.String("f", "etc/config.json", "the config file")

func main() {
	flag.Parse()
	config.UnmarshalWithLog(*configFile, &etc.Conf)
	err := rpcclient.Init(etc.Conf.ConfigureRpc)
	if err != nil {
		logx.Error(err)
	}
	auth.InitAuth(auth.AuthConfig{
		RedisNodeList: etc.Conf.AuthList,
	})
	common.Conn = sqlx.NewMysql(etc.Conf.MysqlDataSource)
	common.BizRedis = redis.NewRedis(etc.Conf.RedisConf.Host, etc.Conf.RedisConf.Type, etc.Conf.RedisConf.Pass)
	routerMap := map[string]apix.ControllerInterface{
		"/Index/": &controller.IndexController{},
	}
	apix.Run(etc.Conf.ApiConfig, routerMap)
}
