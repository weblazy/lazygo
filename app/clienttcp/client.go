package main

import (
	tp "github.com/weblazy/teleport"
	"lazygo/core/logx"
	"time"
)

func main() {
	defer tp.SetLoggerLevel("ERROR")()

	cli := tp.NewPeer(tp.PeerConfig{}, nil)
	defer cli.Close()
	// cli.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	cli.RoutePush(new(Client))

	sess, stat := cli.Dial(":5555")
	if !stat.OK() {
		tp.Fatalf("%v", stat)
	}
	stat = sess.Push(
		"/client/receive",
		"send",
	)

	for {
		time.Sleep(10 * time.Minute)
	}
}

// Push push handler
type Client struct {
	tp.PushCtx
}

// Push handles '/push/status' message
func (c *Client) Onreceive(arg *string) *tp.Status {
	logx.Errorf("收到数据:%s", *arg)
	return nil
}
