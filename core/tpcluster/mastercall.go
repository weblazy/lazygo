package tpcluster

import (
	// "encoding/json"
	// "fmt"
	"github.com/weblazy/teleport"
	"lazygo/core/logx"
)

type (
	MasterCall struct {
		tp.CallCtx
	}
	Auth struct {
		TransAddress string
		Password     string
	}
)

func (m *MasterCall) Auth(args *Auth) (int, *tp.Status) {
	session := m.Session()
	sessionId := session.ID()

	peer := m.Peer()
	psession, ok := peer.GetSession(sessionId)
	if args.Password != masterInfo.masterConf.Password && ok {
		logx.Errorf("密码错误，非法链接:%s", sessionId)
		psession.Close()
		return 0, nil
	}
	masterInfo.setSession(psession, args.TransAddress)
	masterInfo.broadcastAddresses()
	return 1, nil
}
