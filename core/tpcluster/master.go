package tpcluster

import (
	// "encoding/json"
	// "fmt"
	"github.com/henrylee2cn/goutil"
	"github.com/weblazy/teleport"
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
	MasterInfo struct {
		masterConf MasterConf
		nodeMap    goutil.Map
		timer      *timingwheel.TimingWheel
		startTime  time.Time
	}
	nodeSession struct {
		session tp.Session
		address string
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
		masterConf: cfg,
		nodeMap:    goutil.AtomicMap(),
		startTime:  time.Now(),
		timer:      timer,
	}
	peer := tp.NewPeer(cfg.MasterPeerConf, globalLeftPlugin...)
	peer.RouteCall(new(MasterCall))
	peer.RoutePush(new(MasterPush))
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

func (mi *MasterInfo) broadcastAddresses() {
	nodeList := make([]string, 0)
	mi.nodeMap.Range(func(k interface{}, v interface{}) bool {
		nodeList = append(nodeList, v.(nodeSession).address)
		return true
	})
	var result int
	len := mi.nodeMap.Len()
	callCmdChan := make(chan tp.CallCmd, len)
	mi.nodeMap.Range(func(k interface{}, v interface{}) bool {
		v.(nodeSession).session.AsyncCall(
			"/node_call/updatenodelist",
			nodeList,
			&result,
			callCmdChan,
		)
		return true
	})
	for callCmd := range callCmdChan {
		_, _ = callCmd.Reply()
	}
}

// set sets a *session.
func (mi *MasterInfo) setSession(sess tp.Session, address string) {
	sid := sess.ID()
	node := &nodeSession{
		address: address,
		session: sess,
	}
	_node, loaded := mi.nodeMap.LoadOrStore(sid, node)
	if !loaded {
		return
	}
	mi.nodeMap.Store(sid, node)
	if oldSess := _node.(*nodeSession).session; sess != oldSess {
		oldSess.Close()
	}
}
