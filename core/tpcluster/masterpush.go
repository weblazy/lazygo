package tpcluster

import (
	// "encoding/json"
	// "fmt"
	"github.com/weblazy/teleport"
	"lazygo/core/logx"
)

type (
	MasterPush struct {
		tp.PushCtx
	}
)

func (m *MasterPush) Ping(ping *string) *tp.Status {
	sessionId := m.Session().ID()
	logx.Errorf("%s:%s", sessionId, *ping)
	return nil
}

func (m *MasterPush) OnClientConnect(ping *string) *tp.Status {
	sessionId := m.Session().ID()
	logx.Errorf("%s:%s", sessionId, *ping)
	return nil
}

func (m *MasterPush) OnClientClose(ping *string) *tp.Status {
	sessionId := m.Session().ID()
	logx.Errorf("%s:%s", sessionId, *ping)
	return nil
}
