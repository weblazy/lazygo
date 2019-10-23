package main

import (
	// "fmt"
	"github.com/weblazy/teleport"
	"lazygo/core/logx"
	"lazygo/core/tpcluster"
	// "time"
)

func main() {
	defer tp.FlushLogger()
	// graceful
	go tp.GraceSignal()

	// server peer
	controllers := make([]interface{}, 0)
	controllers = append(controllers, new(Client))
	tpcluster.StartNode(tpcluster.NodeConf{
		TransPeerConf: tp.PeerConfig{
			LocalIP:     "127.0.0.1",
			CountTime:   true,
			ListenPort:  8080,
			PrintDetail: true,
		},
		ClientPeerConf: tp.PeerConfig{
			LocalIP:     "127.0.0.1",
			CountTime:   true,
			ListenPort:  5555,
			PrintDetail: true,
		},
		MasterAddress: "127.0.0.1:9090",
		Password:      "skdss",
	}, controllers)

}

type Client struct {
	tp.PushCtx
}

// Add handles addition request
func (c *Client) Receive(args *string) *tp.Status {
	logx.Error(*args)
	tpcluster.BindUid("Client1", c.PushCtx)
	tpcluster.SendToUid("Client1", "/client/onreceive", "send message ok")
	return nil
}
