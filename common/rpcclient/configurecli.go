package rpcclient

import (
	"fmt"
	"lazygo/app/configurerpc/configureproto"
	"lazygo/core/rpcx"
)

var (
	ConfigureCli configureproto.ConfigureHandlerClient
	client       *rpcx.DirectClient
)

func Init(configureRpc rpcx.RpcClientConf) error {
	cli, err := rpcx.NewDirectClient(configureRpc)
	if err != nil {
		return err
	}
	client = cli
	return nil
}

func Next() (configureproto.ConfigureHandlerClient, error) {
	conn, exit := client.Next()
	if !exit {
		return nil, fmt.Errorf("configureRpc client not exit")
	}
	return configureproto.NewConfigureHandlerClient(conn), nil
}
