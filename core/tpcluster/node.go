package tpcluster

import (
	"github.com/henrylee2cn/teleport"
	"lazygo/core/logx"
	"lazygo/core/mapreduce"
	"lazygo/core/timingwheel"
	"strconv"
	"time"
)

type (
	NodeConf struct {
		ClientPeerConf tp.PeerConfig
		TransPeerConf  tp.PeerConfig
		TransPort      int64
		MasterAddress  string
		Password       string
	}

	Test struct {
		Reloadable             bool
		PingInterval           int64
		PingNotResPonseLimit   int64
		PingData               string
		CecreteKey             string
		Router                 func()
		SendToWorkerBufferSize int64
		SendToClientBufferSize int64
		nodeSessions           map[string]tp.Session
		startTime              time.Time
		// gatewayConnections map[string]string
		// businessConnections map[string]tp.Session
	}

	NodeInfo struct {
		nodeConf        NodeConf
		nodeSessions    map[string]tp.CtxSession
		clientSessions  map[string]tp.Session
		uidSessions     map[string]map[string]tp.Session
		groupSessions   map[string]map[string]tp.Session
		clientIdBindUid map[string]string
		clientPeer      tp.Peer
		clientAddress   string
		transPeer       tp.Peer
		transAddress    string
		timer           *timingwheel.TimingWheel
		startTime       time.Time
	}
)

var (
	nodeInfo NodeInfo
)

const (
	PERSISTENCE_CONNECTION_PING_INTERVAL = 25
)

// NewPeer creates a new peer.
func StartNode(cfg NodeConf, controllers []interface{}, globalLeftPlugin ...tp.Plugin) {
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
	nodeInfo = NodeInfo{
		nodeConf:        cfg,
		nodeSessions:    make(map[string]tp.CtxSession),
		clientSessions:  make(map[string]tp.Session),
		uidSessions:     make(map[string]map[string]tp.Session),
		groupSessions:   make(map[string]map[string]tp.Session),
		clientIdBindUid: make(map[string]string),
		startTime:       time.Now(),
		timer:           timer,
	}

	nodeInfo.transPeer = tp.NewPeer(cfg.TransPeerConf, globalLeftPlugin...)
	nodeInfo.clientPeer = tp.NewPeer(cfg.ClientPeerConf, globalLeftPlugin...)
	for _, value := range controllers {
		nodeInfo.clientPeer.RoutePush(value)
	}
	nodeInfo.transPeer.RouteCall(new(Node))
	go func() {
		nodeInfo.transPeer.ListenAndServe()
	}()
	go func() {
		sess, stat := nodeInfo.transPeer.Dial(cfg.MasterAddress)
		if !stat.OK() {
			tp.Fatalf("%v", stat)
		}
		port := strconv.FormatInt(int64(cfg.TransPeerConf.ListenPort), 10)
		nodeInfo.transAddress = cfg.TransPeerConf.LocalIP + ":" + port
		var result int
		auth := &Auth{
			Password:     nodeInfo.nodeConf.Password,
			TransAddress: nodeInfo.transAddress,
		}
		stat = sess.Call("/master/auth",
			auth,
			&result,
		).Status()
	}()
	nodeInfo.clientPeer.ListenAndServe()
}

type Node struct {
	tp.CallCtx
}

// Add handles addition request
func (n *Node) Updatenodelist(nodeList *[]string) (int64, *tp.Status) {
	for _, value := range *nodeList {
		sess, stat := nodeInfo.transPeer.Dial(value)
		if !stat.OK() {
			tp.Fatalf("%v", stat)
		}

		var result int
		auth := &Auth{
			Password:     nodeInfo.nodeConf.Password,
			TransAddress: nodeInfo.transAddress,
		}
		stat = sess.Call("/node/auth",
			auth,
			&result,
		).Status()
	}

	return 0, nil
}

func (n *Node) Auth(args *Auth) (int, *tp.Status) {
	session := n.Session()
	sessionId := session.ID()
	peer := n.Peer()
	psession, ok := peer.GetSession(sessionId)
	if args.Password != nodeInfo.nodeConf.Password && ok {
		logx.Errorf("密码错误，非法链接:%s", sessionId)
		psession.Close()
		return 0, nil
	}
	nodeInfo.nodeSessions[sessionId] = session
	return 1, nil
}

func GetSession(context tp.PreCtx) tp.Session {
	sid := context.Session().ID()
	session, _ := context.Peer().GetSession(sid)
	return session
}

func BindUid(uid string, context tp.PreCtx) (int, *tp.Status) {
	sid := context.Session().ID()
	session, _ := context.Peer().GetSession(sid)
	if oldUid, ok := nodeInfo.clientIdBindUid[sid]; ok {
		if oldUid != uid {
			delete(nodeInfo.uidSessions[oldUid], sid)
		} else {
			return 0, nil
		}
	}
	nodeInfo.clientIdBindUid[sid] = uid
	if _, ok := nodeInfo.uidSessions[uid]; ok {
		nodeInfo.uidSessions[uid][sid] = session
	} else {
		sessions := make(map[string]tp.Session)
		sessions[sid] = session
		nodeInfo.uidSessions[uid] = sessions
	}

	return 0, nil
}

func SendToUid(uid string, path string, req interface{}) (int, *tp.Status) {
	if sessions, ok := nodeInfo.uidSessions[uid]; ok {
		mapreduce.MapVoid(func(source chan<- interface{}) {
			for _, session := range sessions {
				source <- session
			}
		}, func(item interface{}) {
			session := item.(tp.Session)
			session.Push(
				path,
				req,
			)
		})
	}
	return 0, nil
}

func JoinGroup(gid string, session tp.Session) (int, *tp.Status) {
	sid := session.ID()
	if _, ok := nodeInfo.groupSessions[gid]; ok {
		nodeInfo.groupSessions[gid][sid] = session
	} else {
		sessions := make(map[string]tp.Session, 0)
		sessions[sid] = session
		nodeInfo.groupSessions[gid] = sessions
	}
	return 0, nil
}

func LeaveGroup(gid string, session tp.Session) (int, *tp.Status) {
	sid := session.ID()
	if _, ok := nodeInfo.groupSessions[gid]; ok {
		delete(nodeInfo.groupSessions[gid], sid)
	}
	return 0, nil
}

func SendToGroup(gid string, path string, req interface{}) (int, *tp.Status) {
	if sessions, ok := nodeInfo.groupSessions[gid]; ok {
		mapreduce.MapVoid(func(source chan<- interface{}) {
			for _, session := range sessions {
				source <- session
			}
		}, func(item interface{}) {
			session := item.(tp.Session)
			session.Push(
				path,
				req,
			)
		})
	}
	return 0, nil
}

