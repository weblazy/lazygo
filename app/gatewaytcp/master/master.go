package main

import (
	// "fmt"
	"github.com/henrylee2cn/teleport"
	"lazygo/core/tpcluster"
	// "time"
)

func main() {
	defer tp.FlushLogger()
	// graceful
	go tp.GraceSignal()

	// server peer
	tpcluster.StartMaster(tpcluster.MasterConf{
		MasterPeerConf: tp.PeerConfig{
			LocalIP:     "127.0.0.1",
			ListenPort:  9090,
			CountTime:   true,
			PrintDetail: true,
		},
		Password: "skdss",
	})

}
