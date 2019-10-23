package tpcluster

import (
	// "encoding/json"
	// "fmt"
	"github.com/henrylee2cn/teleport"
	"lazygo/core/logx"
	"lazygo/core/timingwheel"
	"time"
)

type (
	MasterConf struct {
		MasterPeerConf tp.PeerConfig
		TransPort      int64
		MasterAddress  string
		Password       string
	}
	Master struct {
		tp.CallCtx
	}
	MasterInfo struct {
		masterConf   MasterConf
		nodeSessions map[string]tp.CtxSession
		nodeAddress  map[string]string
		timer        *timingwheel.TimingWheel
		startTime    time.Time
	}
	Auth struct {
		TransAddress string
		Password     string
	}
)

var (
	masterInfo MasterInfo
)

// 启动master节点.
func StartMaster(cfg MasterConf, globalLeftPlugin ...tp.Plugin) {
	timer, err := timingwheel.NewTimingWheel(time.Second, 300, func(k, v interface{}) {
		logx.Errorf("%s auth timeout", k)
		err := v.(tp.Session).Close()
		if err != nil {
			logx.Error(err)
		}
	})
	defer timer.Stop()
	if err != nil {
		logx.Fatal(err)
	}
	masterInfo = MasterInfo{
		masterConf:   cfg,
		nodeSessions: make(map[string]tp.CtxSession),
		nodeAddress:  make(map[string]string),
		startTime:    time.Now(),
		timer:        timer,
	}
	peer := tp.NewPeer(cfg.MasterPeerConf,func(session *tp.Session){

	}, globalLeftPlugin...)
	master := new(Master)
	peer.RouteCall(master)
	peer.OnConnect = func(p *tp.Peer)(){

	}
	peer.ListenAndServe()

}

// func (m *Master) OnConnect(session tp.Session) {
// 	Timer.SetTimer(session.ID(), session, 10*time.Second)
// }

// func (m *Master) OnMessage(session tp.Session, data Data) {
// 	switch data.Event {
// 	case "node_connect":
// 		sessionId := session.ID()
// 		nodeConnections[sessionId] = session
// 		m.broadcastAddresses()
// 	default:
// 		session.Close()
// 	}
// }

// func (m *Master) OnClose(session tp.Session) {
// 	sessionId := session.ID()
// 	if _, ok := nodeConnections[sessionId]; ok {
// 		delete(nodeConnections, sessionId)
// 		m.broadcastAddresses()
// 	}
// }

func (m *Master) Auth(args *Auth) (int, *tp.Status) {
	session := m.Session()
	sessionId := session.ID()

	peer := m.Peer()
	psession, ok := peer.GetSession(sessionId)
	if args.Password != masterInfo.masterConf.Password && ok {
		logx.Errorf("密码错误，非法链接:%s", sessionId)
		psession.Close()
		return 0, nil
	}
	masterInfo.nodeSessions[sessionId] = session
	masterInfo.nodeAddress[sessionId] = args.TransAddress
	// go func(){
	// 	select{
	// 	case <- m.Session().CloseNotify():
	// 		m.OnClose()
	// }
	// }()
	m.broadcastAddresses()
	return 1, nil
}

func (m *Master) broadcastAddresses() {
	nodeList := make([]string, 0)
	for _, value := range masterInfo.nodeAddress {
		nodeList = append(nodeList, value)
	}
	var result int
	for _, value := range masterInfo.nodeSessions {
		value.Call(
			"/node/updatenodelist",
			nodeList,
			&result,
		)
	}
}

func (m *Master)OnClose(){
	
	
}