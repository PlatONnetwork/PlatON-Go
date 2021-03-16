package p2p

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	monitorSchedulerOnce sync.Once
	monitorSchedulerRef  *monitorScheduler
)

type monitorScheduler struct {
	queue                   []*monitorTask
	disconnectMonitorPeerFn disconnectMonitorPeerFn
}

type monitorTask struct {
	flags        connFlag
	dest         *discover.Node
	lastResolved time.Time
	resolveDelay time.Duration
	err          error
}

func MonitorScheduler() *monitorScheduler {
	monitorSchedulerOnce.Do(func() {
		log.Info("Init node monitor scheduler ...")
		monitorSchedulerRef = &monitorScheduler{}
	})
	return monitorSchedulerRef
}

func (tasks *monitorScheduler) InitDisconnectMonitorPeerFn(disconnectMonitorPeerFn disconnectMonitorPeerFn) {
	tasks.disconnectMonitorPeerFn = disconnectMonitorPeerFn
}

// whether the queue is empty
func (tasks *monitorScheduler) isEmpty() bool {
	if len(tasks.queue) == 0 {
		return true
	}
	return false
}

func (tasks *monitorScheduler) ClearMonitorScheduler() {
	tasks.queue = []*monitorTask{}
}

func (tasks *monitorScheduler) AddMonitorTask(task *monitorTask) {
	tasks.queue = append(tasks.queue, task)
}

/*func (tasks *monitorScheduler) AddMonitorNode(node *discover.Node) {
	tasks.queue = append(tasks.queue, &monitorTask{flags: monitorConn, dest: node})
}
*/
/*func (tasks *monitorScheduler) RenewMonitorSchedule(nodeIdList []discover.NodeID) {
	tasks.queue = []*monitorTask{}
	for _, nodeId := range nodeIdList {
		node := discover.NewNode(nodeId, nil, 0, 0)
		tasks.queue = append(tasks.queue, &monitorTask{flags: monitorConn, dest: node})
	}
}*/

func (tasks *monitorScheduler) ListTask() []*monitorTask {
	log.Info("list monitor dial task", "task queue", tasks.description())
	return tasks.queue
}

func (tasks *monitorScheduler) RemoveTask(NodeID discover.NodeID) {
	if !tasks.isEmpty() {
		for i, t := range tasks.queue {
			if t.dest.ID == NodeID {
				tasks.queue = append(tasks.queue[:i], tasks.queue[i+1:]...)
				break
			}
		}
	}
}

func (tasks *monitorScheduler) description() string {
	var description []string
	for _, t := range tasks.queue {
		description = append(description, fmt.Sprintf("%x", t.dest.ID[:8]))
	}
	return strings.Join(description, ",")
}

func (t *monitorTask) Do(srv *Server) {
	if t.dest.Incomplete() {
		if !t.resolve(srv) {
			//MONITOR：没有找发现节点
			SaveNodePingResult(t.dest.ID, t.dest.IP.String(), strconv.FormatUint(uint64(t.dest.TCP), 10), 2)
			return
		}
	}
	err := t.dial(srv, t.dest)
	if err != nil {
		log.Trace("Dial error", "task", t, "err", err)
		//MONITOR：连接节点错误
		SaveNodePingResult(t.dest.ID, t.dest.IP.String(), strconv.FormatUint(uint64(t.dest.TCP), 10), 2)
		// Try resolving the ID of static nodes if dialing failed.
		if _, ok := err.(*dialError); ok && t.flags&staticDialedConn != 0 {
			if t.resolve(srv) {
				t.dial(srv, t.dest)
			}
		}
	}
	t.err = err
}

//如果节点只是monitor，则关闭，并返回true；否则返回false
func (t *monitorTask) TryDisconnect() bool {
	return MonitorScheduler().disconnectMonitorPeerFn(t.dest)
}

func (t *monitorTask) resolve(srv *Server) bool {
	if srv.ntab == nil {
		log.Debug("Can't resolve node", "id", t.dest.ID, "err", "discovery is disabled")
		return false
	}
	if t.resolveDelay == 0 {
		t.resolveDelay = initialResolveDelay
	}
	if time.Since(t.lastResolved) < t.resolveDelay {
		return false
	}
	resolved := srv.ntab.Resolve(t.dest.ID)
	t.lastResolved = time.Now()
	if resolved == nil {
		t.resolveDelay *= 2
		if t.resolveDelay > maxResolveDelay {
			t.resolveDelay = maxResolveDelay
		}
		log.Debug("Resolving node failed", "id", t.dest.ID, "newdelay", t.resolveDelay)
		return false
	}
	// The node was found.
	t.resolveDelay = initialResolveDelay
	t.dest = resolved
	log.Debug("Resolved node", "id", t.dest.ID, "addr", &net.TCPAddr{IP: t.dest.IP, Port: int(t.dest.TCP)})

	return true
}

// dial performs the actual connection attempt.
func (t *monitorTask) dial(srv *Server, dest *discover.Node) error {
	fd, err := srv.Dialer.Dial(dest)
	if err != nil {
		return &dialError{err}
	}
	mfd := newMeteredConn(fd, false)
	return srv.SetupConn(mfd, t.flags, dest)
}

type Downloading interface {
	HighestBlock() uint64
}

func PostMonitorNodeEvent(eventMux *event.TypeMux, blockNumber uint64, epoch uint64, verifierList []*staking.Validator, downloading Downloading) {
	nodeIdList := ConvertToNodeIdList(verifierList)
	//nodeIdStringList := xcom.ConvertToNodeIdStringList(verifierList)
	//MONITOR，保存下一轮结算周期的新101名单
	SaveEpochElection(epoch+1, nodeIdList)

	if blockNumber > downloading.HighestBlock() {
		//说明区块是共识协议得到的，需要执行monitor任务
		//MONITOR，保存新101名单
		InitNodePing(nodeIdList)
		eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdList})
	} else {
		//此次需要同步的最高块，可以认为是链上当前块，如果当前结算周期和链上结算周期一致，则开始执行monitor任务
		chainEpoch := xutil.CalculateEpoch(downloading.HighestBlock())
		if chainEpoch == epoch {
			InitNodePing(nodeIdList)
			eventMux.Post(cbfttypes.ElectNextEpochVerifierEvent{NodeIdList: nodeIdList})
		}
	}
}

func ConvertToNodeIdStringList(verifierList []*staking.Validator) []string {
	nodeIdStringList := make([]string, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdStringList[i] = verifier.NodeId.String()
	}
	return nodeIdStringList
}

func ConvertToNodeIdList(verifierList []*staking.Validator) []discover.NodeID {
	nodeIdList := make([]discover.NodeID, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdList[i] = verifier.NodeId
	}
	return nodeIdList
}

/*func (t *monitorTask) dialAndClose(srv *Server, dest *discover.Node) error {
	fd, err := srv.Dialer.Dial(dest)
	if err != nil {
		return &dialError{err}
	}
	fd.Close()
	return nil
}
*/
