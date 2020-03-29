package main

import (
	"flag"
	"fmt"
	"lazygo/app/configurerpc/configureproto"
	"lazygo/app/configurerpc/etc"
	"log"

	"github.com/weblazy/core/config"
	"github.com/weblazy/core/logx"
	"github.com/weblazy/core/rpcx"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/config.json", "the config file")

func main() {
	flag.Parse()
	config.UnmarshalWithLog(*configFile, &etc.Conf)
	cs, err := configureproto.NewConfigureHandler()
	if err != nil {
		log.Fatal(err)
	}

	server, err := rpcx.NewRpcServer(etc.Conf.RpcServerConf, func(grpcServer *grpc.Server) {
		configureproto.RegisterConfigureHandlerServer(grpcServer, cs)
	})
	if err != nil {
		logx.Fatal(err)
	}
	fmt.Printf("Starting rpc server at %s...\n", etc.Conf.RpcServerConf.ListenOn)
	server.Start()
}
