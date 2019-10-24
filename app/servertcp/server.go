package main

import (
	"fmt"
	"time"

	tp "github.com/weblazy/teleport"
)

func main() {
	defer tp.FlushLogger()
	// graceful
	go tp.GraceSignal()

	// server peer
	srv := tp.NewPeer(tp.PeerConfig{
		CountTime:   true,
		ListenPort:  9090,
		PrintDetail: true,
	})
	// srv.SetTLSConfig(tp.GenerateTLSConfigForServer())

	// router
	srv.RouteCall(new(Math))
	srv.RouteCall(new(Test))
	// broadcast per 5s
	go func() {
		for {
			time.Sleep(time.Second * 5)
			srv.RangeSession(func(sess tp.Session) bool {
				sess.Push(
					"/push/status",
					fmt.Sprintf("this is a broadcast, server time: %v", time.Now()),
				)
				return true
			})
		}
	}()

	// listen and serve
	srv.ListenAndServe()
}

// Math handler
type Math struct {
	tp.CallCtx
}

// Add handles addition request
func (m *Math) Add(arg *[]int) (int, *tp.Status) {
	// test meta
	tp.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}

// Math handler
type Test struct {
	tp.CallCtx
}

// Add handles addition request
func (m *Test) Test(arg *[]int) (int, *tp.Status) {
	// test meta
	tp.Infof("test: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
